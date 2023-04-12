package goratecounter

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

// NewRateCounter creates a new RateCounter instance supporting multiple counters
func NewRateCounter() *RateCounter {
	rc := &RateCounter{
		interval: 60 * time.Second, // TODO: Make sure to mention in documentation that this is the default value
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
func (rc *RateCounter) WithName(name string) (*RateCounter, error) {
	if rc.counters == nil {
		rc.counters = make(map[string]*Counter, 0)
	}
	if _, ok := rc.counters[name]; ok {
		log.Println(binName + ": name already exists")
		return nil, fmt.Errorf(binName + ": name already exists")
	}
	rc.counters[name] = &Counter{
		active: true,
	}
	rc.restart()
	return rc, nil
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
