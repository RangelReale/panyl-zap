# Panyl integration with Uber Zap logger

```go
import (
    "github.com/RangelReale/panyl/v2"
    "github.com/RangelReale/panyl-zap/v2"
    "github.com/RangelReale/panyl/v2/plugins/structure"
    panylzap "github.com/RangelReale/panyl-zap/v2"
    "github.com/RangelReale/panyl-plugins-ansi/v2/output"
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

    logger := zap.New(core)
}
```

## Author

Rangel Reale (rangelreale@gmail.com)
