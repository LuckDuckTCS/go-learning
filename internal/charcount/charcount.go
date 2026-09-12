// Package charcount предоставляет функцию, которая считает символы из переданного файла
package charcount

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

type Counts struct {
	Runes   map[rune]int
	Types   map[string]int
	UTFLen  [utf8.UTFMax + 1]int
	Invalid int
}

func CountChars(r io.Reader) (Counts, error) {
	counts := make(map[rune]int)       // кол-во символов unicode
	typeCounts := make(map[string]int) // кол-во типов символов (буква, цифра ...)
	var utfLen [utf8.UTFMax + 1]int    //кол-во длин кодировок utf8
	invalid := 0

	in := bufio.NewReader(r)

	for {
		r, n, err := in.ReadRune()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return Counts{}, err
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utfLen[n]++
		switch {
		case unicode.IsNumber(r):
			typeCounts["Numbers"]++
		case unicode.IsLetter(r):
			typeCounts["Letters"]++
		case unicode.IsPunct(r):
			typeCounts["Puncts"]++
		case unicode.IsSymbol(r):
			typeCounts["Symbols"]++
		case unicode.IsSpace(r):
			typeCounts["Spaces"]++
		default:
			typeCounts["Other"]++
		}
	}
	return Counts{Runes: counts, Types: typeCounts, UTFLen: utfLen, Invalid: invalid}, nil
}

func PrintCount(counts Counts, w io.Writer) error {

	fmt.Fprintf(w, "rune\tcount\n")
	for c, n := range counts.Runes {
		fmt.Fprintf(w, "%q\t%d\n", c, n)
	}
	fmt.Fprintf(w, "Type\tcount\n")
	for c, n := range counts.Types {
		fmt.Fprintf(w, "%s\t%d\n", c, n)
	}
	fmt.Fprintf(w, "\nlen\tcount\n")
	for i, n := range counts.UTFLen {
		if i > 0 {
			fmt.Fprintf(w, "%d\t%d\n", i, n)
		}
	}

	if counts.Invalid > 0 {
		fmt.Fprintf(w, "\n%d неверных символов UTF-8\n", counts.Invalid)
	}
	return nil
}
