package collector

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gmghub/smgo/internal/model"
)

const (
	CollectorNameSysStat = "sysstat"
	// Base period to run collector.
	collectorSysStatPeriod = 1
)

type SysStat struct {
	Time time.Time
	// System load average 1m
	LoadAvg1m float32
}

func NewSysStatCollector(size int) *Collector {
	return &Collector{
		name:   CollectorNameSysStat,
		period: collectorSysStatPeriod,
		buffer: *NewRingBuffer(size),
		fun: func() interface{} {
			return GetSysStat()
		},
		funStatJSON: SysStatJSON,
	}
}

func SysStatJSON(c *Collector, period int) []byte {
	stats := c.buffer.GetN(period)
	if len(stats) < period {
		return nil
	}
	var sum float32
	for _, s := range stats {
		sum += s.(SysStat).LoadAvg1m
	}
	stat := model.SysStat{LoadAvg: sum / float32(period)}
	j, err := json.Marshal(stat)
	if err != nil {
		log.Println(CollectorNameSysStat, ":", err)
		return nil
	}
	return j
}
