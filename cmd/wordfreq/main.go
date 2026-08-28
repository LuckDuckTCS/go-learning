package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/LuckDuckTCS/go-learning/internal/wordfreq"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "wordfreq: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("wordfreq", flag.ContinueOnError)
	format := fs.String("format", "text", "output format: text or json")
	top := fs.Int("top", 0, "show only N most frequent words (0 = all)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	var outFormat wordfreq.Formatter
	switch *format {
	case "text":
		outFormat = wordfreq.TextFormatter{}
	case "json":
		outFormat = wordfreq.JSONFormatter{}
	default:
		return fmt.Errorf("ожидался формат text или json, получен %s", *format)
	}

	var r io.Reader = os.Stdin
	rest := fs.Args()
	switch len(rest) {
	case 0:
		// r уже os.Stdin
	case 1:
		f, err := os.Open(rest[0])
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	default:
		return fmt.Errorf("ожидался 1 файл, получено %d", len(rest))
	}

	pairs, err := wordfreq.Count(r)
	if err != nil {
		return err
	}

	if *top > 0 && len(pairs) > *top {
		pairs = pairs[:*top]
	}

	err = outFormat.Format(out, pairs)
	if err != nil {
		return err
	}
	return nil
}
