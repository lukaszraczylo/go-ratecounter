package goratecounter

import (
	"time"
)

func (rc *RateCounter) start() {
	ticker := time.NewTicker(rc.interval)
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
	cutoff := now.Add(-rc.interval)
	newValues := make([]values, 0, len(counter.values))
	newTicks := make([]ticks, 0, len(counter.ticks))
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
}

func (c *Counter) addValue(value int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	t := time.Now()
	c.values = append(c.values, values{value: value, timestamp: t})
	c.ticks = append(c.ticks, ticks{timestamp: t})
}

func (c *Counter) getValue() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	sum := int64(0)
	for _, value := range c.values {
		sum += value.value
	}
	return sum
}
