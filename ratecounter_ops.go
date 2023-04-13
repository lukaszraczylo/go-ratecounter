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
	rc.counters["default"].addValue(v)
}

// Incr increments the counter by the given value
func (c *Counter) Incr(v int64) {
	c.addValue(v)
}

// IncrByName increments the counter by the given value using the given name
func (rc *RateCounter) IncrByName(name string, v int64) {
	if _, ok := rc.counters[name]; !ok {
		rc.WithName(name)
	}
	rc.counters[name].addValue(v)
}

// Get returns the current value of the default counter
func (rc *RateCounter) Get() int64 {
	return rc.GetByName("default")
}

// Get returns the current value of the counter
func (c *Counter) Get() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	sum := int64(0)
	for _, value := range c.values {
		sum += value.value
	}
	return sum
}

// GetByName returns the current value of the counter with the given name
func (rc *RateCounter) GetByName(name string) int64 {
	if _, ok := rc.counters[name]; !ok {
		rc.WithName(name)
	}
	rc.counters[name].mu.Lock()
	defer rc.counters[name].mu.Unlock()
	sum := int64(0)
	for _, value := range rc.counters[name].values {
		sum += value.value
	}
	return sum
}

// GetTicks returns the current number of events of the default counter
func (rc *RateCounter) GetTicks() int64 {
	return rc.GetTicksByName("default")
}

// GetTicks returns the current number of events of the counter
func (c *Counter) GetTicks() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return int64(len(c.ticks))
}

// GetTicksByName returns the current number of events of the counter with the given name
func (rc *RateCounter) GetTicksByName(name string) int64 {
	if _, ok := rc.counters[name]; !ok {
		rc.WithName(name)
	}
	rc.counters[name].mu.Lock()
	defer rc.counters[name].mu.Unlock()
	return int64(len(rc.counters[name].ticks))
}

// GetRate returns the current rate of the default counter over period of time
func (rc *RateCounter) GetRate() float64 {
	if rc.Get() == 0 {
		return 0
	}
	return float64(rc.Get()) / float64(rc.interval.Seconds())
}

// GetRate returns the current rate of the counter over period of time
func (c *Counter) GetRate() float64 {
	if c.Get() == 0 {
		return 0
	}
	return float64(c.parent.interval.Seconds()) / float64(c.Get())
}

// GetRateByName returns the current rate of the counter with the given name
func (rc *RateCounter) GetRateByName(name string) float64 {
	if rc.GetByName(name) == 0 {
		return 0
	}
	return float64(rc.GetByName(name)) / float64(rc.interval.Seconds())
}

// Average returns the average value of the default counter over period of time
func (rc *RateCounter) Average() float64 {
	if rc.GetTicks() == 0 {
		return 0
	}
	return float64(rc.Get()) / float64(rc.GetTicks())
}

// Average returns the average value of the counter over period of time
func (c *Counter) Average() float64 {
	if c.GetTicks() == 0 {
		return 0
	}
	return float64(c.Get()) / float64(c.GetTicks())
}

// AverageByName returns the average value of the counter with the given name
func (rc *RateCounter) AverageByName(name string) float64 {
	if rc.GetTicksByName(name) == 0 {
		return 0
	}
	return float64(rc.GetByName(name)) / float64(rc.GetTicksByName(name))
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
	if _, ok := rc.counters[name]; !ok {
		rc.WithName(name)
	}
	rc.counters["default"].mu.Lock()
	rc.counters["default"].ticks = append(rc.counters["default"].ticks, ticks{timestamp: time.Now()})
	rc.counters["default"].mu.Unlock()
}

// Get ping rate for default counter
func (rc *RateCounter) GetPingRate() float64 {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	rc.counters["default"].mu.RLock()
	defer rc.counters["default"].mu.RUnlock()
	if len(rc.counters["default"].ticks) == 0 {
		return 0
	}
	return float64(len(rc.counters["default"].ticks)) / float64(rc.interval.Seconds())
}

// Get ping rate for the counter
func (c *Counter) GetPingRate() float64 {
	if len(c.ticks) == 0 {
		return 0
	}
	return float64(c.parent.interval.Seconds()) / float64(len(c.ticks))
}

// Get ping rate for the counter with the given name
func (rc *RateCounter) GetPingRateByName(name string) float64 {
	if rc.GetTicksByName(name) == 0 {
		return 0
	}
	return float64(rc.interval.Seconds()) / float64(len(rc.counters[name].ticks))
}
