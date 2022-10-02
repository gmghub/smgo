//go:build linux
// +build linux

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

	// top -b -n1 | grep -E "Cpu"
	// %Cpu(s): 17.9 us,  3.0 sy,  0.0 ni, 79.1 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
	// CPU:  53% usr   4% sys   0% nic  41% idle   0% io   0% irq   0% sirq
	out, err := exec.Command("top", "-b", "-n1").Output()
	if err != nil {
		log.Println(CollectorNameCPUStat, ":", err)
		return stat
	}
	rex := regexp.MustCompile(`[Cc][Pp][Uu].{0,5}:\s+([0-9.]+)\%?\s+us.\s+([0-9.]+)\%?\s+sy.\s+([0-9.]+)\%?\s+ni.\s+([0-9.]+)\%?\s+id`)
	m := rex.FindStringSubmatch(string(out))
	if len(m) < 5 {
		log.Println(CollectorNameCPUStat, ":", ErrInvalidData)
		return stat
	}

	var umode, smode, nmode, idle float64
	if umode, err = strconv.ParseFloat(m[1], 32); err != nil {
		log.Println(CollectorNameCPUStat, ":", err)
	}
	if smode, err = strconv.ParseFloat(m[2], 32); err != nil {
		log.Println(CollectorNameCPUStat, ":", err)
	}
	if nmode, err = strconv.ParseFloat(m[3], 32); err != nil {
		log.Println(CollectorNameCPUStat, ":", err)
	}
	if idle, err = strconv.ParseFloat(m[4], 32); err != nil {
		log.Println(CollectorNameCPUStat, ":", err)
	}

	stat.UserMode = float32(umode) + float32(nmode)
	stat.SysMode = float32(smode)
	stat.Idle = float32(idle)
	return stat
}
