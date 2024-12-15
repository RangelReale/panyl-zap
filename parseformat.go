package panylzap

import (
	"context"
	"fmt"
	"time"

	"github.com/RangelReale/panyl/v2"
	"go.uber.org/zap/zapcore"
)

var _ panyl.PluginParseFormat = (*ZapJSON)(nil)

const ZapJSONFormat = "zap_json"

type ZapJSON struct {
	EncoderConfig zapcore.EncoderConfig
}

// example: {"level":"debug","timestamp":"2022-08-05T18:38:26.088Z","caller":"aws/client.go:10","message":"client response"}

var (
	zapTimestampFormat = "2006-01-02T15:04:05.000Z07:00"
)

func (c ZapJSON) ParseFormat(ctx context.Context, item *panyl.Item) (bool, error) {
	// only if json
	if item.Metadata.StringValue(panyl.MetadataStructure) == panyl.MetadataStructureJSON {
		if (c.EncoderConfig.MessageKey == "" || item.Data.HasValue(c.EncoderConfig.MessageKey)) &&
			(c.EncoderConfig.LevelKey == "" || item.Data.HasValue(c.EncoderConfig.LevelKey)) &&
			(c.EncoderConfig.TimeKey == "" || item.Data.HasValue(c.EncoderConfig.TimeKey)) {

			item.Metadata[panyl.MetadataFormat] = ZapJSONFormat

			if c.EncoderConfig.MessageKey != "" && item.Data.HasValue(c.EncoderConfig.MessageKey) {
				item.Metadata[panyl.MetadataMessage] = item.Data.StringValue(c.EncoderConfig.MessageKey)
			}

			if c.EncoderConfig.NameKey != "" && item.Data.HasValue(c.EncoderConfig.NameKey) {
				item.Metadata[panyl.MetadataCategory] = item.Data.StringValue(c.EncoderConfig.NameKey)
			} else if c.EncoderConfig.CallerKey != "" && item.Data.HasValue(c.EncoderConfig.CallerKey) {
				item.Metadata[panyl.MetadataCategory] = item.Data.StringValue(c.EncoderConfig.CallerKey)
			}

			if c.EncoderConfig.TimeKey != "" && item.Data.HasValue(c.EncoderConfig.TimeKey) {
				ts, err := time.Parse(zapTimestampFormat, item.Data.StringValue(c.EncoderConfig.TimeKey))
				if err == nil {
					item.Metadata[panyl.MetadataTimestamp] = ts
				}
			}

			if c.EncoderConfig.LevelKey != "" && item.Data.HasValue(c.EncoderConfig.LevelKey) {
				switch item.Data.StringValue(c.EncoderConfig.LevelKey) {
				case "error", "ERROR", "fatal", "FATAL", "panic", "PANIC", "dpanic", "DPANIC":
					item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelERROR
				case "warn", "WARN":
					item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
				case "info", "INFO":
					item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
				case "debug", "DEBUG":
					item.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
				}
			}

			// if level is error, and have an "error" key, show it as the error message
			if item.Metadata.HasValue(panyl.MetadataMessage) &&
				item.Metadata.HasValue(panyl.MetadataLevel) &&
				(item.Metadata[panyl.MetadataLevel] == panyl.MetadataLevelERROR) &&
				item.Data.HasValue("error") {
				item.Metadata[panyl.MetadataMessage] = fmt.Sprintf("%s [error: %s]",
					item.Metadata[panyl.MetadataMessage], item.Data.StringValue("error"))
			}

			return true, nil
		}
	}
	return false, nil
}

func (c ZapJSON) IsPanylPlugin() {}
