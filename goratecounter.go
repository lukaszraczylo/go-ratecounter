package goratecounter

import (
	"fmt"
	"sync/atomic"
	"time"
)

// NewRateCounter creates a new RateCounter instance supporting multiple counters
func NewRateCounter() *RateCounter {
	rc := &RateCounter{
		interval: 60 * time.Second,
		counters: make(map[string]*Counter),
	}
	rc.counters["default"] = &Counter{
		active: true,
	}
	rc.stop = make(chan bool)
	go rc.start()
	return rc
}

// WithConfig sets the configuration of the RateCounter
// If the RateCounter is already active, it will be modified on the fly.
func (rc *RateCounter) WithConfig(config RateCounterConfig) *RateCounter {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	// for each field of the config struct, set the corresponding field of the RateCounter struct
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
	if rc.counters == nil {
		rc.counters = make(map[string]*Counter, 0)
	}
	if _, ok := rc.counters[name]; !ok {
		rc.counters[name] = &Counter{
			active: true,
			parent: rc,
		}
		rc.restart()
	} else {
		return rc.counters[name], fmt.Errorf("counter with name %s already exists", name)
	}
	return rc.counters[name], nil
}

// Channel and ticker related functions

func (rc *RateCounter) start() {
	ticker := time.NewTicker(rc.interval)
	for {
		select {
		case <-ticker.C:
			for id, _ := range rc.counters {
				atomic.SwapInt64(&rc.counters[id].count, 0)
				atomic.SwapInt64(&rc.counters[id].ticks, 0)
			}
		case <-rc.getStopChan():
			ticker.Stop()
			return
		}
	}
}

func (rc *RateCounter) getCounter(name string) *Counter {
	if c, ok := rc.counters[name]; ok {
		return c
	} else {
		rc.mu.Lock()
		rc.counters[name] = &Counter{
			active: true,
			parent: rc,
		}
		rc.mu.Unlock()
		return rc.counters[name]
	}
}

func (rc *RateCounter) getStopChan() chan bool {
	rc.stopMutex.Lock()
	defer rc.stopMutex.Unlock()
	if rc.stop == nil {
		rc.stop = make(chan bool)
	}
	return rc.stop
}

func (rc *RateCounter) stopTicker() {
	rc.stopMutex.Lock()
	defer rc.stopMutex.Unlock()
	if rc.stop != nil {
		close(rc.stop)
		rc.stop = nil
	}
}

func (rc *RateCounter) restart() {
	rc.stopTicker()
	go rc.start()
}
