package goratecounter

import (
	"time"
)

// Interface allowing incrementation of both RateCounter and Counter
type Incrementer interface {
	Incr(v int64)
	Get() int64
	GetTicks() int64
	GetRate() float64
	Average() float64
	Ping()
	GetPingRate() float64
}

// Incr increments the counter by the given value
func (rc *RateCounter) Incr(v int64) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	rc.counters["default"].addValue(v)
}

// Incr increments the counter by the given value
func (c *Counter) Incr(v int64) {
	c.addValue(v)
}

// IncrByName increments the counter by the given value using the given name
func (rc *RateCounter) IncrByName(name string, v int64) {
	rc.mu.RLock()
	counter, ok := rc.counters[name]
	rc.mu.RUnlock()

	if !ok {
		counter, _ = rc.WithName(name)
	}
	counter.addValue(v)
}

// Get returns the current value of the default counter
func (rc *RateCounter) Get() int64 {
	return rc.GetByName("default")
}

// Get returns the current value of the counter
func (c *Counter) Get() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var sum int64
	for _, value := range c.values {
		sum += value.value
	}
	return sum
}

// GetByName returns the current value of the counter with the given name
func (rc *RateCounter) GetByName(name string) int64 {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if counter, ok := rc.counters[name]; ok {
		counter.mu.RLock()
		defer counter.mu.RUnlock()
		var sum int64
		for _, value := range counter.values {
			sum += value.value
		}
		return sum
	}
	return 0
}

// GetTicks returns the current number of events of the default counter
func (rc *RateCounter) GetTicks() int64 {
	return rc.GetTicksByName("default")
}

// GetTicks returns the current number of events of the counter
func (c *Counter) GetTicks() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return int64(len(c.ticks))
}

// GetTicksByName returns the current number of events of the counter with the given name
func (rc *RateCounter) GetTicksByName(name string) int64 {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if counter, ok := rc.counters[name]; ok {
		counter.mu.RLock()
		defer counter.mu.RUnlock()
		return int64(len(counter.ticks))
	}
	return 0
}

// GetRate returns the current rate of the default counter over period of time
func (rc *RateCounter) GetRate() float64 {
	return rc.GetRateByName("default")
}

// GetRate returns the current rate of the counter over period of time
func (c *Counter) GetRate() float64 {
	if c.Get() == 0 {
		return 0
	}
	return float64(c.Get()) / c.parent.interval.Seconds()
}

// GetRateByName returns the current rate of the counter with the given name
func (rc *RateCounter) GetRateByName(name string) float64 {
	total := rc.GetByName(name)
	if total == 0 {
		return 0
	}
	return float64(total) / rc.interval.Seconds()
}

// Average returns the average value of the default counter over period of time
func (rc *RateCounter) Average() float64 {
	return rc.AverageByName("default")
}

// Average returns the average value of the counter over period of time
func (c *Counter) Average() float64 {
	ticks := c.GetTicks()
	if ticks == 0 {
		return 0
	}
	return float64(c.Get()) / float64(ticks)
}

// AverageByName returns the average value of the counter with the given name
func (rc *RateCounter) AverageByName(name string) float64 {
	ticks := rc.GetTicksByName(name)
	if ticks == 0 {
		return 0
	}
	return float64(rc.GetByName(name)) / float64(ticks)
}

// Increase ping counter for the default rate counter
func (rc *RateCounter) Ping() {
	rc.PingByName("default")
}

// Increase ping counter for the counter
func (c *Counter) Ping() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ticks = append(c.ticks, ticks{timestamp: time.Now()})
}

// Increase ping counter for the counter with the given name
func (rc *RateCounter) PingByName(name string) {
	rc.mu.RLock()
	counter, ok := rc.counters[name]
	rc.mu.RUnlock()

	if !ok {
		counter, _ = rc.WithName(name)
	}
	counter.mu.Lock()
	counter.ticks = append(counter.ticks, ticks{timestamp: time.Now()})
	counter.mu.Unlock()
}

// Get ping rate for default counter
func (rc *RateCounter) GetPingRate() float64 {
	return rc.GetPingRateByName("default")
}

// Get ping rate for the counter
func (c *Counter) GetPingRate() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.ticks) == 0 {
		return 0
	}
	return float64(len(c.ticks)) / c.parent.interval.Seconds()
}

// Get ping rate for the counter with the given name
func (rc *RateCounter) GetPingRateByName(name string) float64 {
	ticks := rc.GetTicksByName(name)
	if ticks == 0 {
		return 0
	}
	return float64(ticks) / rc.interval.Seconds()
}
