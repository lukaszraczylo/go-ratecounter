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
			now := time.Now()
			rc.mu.RLock()
			for _, counter := range rc.counters {
				rc.cleanUpOldValues(counter, now)
			}
			rc.mu.RUnlock()
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
	rc.mu.RLock()
	cutoff := now.Add(-rc.interval)
	rc.mu.RUnlock()

	counter.mu.Lock()
	defer counter.mu.Unlock()

	// Filter in place
	newValues := counter.values[:0]
	for _, value := range counter.values {
		if value.timestamp.After(cutoff) {
			newValues = append(newValues, value)
		} else {
			valuesPool.Put(&value)
		}
	}
	counter.values = newValues

	newTicks := counter.ticks[:0]
	for _, tick := range counter.ticks {
		if tick.timestamp.After(cutoff) {
			newTicks = append(newTicks, tick)
		} else {
			ticksPool.Put(&tick)
		}
	}
	counter.ticks = newTicks
}

func (c *Counter) addValue(value int64) {
	t := time.Now()
	v := valuesPool.Get().(*values)
	v.timestamp = t
	v.value = value
	tk := ticksPool.Get().(*ticks)
	tk.timestamp = t

	c.mu.Lock()
	defer c.mu.Unlock()
	c.values = append(c.values, *v)
	c.ticks = append(c.ticks, *tk)
}

func (c *Counter) getValue() int64 {
	sum := int64(0)
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, value := range c.values {
		sum += value.value
	}
	return sum
}
