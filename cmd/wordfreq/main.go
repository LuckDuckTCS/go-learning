package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/LuckDuckTCS/go-learning/internal/wordfreq"
)

type InvalidFlagError struct {
	Flag  string
	Value string
}

func (e *InvalidFlagError) Error() string {
	return fmt.Sprintf("invalid value %q for -%s", e.Value, e.Flag)
}

func main() {

	if err := run(os.Args[1:], os.Stdout); err != nil {
		var ife *InvalidFlagError
		fmt.Fprintf(os.Stderr, "wordfreq: %v\n", err)
		if errors.As(err, &ife) {
			fmt.Fprintf(os.Stderr, "hint: run -help to see valid values for -%s\n", ife.Flag)
		}
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("wordfreq", flag.ContinueOnError)
	format := fs.String("format", "text", "output format: text or json")
	top := fs.Int("top", 0, "show only N most frequent words (0 = all)")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	if *top < 0 {
		return &InvalidFlagError{Flag: "top", Value: strconv.Itoa(*top)}
	}

	var outFormat wordfreq.Formatter
	switch *format {
	case "text":
		outFormat = wordfreq.TextFormatter{}
	case "json":
		outFormat = wordfreq.JSONFormatter{}
	default:
		return fmt.Errorf("unknown format %q (want text or json)", *format)
	}

	var r io.Reader = os.Stdin
	rest := fs.Args()
	switch len(rest) {
	case 0:
		// r уже os.Stdin
	case 1:
		f, err := os.Open(rest[0])
		if err != nil {
			return fmt.Errorf("open input: %w", err)
		}
		defer f.Close()
		r = f
	default:
		return fmt.Errorf("too many files: got %d, want 1", len(rest))
	}

	pairs, err := wordfreq.Count(r)
	if err != nil {
		if errors.Is(err, wordfreq.ErrEmptyInput) {
			fmt.Fprintln(os.Stderr, "wordfreq: no words in input")
			return nil
		}
		return fmt.Errorf("count words: %w", err)
	}

	if *top > 0 && len(pairs) > *top {
		pairs = pairs[:*top]
	}

	err = outFormat.Format(out, pairs)
	if err != nil {
		return fmt.Errorf("format output: %w", err)
	}
	return nil
}
