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

func GetCPUStat() CPUStat {
	stat := CPUStat{Time: time.Now().UTC()}

	// top -l1 | grep -E "^CPU"
	// CPU usage: 2.56% user, 10.76% sys, 86.66% idle
	out, err := exec.Command("top", "-l1").Output()
	if err != nil {
		log.Println(CollectorNameCPUStat, ":", err)
		return stat
	}
	rex := regexp.MustCompile(`CPU usage:\s+([0-9.]+)%\s+user,\s+([0-9.]+)%\s+sys,\s+([0-9.]+)%\s+idle`)
	m := rex.FindStringSubmatch(string(out))
	if len(m) < 4 {
		log.Println(CollectorNameCPUStat, ":", ErrInvalidData)
		return stat
	}

	var umode, smode, idle float64
	umode, err = strconv.ParseFloat(m[1], 32)
	if err != nil {
		log.Println(CollectorNameCPUStat, ":", err)
	}
	smode, err = strconv.ParseFloat(m[2], 32)
	if err != nil {
		log.Println(CollectorNameCPUStat, ":", err)
	}
	idle, err = strconv.ParseFloat(m[3], 32)
	if err != nil {
		log.Println(CollectorNameCPUStat, ":", err)
	}

	stat.UserMode = float32(umode)
	stat.SysMode = float32(smode)
	stat.Idle = float32(idle)
	return stat
}
