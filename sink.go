package panyl_zap

import "github.com/RangelReale/panyl"

type Sink struct {
	job *panyl.Job
}

func NewSink(job *panyl.Job) *Sink {
	return &Sink{
		job: job,
	}
}

func (s *Sink) Write(p []byte) (n int, err error) {
	return len(p), s.job.ProcessLine(string(p))
}

func (s *Sink) Sync() error {
	return nil
}

func (s *Sink) Close() error {
	return s.job.Finish()
}
