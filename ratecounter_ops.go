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

func (rc *RateCounter) Incr(v int64) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	rc.counters["default"].addValue(v)
}

func (c *Counter) Incr(v int64) {
	c.addValue(v)
}

func (rc *RateCounter) IncrByName(name string, v int64) {
	rc.mu.RLock()
	counter, ok := rc.counters[name]
	rc.mu.RUnlock()

	if !ok {
		counter, _ = rc.WithName(name)
	}
	counter.addValue(v)
}

func (rc *RateCounter) Get() int64 {
	return rc.GetByName("default")
}

func (c *Counter) Get() int64 {
	return c.getValue()
}

func (rc *RateCounter) GetByName(name string) int64 {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if counter, ok := rc.counters[name]; ok {
		return counter.getValue()
	}
	return 0
}

func (rc *RateCounter) GetTicks() int64 {
	return rc.GetTicksByName("default")
}

func (c *Counter) GetTicks() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return int64(len(c.ticks))
}

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

func (rc *RateCounter) GetRate() float64 {
	return rc.GetRateByName("default")
}

func (c *Counter) GetRate() float64 {
	if c.getValue() == 0 {
		return 0
	}
	return float64(c.getValue()) / c.parent.interval.Seconds()
}

func (rc *RateCounter) GetRateByName(name string) float64 {
	total := rc.GetByName(name)
	if total == 0 {
		return 0
	}
	return float64(total) / rc.interval.Seconds()
}

func (rc *RateCounter) Average() float64 {
	return rc.AverageByName("default")
}

func (c *Counter) Average() float64 {
	ticks := c.GetTicks()
	if ticks == 0 {
		return 0
	}
	return float64(c.getValue()) / float64(ticks)
}

func (rc *RateCounter) AverageByName(name string) float64 {
	ticks := rc.GetTicksByName(name)
	if ticks == 0 {
		return 0
	}
	return float64(rc.GetByName(name)) / float64(ticks)
}

func (rc *RateCounter) Ping() {
	rc.PingByName("default")
}

func (c *Counter) Ping() {
	t := ticksPool.Get().(*ticks)
	t.timestamp = time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ticks = append(c.ticks, *t)
}

func (rc *RateCounter) PingByName(name string) {
	rc.mu.RLock()
	counter, ok := rc.counters[name]
	rc.mu.RUnlock()

	if !ok {
		counter, _ = rc.WithName(name)
	}
	t := ticksPool.Get().(*ticks)
	t.timestamp = time.Now()
	counter.mu.Lock()
	counter.ticks = append(counter.ticks, *t)
	counter.mu.Unlock()
}

func (rc *RateCounter) GetPingRate() float64 {
	return rc.GetPingRateByName("default")
}

func (c *Counter) GetPingRate() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.ticks) == 0 {
		return 0
	}
	return float64(len(c.ticks)) / c.parent.interval.Seconds()
}

func (rc *RateCounter) GetPingRateByName(name string) float64 {
	ticks := rc.GetTicksByName(name)
	if ticks == 0 {
		return 0
	}
	return float64(ticks) / rc.interval.Seconds()
}
