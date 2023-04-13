package goratecounter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Tests struct {
	suite.Suite
}

var test_rc *RateCounter

func (suite *Tests) SetupTest() {
	test_rc = NewRateCounter()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Tests))
}

func (suite *Tests) TestNewRateCounter() {
	tests := []struct {
		want *RateCounter
		name string
	}{
		{
			name: "NewRateCounterTest",
			want: &RateCounter{
				interval: 60 * time.Second,
				counters: map[string]*Counter{
					"default": &Counter{},
				},
			},
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			assert.Equal(t, test_rc.interval, tt.want.interval)
		})
	}
}

func (suite *Tests) TestRateCounter_WithConfig() {
	type args struct {
		config RateCounterConfig
	}
	tests := []struct {
		want *RateCounter
		name string
		args args
	}{
		{
			name: "WithConfigTest",
			args: args{
				config: RateCounterConfig{
					Interval: 1 * time.Second,
				},
			},
			want: &RateCounter{
				interval: 1 * time.Second,
				counters: map[string]*Counter{
					"default": &Counter{},
				},
			},
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			test_rc.WithConfig(tt.args.config)
			assert.Equal(t, test_rc.interval, tt.want.interval)
		})
	}
}

func (suite *Tests) TestRateCounter_WithName() {
	type fields struct {
		name string
	}
	type args struct {
		duplicate bool
	}
	tests := []struct {
		want   *RateCounter
		name   string
		fields fields
		args   args
	}{
		{
			name: "WithNameTest",
			fields: fields{
				name: "test",
			},
			want: &RateCounter{
				counters: map[string]*Counter{
					"test": &Counter{
						ticks: nil,
					},
				},
			},
		},
		{
			name: "Duplicated name",
			fields: fields{
				name: "test",
			},
			args: args{
				duplicate: true,
			},
			want: &RateCounter{
				counters: map[string]*Counter{
					"test": &Counter{
						ticks: nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			test_rc.counters = map[string]*Counter{}
			_, err := test_rc.WithName(tt.fields.name)
			assert.NoError(t, err)
			if tt.args.duplicate {
				_, err = test_rc.WithName(tt.fields.name)
				assert.Error(t, err)
			}
			// assert equal except of field "parent"
			// assert.Equal(t, test_rc.counters[tt.fields.name].ticks, tt.want.counters[tt.fields.name].ticks)
			assert.Equal(t, len(test_rc.counters[tt.fields.name].values), len(tt.want.counters[tt.fields.name].values))
		})
	}
}
