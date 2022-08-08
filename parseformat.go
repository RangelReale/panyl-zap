package panylzap

import (
	"time"

	"github.com/RangelReale/panyl"
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

func (c ZapJSON) ParseFormat(result *panyl.Process) (bool, error) {
	// only if json
	if result.Metadata.StringValue(panyl.Metadata_Structure) == panyl.MetadataStructure_JSON {
		if (c.EncoderConfig.MessageKey == "" || result.Data.HasValue(c.EncoderConfig.MessageKey)) &&
			(c.EncoderConfig.LevelKey == "" || result.Data.HasValue(c.EncoderConfig.LevelKey)) &&
			(c.EncoderConfig.TimeKey == "" || result.Data.HasValue(c.EncoderConfig.TimeKey)) {

			result.Metadata[panyl.Metadata_Format] = ZapJSONFormat

			if c.EncoderConfig.MessageKey != "" && result.Data.HasValue(c.EncoderConfig.MessageKey) {
				result.Metadata[panyl.Metadata_Message] = result.Data.StringValue(c.EncoderConfig.MessageKey)
			}

			if c.EncoderConfig.NameKey != "" && result.Data.HasValue(c.EncoderConfig.NameKey) {
				result.Metadata[panyl.Metadata_Category] = result.Data.StringValue(c.EncoderConfig.NameKey)
			} else if c.EncoderConfig.CallerKey != "" && result.Data.HasValue(c.EncoderConfig.CallerKey) {
				result.Metadata[panyl.Metadata_Category] = result.Data.StringValue(c.EncoderConfig.CallerKey)
			}

			if c.EncoderConfig.TimeKey != "" && result.Data.HasValue(c.EncoderConfig.TimeKey) {
				ts, err := time.Parse(zapTimestampFormat, result.Data.StringValue(c.EncoderConfig.TimeKey))
				if err == nil {
					result.Metadata[panyl.Metadata_Timestamp] = ts
				}
			}

			if c.EncoderConfig.LevelKey != "" && result.Data.HasValue(c.EncoderConfig.LevelKey) {
				switch result.Data.StringValue(c.EncoderConfig.LevelKey) {
				case "fatal", "FATAL", "panic", "PANIC", "dpanic", "DPANIC":
					result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_FATAL
				case "error", "ERROR":
					result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_ERROR
				case "warn", "WARN":
					result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_WARNING
				case "info", "INFO":
					result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_INFO
				case "debug", "DEBUG":
					result.Metadata[panyl.Metadata_Level] = panyl.MetadataLevel_DEBUG
				}
			}
			return true, nil
		}
	}
	return false, nil
}

func (c ZapJSON) IsPanylPlugin() {}
