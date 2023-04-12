package goratecounter

import (
	"sync"
	"time"
)

var (
	binName = "goratecounter"
)

type Counter struct {
	active bool
	ticks  int64
	count  int64
}

type RateCounter struct {
	counters     map[string]*Counter
	stop         chan bool
	interval     time.Duration
	stopMutex    sync.Mutex
	mu           sync.Mutex
	customConfig bool
}

type RateCounterConfig struct {
	Interval time.Duration
}
