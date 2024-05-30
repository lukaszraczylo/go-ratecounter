package goratecounter

import (
	"fmt"
	"time"
)

// NewRateCounter creates a new RateCounter instance supporting multiple counters
func NewRateCounter() *RateCounter {
	rc := &RateCounter{
		interval: 60 * time.Second,
		counters: make(map[string]*Counter),
		stop:     make(chan bool),
	}
	rc.counters["default"] = &Counter{
		ticks:  make([]ticks, 0),
		values: make([]values, 0),
		parent: rc,
	}
	go rc.start()
	return rc
}

// WithConfig sets the configuration of the RateCounter
// If the RateCounter is already active, it will be modified on the fly.
func (rc *RateCounter) WithConfig(config RateCounterConfig) *RateCounter {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.customConfig = true
	rc.interval = config.Interval
	rc.restart()
	return rc
}

// WithName creates a new counter with the given name
// If the counter already exists - it will not be modified and error will be returned.
func (rc *RateCounter) WithName(name string) (*Counter, error) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if _, ok := rc.counters[name]; !ok {
		rc.counters[name] = &Counter{
			ticks:  make([]ticks, 0),
			values: make([]values, 0),
			parent: rc,
		}
		rc.restart()
		return rc.counters[name], nil
	}
	return rc.counters[name], fmt.Errorf("counter with name %s already exists", name)
}
