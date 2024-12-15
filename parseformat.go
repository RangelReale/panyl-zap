package panylzap

import (
	"context"
	"fmt"
	"time"

	"github.com/RangelReale/panyl/v2"
	"go.uber.org/zap/zapcore"
)

const ZapJSONFormat = "zap_json"

type ZapJSON struct {
	EncoderConfig zapcore.EncoderConfig
}

// example: {"level":"debug","timestamp":"2022-08-05T18:38:26.088Z","caller":"aws/client.go:10","message":"client response"}

var (
	zapTimestampFormat = "2006-01-02T15:04:05.000Z07:00"
)

var _ panyl.PluginParseFormat = ZapJSON{}

func (m ZapJSON) ParseFormat(ctx context.Context, item *panyl.Item) (bool, error) {
	// only if json
	if item.Metadata.StringValue(panyl.MetadataStructure) == panyl.MetadataStructureJSON {
		if (m.EncoderConfig.MessageKey == "" || item.Data.HasValue(m.EncoderConfig.MessageKey)) &&
			(m.EncoderConfig.LevelKey == "" || item.Data.HasValue(m.EncoderConfig.LevelKey)) &&
			(m.EncoderConfig.TimeKey == "" || item.Data.HasValue(m.EncoderConfig.TimeKey)) {

			item.Metadata[panyl.MetadataFormat] = ZapJSONFormat

			if m.EncoderConfig.MessageKey != "" && item.Data.HasValue(m.EncoderConfig.MessageKey) {
				item.Metadata[panyl.MetadataMessage] = item.Data.StringValue(m.EncoderConfig.MessageKey)
			}

			if m.EncoderConfig.NameKey != "" && item.Data.HasValue(m.EncoderConfig.NameKey) {
				item.Metadata[panyl.MetadataCategory] = item.Data.StringValue(m.EncoderConfig.NameKey)
			} else if m.EncoderConfig.CallerKey != "" && item.Data.HasValue(m.EncoderConfig.CallerKey) {
				item.Metadata[panyl.MetadataCategory] = item.Data.StringValue(m.EncoderConfig.CallerKey)
			}

			if m.EncoderConfig.TimeKey != "" && item.Data.HasValue(m.EncoderConfig.TimeKey) {
				ts, err := time.Parse(zapTimestampFormat, item.Data.StringValue(m.EncoderConfig.TimeKey))
				if err == nil {
					item.Metadata[panyl.MetadataTimestamp] = ts
				}
			}

			if m.EncoderConfig.LevelKey != "" && item.Data.HasValue(m.EncoderConfig.LevelKey) {
				switch item.Data.StringValue(m.EncoderConfig.LevelKey) {
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

func (m ZapJSON) IsPanylPlugin() {}
