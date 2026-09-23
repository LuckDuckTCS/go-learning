package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

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
		result = append(result, path)
		return nil
	})
	if err != nil {
		return result, fmt.Errorf("scan dir %s: %w", dir, err)
	}
	return result, nil
}

type result struct {
	path string
	sum  string
	err  error
}

func main() {
	start := time.Now()
	pathDir := os.Args[1:]
	if len(pathDir) == 0 {
		log.Fatal("error: empty dir")
	}
	if len(pathDir) > 1 {
		log.Fatal("error: more than one dir")
	}

	paths, err := ScanDir(pathDir[0])
	if err != nil {
		log.Fatalf("scan dir: %v", err)
		os.Exit(2)
	}
	results := make([]result, len(paths))

	var wg sync.WaitGroup

	for i, path := range paths {
		wg.Add(1)
		go func(i int, path string) {
			defer wg.Done()
			f, err := os.Open(path)
			if err != nil {
				log.Printf("open file: %v", err)
				results[i] = result{path: "", sum: "", err: err}
				return
			}
			defer f.Close()
			h := sha256.New()
			_, err = io.Copy(h, f)
			if err != nil {
				log.Printf("copy from file: %v", err)
				results[i] = result{path: "", sum: "", err: err}
				return
			}
			sum := h.Sum(nil)
			results[i] = result{path: path, sum: fmt.Sprintf("%x", sum), err: nil}
		}(i, path)

	}
	wg.Wait()
	fmt.Fprintln(os.Stderr, time.Since(start))
	/*for _, data := range results {
		fmt.Printf("path: %s\thash: %s\n", data.path, data.sum)
	}*/
}
