package wordfreq

import (
	"bufio"
	"cmp"
	"io"
	"slices"
	"strings"
	"unicode"
)

type Pair struct {
	Word  string
	Count int
}

func Count(r io.Reader) ([]Pair, error) {
	result := []Pair{}

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	m := make(map[string]int)
	// читаем
	for scanner.Scan() {
		word := scanner.Text()

		word = strings.ToLower(word)
		word = strings.TrimFunc(word, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r)
		})
		m[word]++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	// из карты в срез
	for key, value := range m {
		result = append(result, Pair{Word: key, Count: value})
	}

	// Сортируем
	slices.SortFunc(result, func(a, b Pair) int {
		if n := cmp.Compare(b.Count, a.Count); n != 0 {
			return n
		}
		return cmp.Compare(a.Word, b.Word)
	})

	return result, nil
}
