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

func (c ZapJSON) ParseFormat(ctx context.Context, result *panyl.Process) (bool, error) {
	// only if json
	if result.Metadata.StringValue(panyl.MetadataStructure) == panyl.MetadataStructureJSON {
		if (c.EncoderConfig.MessageKey == "" || result.Data.HasValue(c.EncoderConfig.MessageKey)) &&
			(c.EncoderConfig.LevelKey == "" || result.Data.HasValue(c.EncoderConfig.LevelKey)) &&
			(c.EncoderConfig.TimeKey == "" || result.Data.HasValue(c.EncoderConfig.TimeKey)) {

			result.Metadata[panyl.MetadataFormat] = ZapJSONFormat

			if c.EncoderConfig.MessageKey != "" && result.Data.HasValue(c.EncoderConfig.MessageKey) {
				result.Metadata[panyl.MetadataMessage] = result.Data.StringValue(c.EncoderConfig.MessageKey)
			}

			if c.EncoderConfig.NameKey != "" && result.Data.HasValue(c.EncoderConfig.NameKey) {
				result.Metadata[panyl.MetadataCategory] = result.Data.StringValue(c.EncoderConfig.NameKey)
			} else if c.EncoderConfig.CallerKey != "" && result.Data.HasValue(c.EncoderConfig.CallerKey) {
				result.Metadata[panyl.MetadataCategory] = result.Data.StringValue(c.EncoderConfig.CallerKey)
			}

			if c.EncoderConfig.TimeKey != "" && result.Data.HasValue(c.EncoderConfig.TimeKey) {
				ts, err := time.Parse(zapTimestampFormat, result.Data.StringValue(c.EncoderConfig.TimeKey))
				if err == nil {
					result.Metadata[panyl.MetadataTimestamp] = ts
				}
			}

			if c.EncoderConfig.LevelKey != "" && result.Data.HasValue(c.EncoderConfig.LevelKey) {
				switch result.Data.StringValue(c.EncoderConfig.LevelKey) {
				case "error", "ERROR", "fatal", "FATAL", "panic", "PANIC", "dpanic", "DPANIC":
					result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelERROR
				case "warn", "WARN":
					result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelWARNING
				case "info", "INFO":
					result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelINFO
				case "debug", "DEBUG":
					result.Metadata[panyl.MetadataLevel] = panyl.MetadataLevelDEBUG
				}
			}

			// if level is error, and have an "error" key, show it as the error message
			if result.Metadata.HasValue(panyl.MetadataMessage) &&
				result.Metadata.HasValue(panyl.MetadataLevel) &&
				(result.Metadata[panyl.MetadataLevel] == panyl.MetadataLevelERROR) &&
				result.Data.HasValue("error") {
				result.Metadata[panyl.MetadataMessage] = fmt.Sprintf("%s [error: %s]",
					result.Metadata[panyl.MetadataMessage], result.Data.StringValue("error"))
			}

			return true, nil
		}
	}
	return false, nil
}

func (c ZapJSON) IsPanylPlugin() {}
