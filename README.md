# Panyl integration with Uber Zap logger

```go
import (
    "github.com/RangelReale/panyl"
    "github.com/RangelReale/panyl-zap"
    "github.com/RangelReale/panyl/plugins/structure"
    panylzap "github.com/RangelReale/panyl-zap"
    "github.com/RangelReale/panyl-plugins-ansi/output"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func main() {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	
	processor := panyl.NewProcessor(
		panyl.WithPlugins(
			&structure.JSON{},
			&panylzap.ZapJSON{
				EncoderConfig: encoderConfig,
			},
		),
	)
    output := output.NewAnsiOutput(true)
    job := panyl.NewJob(processor, output, panyl.WithIncludeSource(true))

    sink := panylzap.NewSink(job)
    encoder := zapcore.NewJSONEncoder(encoderConfig)
    core := zapcore.NewCore(encoder, zapcore.Lock(sink),  zapcore.DebugLevel)
}	
```

## Author

Rangel Reale (rangelreale@gmail.com)
