package counters

import (
	"bytes"
	"io"
	"unicode"
	"unicode/utf8"
)

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // Преобразование int в ByteCounter
	return len(p), nil
}

//Используя идеи из ByteCounter, реализуйте счетчики для слов и строк. Вам пригодится функция bufio.ScanWords.

type WordCounter struct {
	Count  int
	InWord bool
}

func (c *WordCounter) Write(p []byte) (int, error) {
	for i := 0; i < len(p); {
		r, size := utf8.DecodeRune(p[i:])
		if unicode.IsSpace(r) {
			c.InWord = false
		} else {
			if !c.InWord {
				c.Count++
			}
			c.InWord = true
		}
		i += size
	}
	return len(p), nil
}

type LineCounter int

func (c *LineCounter) Write(p []byte) (int, error) {
	*c = LineCounter(bytes.Count(p, []byte("\n")))
	return len(p), nil
}

type countingWriter struct {
	w io.Writer
	n int64
}

func (c *countingWriter) Write(p []byte) (int, error) {
	i, err := c.w.Write(p)
	c.n += int64(i)
	return i, err
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	var c countingWriter
	c.w = w
	return &c, &c.n
}
