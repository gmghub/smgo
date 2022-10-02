package collector

import (
	"errors"
	"time"
)

var ErrInvalidData = errors.New("invalid data to parse")

type Collector struct {
	name        string
	period      int
	buffer      RingBuffer
	fun         func() interface{}
	funStatJSON func(*Collector, int) []byte
	quit        chan struct{}
}

func (c *Collector) Name() string {
	return c.name
}

func (c *Collector) Start() {
	// run goroutine
	c.quit = make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Duration(c.period) * time.Second)
		for {
			select {
			case <-c.quit:
				ticker.Stop()
				return
			case <-ticker.C:
				c.Collect()
			}
		}
	}()
}

func (c *Collector) Stop() {
	if c.quit != nil {
		close(c.quit)
	}
}

func (c *Collector) Collect() {
	c.buffer.Put(c.fun())
}

func (c *Collector) StatJSON(period int) []byte {
	return c.funStatJSON(c, period)
}
