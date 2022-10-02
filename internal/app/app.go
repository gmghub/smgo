package app

import (
	"fmt"
	"strings"
)

type Collector interface {
	Name() string
	Start()
	Stop()
	StatJSON(int) []byte
}

type App struct {
	collectors []Collector
}

func NewApp() *App {
	return &App{}
}

func (a *App) Add(collector Collector) {
	a.collectors = append(a.collectors, collector)
}

func (a *App) MapCollectors() map[string]int {
	m := make(map[string]int, len(a.collectors))
	for _, c := range a.collectors {
		m[c.Name()] = 1
	}
	return m
}

func (a *App) Start() {
	for _, c := range a.collectors {
		c.Start()
	}
}

func (a *App) Close() {
	for _, c := range a.collectors {
		c.Stop()
	}
}

// get all stats and return JSON.
func (a *App) StatJSON(period int, collectors ...string) []byte {
	jsons := make([]string, 0, len(a.collectors))
	for _, c := range a.collectors {
		cname := c.Name()
		if len(collectors) > 0 {
			inlist := false
			for _, v := range collectors {
				if v == cname {
					inlist = true
					break
				}
			}
			if !inlist {
				continue
			}
		}
		j := c.StatJSON(period)
		if len(j) > 0 {
			jsons = append(jsons, fmt.Sprintf(`"%s":%s`, cname, j))
		}
	}

	bld := strings.Builder{}
	if len(jsons) > 0 {
		bld.WriteString("{")
		bld.WriteString(strings.Join(jsons, ","))
		bld.WriteString("}")
	}
	return []byte(bld.String())
}
