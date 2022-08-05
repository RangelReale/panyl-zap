package panyl_zap

import (
	"github.com/RangelReale/panyl"
)

type ChannelLineProvider struct {
	c <-chan string
}

func NewChannelLineProvider(c <-chan string) panyl.LineProvider {
	return &ChannelLineProvider{c: c}
}

func (r *ChannelLineProvider) Err() error {
	return nil
}

func (r *ChannelLineProvider) Line() interface{} {
	return <-r.c
}

func (r *ChannelLineProvider) Scan() bool {
	return true
}
