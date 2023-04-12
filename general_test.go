package goratecounter

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Run the test concurrently 100 times adding random values to the counter and making sure the end value is correct
func (suite *Tests) TestRateCounter_StressIncr() {
	tests := []struct {
		name       string
		cycles     int
		routines   int
		wantPasses int64
	}{
		{
			name:       "StressIncrTest",
			cycles:     1000,
			routines:   1000,
			wantPasses: 1000 * 1000,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			currentCount := int64(0)

			for i := 0; i < tt.routines; i++ {
				wg.Add(1)
				go func(w *sync.WaitGroup) {
					c := int64(0)
					for j := 0; j < tt.cycles; j++ {
						w.Add(1)
						c = rand.Int63n(100)
						atomic.AddInt64(&currentCount, c)
						test_rc.Incr(c)
						w.Done()
					}
					w.Done()
				}(&wg)
			}

			wg.Wait()
			rcRes := test_rc.Get()
			assert.Equal(t, rcRes, currentCount, "Mismatched counter values")
			assert.Equal(t, tt.wantPasses, test_rc.GetTicks(), "Mismatched ticks")
		})
	}
}
