package loganalyzer

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const maxExampleLen = 80

func ScanDir(dir string) ([]string, error) {
	result := make([]string, 0, 10)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if d != nil && d.IsDir() {
				fmt.Fprintf(os.Stderr, "skip %s: %v\n", path, err)
				return fs.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(strings.ToLower(d.Name()), ".log") {
			result = append(result, path)
		}
		return nil
	})
	if err != nil {
		return result, fmt.Errorf("scan dir %s: %w", dir, err)
	}
	return result, nil
}

func ProcessFile(path string, a *Aggregator, from, to time.Time) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file %s: %w", path, err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close input %s: %w", path, closeErr)
		}
	}()
	err = Process(f, a, from, to)
	if err != nil {
		return fmt.Errorf("process %s: %w", path, err)
	}
	return nil
}

func Process(r io.Reader, a *Aggregator, from, to time.Time) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		e, err := ParseLine(line)
		if err != nil {
			a.AddBroken(&ParseError{Line: lineNum, Text: truncate(line), Err: err})
			continue
		}

		if !from.IsZero() && e.Time.Before(from) {
			continue
		}
		if !to.IsZero() && !e.Time.Before(to) {
			continue
		}
		a.Add(e)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanning lines: %w", err)
	}
	return nil
}

func truncate(s string) string {
	if maxExampleLen >= len(s) {
		return s
	}
	return s[:maxExampleLen] + "..."
}
