package panyl_zap

type Sink struct {
	c chan<- string
}

func NewSink(c chan<- string) *Sink {
	return &Sink{
		c: c,
	}
}

func (s *Sink) Write(p []byte) (n int, err error) {
	s.c <- string(p)
	return len(p), nil
}

func (s *Sink) Sync() error {
	return nil
}

func (s *Sink) Close() error {
	return nil
}
