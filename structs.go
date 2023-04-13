package goratecounter

import (
	"sync"
	"time"
)

var (
	binName = "goratecounter"
)

type ticks struct {
	timestamp time.Time
}

type values struct {
	timestamp time.Time
	value     int64
}

type Counter struct {
	parent *RateCounter
	ticks  []ticks
	values []values
	mu     sync.RWMutex
}

type RateCounter struct {
	counters     map[string]*Counter
	stop         chan bool
	interval     time.Duration
	stopMutex    sync.Mutex
	mu           sync.RWMutex
	customConfig bool
}

type RateCounterConfig struct {
	Interval time.Duration
	MaxSize  int
}
