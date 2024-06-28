package goratecounter

import (
	"fmt"
	"sync"
	"time"
)

var (
	ticksPool = sync.Pool{
		New: func() interface{} {
			return &ticks{}
		},
	}

	valuesPool = sync.Pool{
		New: func() interface{} {
			return &values{}
		},
	}
)

func NewRateCounter() *RateCounter {
	rc := &RateCounter{
		interval: 60 * time.Second,
		counters: make(map[string]*Counter),
		stop:     make(chan bool),
	}
	rc.counters["default"] = &Counter{
		ticks:  make([]ticks, 0, 1000),
		values: make([]values, 0, 1000),
		parent: rc,
	}
	go rc.start()
	return rc
}

func (rc *RateCounter) WithConfig(config RateCounterConfig) *RateCounter {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.customConfig = true
	rc.interval = config.Interval
	rc.restart()
	return rc
}

func (rc *RateCounter) WithName(name string) (*Counter, error) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if _, ok := rc.counters[name]; !ok {
		rc.counters[name] = &Counter{
			ticks:  make([]ticks, 0, 1000),
			values: make([]values, 0, 1000),
			parent: rc,
		}
		rc.restart()
		return rc.counters[name], nil
	}
	return rc.counters[name], fmt.Errorf("counter with name %s already exists", name)
}
