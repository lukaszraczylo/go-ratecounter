package goratecounter

import "sync/atomic"

// Interface allowing incrementation of both RateCounter and Counter
type Incrementer interface {
	Incr(v int64)
	Get() int64
	GetTicks() int64
	GetRate() float64
	Average() float64
}

// Incr increments the counter by the given value
func (rc *RateCounter) Incr(v int64) {
	atomic.AddInt64(&rc.counters["default"].count, int64(v))
	atomic.AddInt64(&rc.counters["default"].ticks, 1)
}

// Incr increments the counter by the given value
func (c *Counter) Incr(v int64) {
	atomic.AddInt64(&c.count, int64(v))
	atomic.AddInt64(&c.ticks, 1)
}

// IncrByName increments the counter by the given value using the given name
func (rc *RateCounter) IncrByName(name string, v int64) {
	if _, ok := rc.counters[name]; !ok {
		rc.WithName(name)
	}
	atomic.AddInt64(&rc.counters[name].ticks, 1)
	atomic.AddInt64(&rc.counters[name].count, int64(v))
}

// Get returns the current value of the default counter
func (rc *RateCounter) Get() int64 {
	return atomic.LoadInt64(&rc.counters["default"].count)
}

// Get returns the current value of the counter
func (c *Counter) Get() int64 {
	return atomic.LoadInt64(&c.count)
}

// GetByName returns the current value of the counter with the given name
func (rc *RateCounter) GetByName(name string) int64 {
	if _, ok := rc.counters[name]; !ok {
		rc.WithName(name)
	}
	return atomic.LoadInt64(&rc.counters[name].count)
}

// GetTicks returns the current number of events of the default counter
func (rc *RateCounter) GetTicks() int64 {
	return atomic.LoadInt64(&rc.counters["default"].ticks)
}

// GetTicks returns the current number of events of the counter
func (c *Counter) GetTicks() int64 {
	return atomic.LoadInt64(&c.ticks)
}

// GetTicksByName returns the current number of events of the counter with the given name
func (rc *RateCounter) GetTicksByName(name string) int64 {
	if _, ok := rc.counters[name]; !ok {
		rc.WithName(name)
	}
	return atomic.LoadInt64(&rc.counters[name].ticks)
}

// GetRate returns the current rate of the default counter over period of time
func (rc *RateCounter) GetRate() float64 {
	return float64(rc.Get()) / float64(rc.interval.Seconds())
}

// GetRate returns the current rate of the counter over period of time
func (c *Counter) GetRate() float64 {
	return float64(c.Get()) / float64(c.parent.interval.Seconds())
}

// GetRateByName returns the current rate of the counter with the given name
func (rc *RateCounter) GetRateByName(name string) float64 {
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
