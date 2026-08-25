package intset

import (
	"bytes"
	"fmt"
)

type IntSet struct {
	words []uint64
}

func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith делает множество s равным объединению множеств s и t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// String возвращает множество как строку вида "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(uint64(1)<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// Возвращает количество элементов
func (s *IntSet) Len() int {
	var result int
	for _, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				result++
			}
		}
	}
	return result
}

// Удаляет x из множества
func (s *IntSet) Remove(x int) {
	word, bit := x/64, uint(x%64)
	if word < len(s.words) {
		s.words[word] &^= (1 << bit)
	}
}

// Удаляет все элементы множества
func (s *IntSet) Clear() {
	for i := range s.words {
		s.words[i] = 0
	}
}

// Возвращает копию множества
func (s *IntSet) Copy() *IntSet {
	result := &IntSet{words: make([]uint64, len(s.words))}
	copy(result.words, s.words)
	return result
}

// Добавляет группу элементов
func (s *IntSet) AddAll(nums ...int) {
	for _, num := range nums {
		s.Add(num)
	}
}

// упр 6.3 методы IntersectWith, DifferenceWith и SymmetricDifference

// IntersectWith делает множество s равным пересечению множеств s и t.
func (s *IntSet) IntersectWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &= tword
		}
	}
	if len(s.words) > len(t.words) {
		clear(s.words[len(t.words):])
		//s.words = s.words[:len(t.words)]
	}
}

// DifferenceWith делает множество s равным разности множеств s и t.
func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &^= tword
		}
	}
}

// SymmetricDifference делает множество s равным симметричной разности множеств s и t.
func (s *IntSet) SymmetricDifference(t *IntSet) {
	sCopy := s.Copy()
	tCopy := t.Copy()
	sCopy.DifferenceWith(t)
	tCopy.DifferenceWith(s)
	sCopy.UnionWith(tCopy)
	*s = *sCopy
}

// Elems возвращает срез, содержащий элементы множества и годящийся для итерирования с использованием цикла по диапазону range.
func (s *IntSet) Elems() []int {
	result := []int{} //можно result = make([]int, len(s.words)/64)
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(uint64(1)<<uint(j)) != 0 {
				result = append(result, 64*i+j)
			}
		}
	}
	return result
}
