package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"

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

	fs := flag.NewFlagSet("loggen", flag.ContinueOnError)
	path := fs.String("in", "testdata/big.log", "path to log file")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	f, err := os.Open(*path)
	if err != nil {
		return fmt.Errorf("open file %s: %w", *path, err)
	}

	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close file: %w", closeErr)
		}
	}()

	r := bufio.NewReader(f)

	n, errScan := loganalyzer.ScanLines(r)
	if errScan != nil {
		return errScan
	}
	fmt.Printf("count of lines: %d", n)
	return nil
}
