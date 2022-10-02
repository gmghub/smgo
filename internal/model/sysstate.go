package model

import "time"

type SysStat struct {
	// System load average %
	LoadAvg float32 `json:"loadavg"`
}

type CPUStat struct {
	// CPU user mode average %
	UserModeAvg float32 `json:"usermodeavg"`
	// CPU system mode average %
	SysModeAvg float32 `json:"sysmodeavg"`
	// CPU idle average %
	IdleAvg float32 `json:"idleavg"`
}

type DiskStat struct {
	Time time.Time
	// TODO
	// Disk trasfers per sec
	// 	DiskTps float32
	// 	// Disk KB/s
	// 	DiskRWkbs         float32
	// 	DiskUsedMb        float32
	// 	DiskUsedPct       float32
	// 	DiskInodesUsed    float32
	// 	DiskInodesUsedPct float32
}
