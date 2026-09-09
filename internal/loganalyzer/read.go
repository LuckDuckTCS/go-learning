package loganalyzer

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func ScanLines(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	n := 0
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		n++
		//if n%500000 == 0 {
		//	var ms runtime.MemStats
		//	runtime.ReadMemStats(&ms)
		//	fmt.Printf("строк %d, куча %d МБ\n", n, ms.HeapAlloc/1024)
		//}
	}
	if err := scanner.Err(); err != nil {
		return n, fmt.Errorf("scanning lines: %w", err)
	}
	return n, nil
}

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

func ScanLinesSlice(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines)%500000 == 0 {
			var ms runtime.MemStats
			runtime.ReadMemStats(&ms)
			fmt.Printf("строк %d, куча %d МБ\n", len(lines), ms.HeapAlloc/1024/1024)
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scanning lines: %w", err)
	}
	return len(lines), nil
}

func CountFile(path string) (n int, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("open file: %w", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close input: %w", closeErr)
		}
	}()
	n, err = ScanLines(f)
	if err != nil {
		return 0, fmt.Errorf("scan lines: %w", err)
	}
	return n, nil
}
