//go:build linux
// +build linux

package collector

import (
	"log"
	"os"
	"strconv"
	"time"
)

func GetSysStat() SysStat {
	stat := SysStat{Time: time.Now().UTC()}

	// /proc/loadavg
	// 2.58 2.27 2.10 5/1513 18119
	cont, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		log.Println(CollectorNameSysStat, ":", err)
		return stat
	}

	var la1m float64
	if la1m, err = strconv.ParseFloat(string(cont[0:4]), 32); err != nil {
		log.Println(CollectorNameSysStat, ":", err)
		return stat
	}
	stat.LoadAvg1m = float32(la1m)

	return stat
}
