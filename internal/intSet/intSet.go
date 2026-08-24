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
			if word&(1<<uint(j)) != 0 {
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
