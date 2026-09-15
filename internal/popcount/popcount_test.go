package popcount

import (
	"math/bits"
	"testing"
)

func TestPopCount(t *testing.T) {
	values := []uint64{
		0,
		1,
		0xFFFFFFFFFFFFFFFF,
		0x8000000000000000, // только старший бит
		0x1234567890ABCDEF,
		0xAAAAAAAAAAAAAAAA, // через один
	}

	impls := map[string]func(uint64) int{
		"table":     PopCountTable,
		"cycle":     PopCountCycle,
		"kernighan": PopCountKernighan,
	}

	for _, v := range values {
		want := bits.OnesCount64(v)
		for name, f := range impls {
			if got := f(v); got != want {
				t.Errorf("%s(%#016x) = %d, want %d", name, v, got, want)
			}
		}
	}
}

var benchInput uint64 = 0x1234567890ABCDEF

func BenchmarkPopCountTable(b *testing.B) {
	for b.Loop() {
		PopCountTable(benchInput)
	}
}

func BenchmarkPopCountCycle(b *testing.B) {
	for b.Loop() {
		PopCountCycle(benchInput)
	}
}

func BenchmarkPopCountKernighan(b *testing.B) {
	cases := map[string]uint64{
		"zero":   0,
		"sparse": 0x0000000000000001,
		"half":   0x1234567890ABCDEF,
		"full":   0xFFFFFFFFFFFFFFFF,
	}
	for name, v := range cases {
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				PopCountKernighan(v)
			}
		})
	}
}

func BenchmarkPopCountBits(b *testing.B) {
	for b.Loop() {
		bits.OnesCount64(benchInput)
	}
}
