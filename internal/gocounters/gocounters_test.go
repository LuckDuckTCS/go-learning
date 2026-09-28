package gocounters

import "testing"

func BenchmarkMutexCounterInc(b *testing.B) {
	c := MutexCounter{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Inc()
		}
	})
}
func BenchmarkMutexCounterValue(b *testing.B) {
	c := MutexCounter{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Value()
		}
	})
}

func BenchmarkChanCounterInc(b *testing.B) {
	c := NewChanCounter()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Inc()
		}
	})
}

func BenchmarkChanCounterValue(b *testing.B) {
	c := NewChanCounter()
	defer c.Close()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Value()
		}
	})
}
