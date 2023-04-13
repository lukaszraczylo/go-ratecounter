package goratecounter

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

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

func (suite *Tests) TestRateCounter_Daisychain() {
	type args struct {
		name      string
		increment int64
	}
	tests := []struct {
		want interface{}
		name string
		args args
	}{
		{
			name: "DaisyChain",
			args: args{
				name:      "test",
				increment: 7,
			},
			want: &Counter{
				values: []values{
					{timestamp: time.Now(), value: 7},
				},
				ticks: []ticks{
					{timestamp: time.Now()},
				},
			},
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			r, err := test_rc.WithName(tt.args.name)
			assert.NoError(t, err)
			r.Incr(tt.args.increment)

			// Test daisy chained counter
			assert.Equal(t, r.Get(), tt.want.(*Counter).getValue())
			assert.Equal(t, r.GetTicks(), int64(len(tt.want.(*Counter).ticks)))
			assert.Equal(t, r.Average(), float64(tt.want.(*Counter).getValue()/int64(len(tt.want.(*Counter).ticks))))

			// Test base counter
			assert.Equal(t, test_rc.Get(), int64(0))
			assert.Equal(t, test_rc.GetTicks(), int64(0))
			assert.Equal(t, test_rc.Average(), float64(0))
			assert.Equal(t, test_rc.GetRate(), float64(0))
		})
	}
}
