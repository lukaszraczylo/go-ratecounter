package goratecounter

import (
	"sync/atomic"
	"time"
)

func (rc *RateCounter) start() {
	ticker := time.NewTicker(rc.interval)
	for {
		select {
		case <-ticker.C:
			for id, _ := range rc.counters {
				atomic.SwapInt64(&rc.counters[id].count, 0)
				atomic.SwapInt64(&rc.counters[id].ticks, 0)
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
