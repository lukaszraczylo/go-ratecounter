package goratecounter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *Tests) TestRateCounter_Incr() {
	type fields struct {
		counters     map[string]*Counter
		interval     time.Duration
		active       bool
		customConfig bool
	}
	type args struct {
		v int64
	}
	tests := []struct {
		want   *RateCounter
		name   string
		fields fields
		args   args
	}{
		{
			name: "IncrTest - 1337",
			fields: fields{
				active:   true,
				counters: map[string]*Counter{},
			},
			args: args{
				v: 1337,
			},
			want: &RateCounter{
				interval: 60 * time.Second,
				counters: map[string]*Counter{
					"default": &Counter{
						values: []values{
							{
								value:     1337,
								timestamp: time.Now(),
							},
						},
						ticks: []ticks{
							{
								timestamp: time.Now(),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			test_rc.Incr(tt.args.v)
			assert.Equal(t, test_rc.Get(), tt.want.counters["default"].getValue())
			assert.Equal(t, test_rc.GetTicks(), int64(len(tt.want.counters["default"].ticks)))
		})
	}
}

func (suite *Tests) TestRateCounter_IncrByName() {
	type fields struct {
	}
	type args struct {
		name     string
		v        int64
		withName bool
	}
	tests := []struct {
		want *RateCounter
		name string
		args args
	}{
		{
			name: "IncrByNameTest - 1337",
			args: args{
				withName: false,
			},
			want: &RateCounter{
				interval: 60 * time.Second,
				counters: map[string]*Counter{
					"default": &Counter{
						values: []values{
							{
								timestamp: time.Now(),
								value:     1337,
							},
						},
						ticks: []ticks{
							{
								timestamp: time.Now(),
							},
						},
					},
				},
			},
		},
		{
			name: "IncrByNameTestWithSetup - 1337",
			args: args{
				withName: true,
				name:     "test",
			},
			want: &RateCounter{
				interval: 60 * time.Second,
				counters: map[string]*Counter{
					"default": &Counter{
						values: []values{
							{
								timestamp: time.Now(),
								value:     1,
							},
							{
								timestamp: time.Now(),
								value:     1337,
							},
						},
						ticks: []ticks{
							{
								timestamp: time.Now(),
							}, {
								timestamp: time.Now(),
							},
						},
					},
					"test": &Counter{
						values: []values{
							{
								timestamp: time.Now(),
								value:     8,
							},
						},
						ticks: []ticks{
							{
								timestamp: time.Now(),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			tc := &Counter{}
			if !tt.args.withName {
			} else {
				tc, _ = test_rc.WithName(tt.args.name)
			}

			test_rc.IncrByName(tt.args.name, tt.args.v)
			assert.Equal(t, test_rc.GetByName(tt.args.name), tc.getValue())
		})
	}
}

func (suite *Tests) TestRateCounter_IncrBy() {
	type fields struct {
	}
	type args struct {
		name    string
		v       int64
		rate    float64
		average float64
	}
	tests := []struct {
		want *RateCounter
		name string
		args args
	}{
		{
			name: "IncrByTest - 1337",
			args: args{
				name:    "t123",
				v:       1337,
				rate:    22.283333333333335, // ( 1337 / 60 )
				average: 1337,
			},
			want: &RateCounter{
				interval: 60 * time.Second,
				counters: map[string]*Counter{
					"t123": &Counter{
						values: []values{
							{
								value:     1337,
								timestamp: time.Now(),
							},
						},
						ticks: []ticks{
							{
								timestamp: time.Now(),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			test_rc.Incr(tt.args.v)
			assert.Equal(t, test_rc.Get(), tt.args.v, "Value does not match")
			assert.Equal(t, test_rc.GetRate(), tt.args.rate, "Rate does not match")
			assert.Equal(t, test_rc.Average(), tt.args.average, "Average does not match")

			test_rc.IncrByName(tt.args.name, tt.args.v)
			assert.Equal(t, test_rc.GetByName(tt.args.name), tt.args.v, "Value does not match")
			assert.Equal(t, test_rc.GetRateByName(tt.args.name), tt.args.rate, "Rate does not match")
			assert.Equal(t, test_rc.AverageByName(tt.args.name), tt.args.average, "Average does not match")
		})
	}
}

func (suite *Tests) TestRateCounter_WithExpiry() {
	type fields struct {
	}
	type args struct {
		name     string
		v        int64
		rate     float64
		average  float64
		interval time.Duration
	}
	tests := []struct {
		want *RateCounter
		name string
		args args
	}{
		{
			name: "IncrByTest - 1337",
			args: args{
				name:     "t123",
				v:        1337,
				rate:     668.5, // ( 1337 / interval )
				average:  1337,
				interval: 2,
			},
			want: &RateCounter{
				interval: 60 * time.Second,
				counters: map[string]*Counter{
					"t123": &Counter{
						values: []values{
							{
								value:     1337,
								timestamp: time.Now(),
							},
						},
						ticks: []ticks{
							{
								timestamp: time.Now(),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			test_rc.WithConfig(RateCounterConfig{
				Interval: time.Duration(tt.args.interval) * time.Second,
			})
			test_rc.Incr(tt.args.v)
			assert.Equal(t, tt.args.v, test_rc.Get(), "Value does not match")
			assert.Equal(t, tt.args.rate, test_rc.GetRate(), "Rate does not match")
			assert.Equal(t, tt.args.average, test_rc.Average(), "Average does not match")

			time.Sleep(time.Duration(tt.args.interval+5) * time.Second)

			assert.Equal(t, int64(0), test_rc.Get(), "Value does not match")
			assert.Equal(t, float64(0), test_rc.GetRate(), "Rate does not match")
			assert.Equal(t, float64(0), test_rc.Average(), "Average does not match")
		})
	}
}

func (suite *Tests) TestRateCounter_PingWithExpiry() {
	type fields struct {
	}
	type args struct {
		name     string
		rate     float64
		interval time.Duration
	}
	tests := []struct {
		want *RateCounter
		name string
		args args
	}{
		{
			name: "IncrByTest - 1337",
			args: args{
				name:     "test-rate",
				rate:     0.5,
				interval: 2,
			},
			want: &RateCounter{
				interval: 60 * time.Second,
				counters: map[string]*Counter{
					"t123": &Counter{
						values: []values{
							{
								value:     1337,
								timestamp: time.Now(),
							},
						},
						ticks: []ticks{
							{
								timestamp: time.Now(),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			test_rc.WithConfig(RateCounterConfig{
				Interval: time.Duration(tt.args.interval) * time.Second,
			})
			test_rc.Ping()
			assert.Equal(t, tt.args.rate, test_rc.GetPingRate(), "Average does not match")
			time.Sleep(time.Duration(tt.args.interval*2) * time.Second)
			assert.Equal(t, float64(0), test_rc.GetPingRate(), "Rate does not match")
		})
	}
}
