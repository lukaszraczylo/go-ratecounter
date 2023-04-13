package goratecounter

import (
	"time"
)

func (rc *RateCounter) start() {
	rc.mu.RLock()
	ticker := time.NewTicker(rc.interval)
	rc.mu.RUnlock()
	for {
		select {
		case <-ticker.C:
			for _, counter := range rc.counters {
				now := time.Now()
				rc.cleanUpOldValues(counter, now)
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

func (rc *RateCounter) cleanUpOldValues(counter *Counter, now time.Time) {
	newValues := make([]values, 0)
	newTicks := make([]ticks, 0)

	rc.mu.RLock()
	cutoff := now.Add(-rc.interval)
	rc.mu.RUnlock()
	counter.mu.Lock()
	for _, value := range counter.values {
		if value.timestamp.After(cutoff) {
			newValues = append(newValues, value)
		}
	}
	for _, tick := range counter.ticks {
		if tick.timestamp.After(cutoff) {
			newTicks = append(newTicks, tick)
		}
	}
	counter.values = newValues
	counter.ticks = newTicks
	counter.mu.Unlock()
}

func (c *Counter) addValue(value int64) {
	t := time.Now()
	c.mu.Lock()
	c.values = append(c.values, values{value: value, timestamp: t})
	c.ticks = append(c.ticks, ticks{timestamp: t})
	c.mu.Unlock()
}

func (c *Counter) getValue() int64 {
	sum := int64(0)
	c.mu.RLock()
	for _, value := range c.values {
		sum += value.value
	}
	c.mu.RUnlock()
	return sum
}
