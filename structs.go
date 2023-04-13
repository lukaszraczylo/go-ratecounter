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
	value     int64
	timestamp time.Time
}

type Counter struct {
	ticks      []ticks
	values     []values
	parent     *RateCounter
	mu         sync.Mutex
	valuesPool sync.Pool
	ticksPool  sync.Pool
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
	MaxSize  int
}
