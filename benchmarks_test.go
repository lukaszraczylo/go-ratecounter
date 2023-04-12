package goratecounter

import (
	"math/rand"
	"testing"
)

func BenchmarkRateCounter_Default(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := rand.Int63n(100)
		rc := NewRateCounter()
		rc.Incr(c)
	}
}

func BenchmarkRateCounter_WithName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := rand.Int63n(100)
		rc := NewRateCounter()
		rc.IncrByName("test", c)
	}
}

func BenchmarkRateCounter_WithMultipleNames(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := rand.Int63n(100)
		rc := NewRateCounter()
		rc.IncrByName("test", c)
		rc.IncrByName("test2", c)
	}
}

func BenchmarkRateCounter_WithMultipleNamesAndDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := rand.Int63n(100)
		rc := NewRateCounter()
		rc.Incr(c)
		rc.IncrByName("test", c)
		rc.IncrByName("test2", c)
	}
}
