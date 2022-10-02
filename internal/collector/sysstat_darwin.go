//go:build darwin
// +build darwin

package collector

import (
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

func GetSysStat() SysStat {
	stat := SysStat{Time: time.Now().UTC()}

	// top -l1 | grep -E "^Load Avg"
	// Load Avg: 1.81, 1.74, 1.87
	// out, err := exec.Command("top", "-l1").Output()
	// sysctl -n vm.loadavg
	// { 1.81 1.74 1.87 }
	out, err := exec.Command("sysctl", "-n", "vm.loadavg").Output()
	if err != nil {
		log.Println(CollectorNameSysStat, ":", err)
		return stat
	}
	rex := regexp.MustCompile(`^{\s+([0-9.]+)\s+`)
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
