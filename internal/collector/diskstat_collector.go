package collector

import "time"

const (
	CollectorNameDiskStat = "diskstat"
	// Base period to run collector.
	collectorPeriodDiskStat = 1
)

type DiskStat struct {
	Time time.Time
}

func NewDiskStatCollector(size int) *Collector {
	return &Collector{
		name:   "cpustat",
		period: collectorPeriodDiskStat,
		buffer: *NewRingBuffer(size),
		fun: func() interface{} {
			return GetDiskStat()
		},
		funStatJSON: DiskStatJSON,
	}
}

func DiskStatJSON(c *Collector, period int) []byte {
	return nil
}

func GetDiskStat() DiskStat {
	stat := DiskStat{Time: time.Now().UTC()}
	return stat
}
