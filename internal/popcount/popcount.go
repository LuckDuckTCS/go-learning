// Package popcount предоставляет функции подсчёта единичных битов.
package popcount

import "math/bits"

// pc[i] — количество единичных битов в байте i.
var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCountTable возвращает количество единичных битов в x,
// используя предвычисленную таблицу.
func PopCountTable(x uint64) int {
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

func PopCountCycle(x uint64) int {
	count := 0
	for i := range 64 {
		if x&(uint64(1)<<i) != 0 {
			count++
		}
	}
	return count
}

func PopCountKernighan(x uint64) int {
	count := 0
	for x != 0 {
		x = x & (x - 1)
		count++
	}
	return count
}

func PopCountBits(x uint64) int {
	return bits.OnesCount64(x)
}
