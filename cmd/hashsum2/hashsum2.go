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

const N = 10 // число воркеров

func ScanDir(dir string, paths chan<- string) {
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if d != nil && d.IsDir() {
				log.Printf("skip %s: %v\n", path, err)
				return fs.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		paths <- path
		return nil
	})
	if err != nil {
		log.Printf("scan dir %s: %v", dir, err)
	}
	close(paths)
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

	paths := make(chan string)
	results := make(chan result)

	go ScanDir(pathDir[0], paths)

	var wg sync.WaitGroup

	for range N {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range paths {
				results <- hashFile(path)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for data := range results {
		if data.err != nil {
			fmt.Printf("!error!\tpath: %s\terr: %s\n", data.path, data.err)
		} else {
			fmt.Printf("path: %s\thash: %s\n", data.path, data.sum)
		}
	}
	fmt.Fprintln(os.Stderr, time.Since(start))
}

func hashFile(path string) result {
	f, err := os.Open(path)
	if err != nil {
		return result{path: path, err: err}
	}
	defer f.Close()

	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return result{path: path, sum: "", err: err}
	}
	sum := h.Sum(nil)
	return result{path: path, sum: fmt.Sprintf("%x", sum), err: nil}
}
