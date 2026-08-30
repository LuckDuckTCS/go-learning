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

func CharCount(r io.Reader, w io.Writer) error {
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
			return err
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
	fmt.Fprintf(w, "rune\tcount\n")
	for c, n := range counts {
		fmt.Fprintf(w, "%q\t%d\n", c, n)
	}
	fmt.Fprintf(w, "Type\tcount\n")
	for c, n := range typeCounts {
		fmt.Fprintf(w, "%s\t%d\n", c, n)
	}
	fmt.Fprintf(w, "\nlen\tcount\n")
	for i, n := range utfLen {
		if i > 0 {
			fmt.Fprintf(w, "%d\t%d\n", i, n)
		}
	}

	if invalid > 0 {
		fmt.Fprintf(w, "\n%d неверных символов UTF-8\n", invalid)
	}
	return nil
}
