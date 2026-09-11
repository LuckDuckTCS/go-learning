package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/LuckDuckTCS/go-learning/internal/loganalyzer"
)

type InvalidFlagError struct {
	Flag  string
	Value string
}

func (e *InvalidFlagError) Error() string {
	return fmt.Sprintf("invalid value %q for -%s", e.Value, e.Flag)
}

func main() {

	if err := run(os.Args[1:]); err != nil {
		var ife *InvalidFlagError
		fmt.Fprintf(os.Stderr, "loganalyzer: %v\n", err)
		if errors.As(err, &ife) {
			fmt.Fprintf(os.Stderr, "hint: run -help to see valid values for -%s\n", ife.Flag)
		}
		os.Exit(1)
	}
}

func run(args []string) (err error) {

	fs := flag.NewFlagSet("loganalyzer", flag.ContinueOnError)
	format := fs.String("format", "text", "output format: text|json|csv")
	top := fs.Int("top", 0, "number of entries in top lists (0 = all)")
	outPath := fs.String("out", "-", "output path (- for stdout)")
	fromStr := fs.String("from", "", "start time (2026-08-01)")
	toStr := fs.String("to", "", "finish time (default: no limit)")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	// валидация top
	if *top < 0 {
		return &InvalidFlagError{Flag: "top", Value: strconv.Itoa(*top)}
	}
	// формат
	var outFormat loganalyzer.Formatter
	switch *format {
	case "text":
		outFormat = loganalyzer.TextFormatter{}
	case "json":
		outFormat = loganalyzer.JSONFormatter{}
	case "csv":
		outFormat = loganalyzer.CSVFormatter{}
	default:
		return fmt.Errorf("unknown format %q (want text|json|csv)", *format)
	}

	// разбор дат
	var fromDate, toDate time.Time

	if *fromStr != "" {
		fromDate, err = time.ParseInLocation(time.DateOnly, *fromStr, time.Local)
		if err != nil {
			return fmt.Errorf("parse -from: %w", err)
		}
	}
	if *toStr != "" {
		toDate, err = time.ParseInLocation(time.DateOnly, *toStr, time.Local)
		if err != nil {
			return fmt.Errorf("parse -to: %w", err)
		}
		// добавим сутки, чтобы поиск был интуитивным 2026-08-01 -- 2026-08-01 станет одними сутками, а не нулём
		toDate = toDate.AddDate(0, 0, 1)
	}

	if !fromDate.IsZero() && !toDate.IsZero() && fromDate.After(toDate) {
		return fmt.Errorf("-from %s is after -to %s", *fromStr, *toStr)
	}

	a := loganalyzer.NewAggregator()
	rest := fs.Args()

	if len(rest) == 0 {
		err = loganalyzer.Process(os.Stdin, a, fromDate, toDate) // stdin напрямую
		if err != nil {
			return err
		}
	} else {
		// определить: файл или директория
		info, statErr := os.Stat(rest[0])
		if statErr != nil {
			return fmt.Errorf("stat input: %w", statErr)
		}
		var paths []string
		if info.IsDir() {
			paths, err = loganalyzer.ScanDir(rest[0])
			if err != nil {
				return fmt.Errorf("scan dir: %w", err)
			}
		} else {
			paths = []string{rest[0]}
		}

		for _, p := range paths {
			err := loganalyzer.ProcessFile(p, a, fromDate, toDate) // обёртка
			if err != nil {
				return err
			}
		}
	}

	// формируем writer

	var w io.Writer

	switch *outPath {
	case "-":
		w = os.Stdout
	default:
		if err := os.MkdirAll(filepath.Dir(*outPath), 0o755); err != nil {
			return fmt.Errorf("create output dir: %w", err)
		}

		f, createErr := os.Create(*outPath)
		if createErr != nil {
			return fmt.Errorf("create %s: %w", *outPath, createErr)
		}

		defer func() {
			if closeErr := f.Close(); closeErr != nil && err == nil {
				err = fmt.Errorf("close output: %w", closeErr)
			}
		}()
		w = f
	}

	report := a.Report(*top)
	return outFormat.Format(w, report)

}
