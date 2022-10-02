package collector

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gmghub/smgo/internal/model"
)

const (
	CollectorNameCPUStat = "cpustat"
	// Base period to run collector.
	collectorPeriodCPUStat = 1
)

type CPUStat struct {
	Time time.Time
	// CPU user mode %
	UserMode float32
	// CPU system mode %
	SysMode float32
	// CPU idle %
	Idle float32
}

func NewCPUStatCollector(size int) *Collector {
	return &Collector{
		name:   "cpustat",
		period: collectorPeriodCPUStat,
		buffer: *NewRingBuffer(size),
		fun: func() interface{} {
			return GetCPUStat()
		},
		funStatJSON: CPUStatJSON,
	}
}

func CPUStatJSON(c *Collector, period int) []byte {
	stats := c.buffer.GetN(period)

	if len(stats) < period {
		return nil
	}
	var sumuser, sumsys, sumidle float32
	for _, s := range stats {
		sumuser += s.(CPUStat).UserMode
		sumsys += s.(CPUStat).SysMode
		sumidle += s.(CPUStat).Idle
	}
	stat := model.CPUStat{
		UserModeAvg: sumuser / float32(period),
		SysModeAvg:  sumsys / float32(period),
		IdleAvg:     sumidle / float32(period),
	}

	j, err := json.Marshal(stat)
	if err != nil {
		log.Println(CollectorNameCPUStat, ":", err)
		return nil
	}
	return j
}
