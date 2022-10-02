//go:build windows
// +build windows

package collector

import (
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

func GetSysStat() SysStat {
	var err error
	stat := SysStat{Time: time.Now().UTC()}

	// C:\> wmic cpu get loadpercentage
	// LoadPercentage\r\n10
	out, err := exec.Command("wmic", "cpu", "get", "LoadPercentage").Output()
	if err != nil {
		log.Println(CollectorNameSysStat, ":", err)
		return stat
	}
	rex := regexp.MustCompile(`^LoadPercentage\s+(\d+)`)
	m := rex.FindStringSubmatch(string(out))
	if len(m) < 2 {
		log.Println(CollectorNameSysStat, ":", ErrInvalidData)
		return stat
	}
	var la1m float64
	la1m, err = strconv.ParseFloat(m[1], 32)
	if err != nil {
		log.Println(CollectorNameSysStat, ":", err)
		return stat
	}
	stat.LoadAvg1m = float32(la1m)
	return stat
}
