package panylzap

import (
	"context"
	"io"

	"github.com/RangelReale/panyl/v2"
	"go.uber.org/zap"
)

type Sink struct {
	job    *panyl.Job
	closer io.Closer
}

var _ zap.Sink = (*Sink)(nil)

type SinkOption func(*Sink)

// NewSink creates a new zap sink using a panyl job
// Must use "zapcore.Lock" to avoid sync issues
func NewSink(job *panyl.Job, options ...SinkOption) *Sink {
	ret := &Sink{
		job: job,
	}
	for _, o := range options {
		o(ret)
	}
	return ret
}

func (s *Sink) Write(p []byte) (n int, err error) {
	return len(p), s.job.ProcessLine(context.Background(), string(p))
}

func (s *Sink) Sync() error {
	return nil
}

func (s *Sink) Close() error {
	err := s.job.Finish(context.Background())
	if err != nil {
		return err
	}
	if s.closer != nil {
		return s.closer.Close()
	}
	return nil
}

func WithCloser(closer io.Closer) SinkOption {
	return func(sink *Sink) {
		sink.closer = closer
	}
}
