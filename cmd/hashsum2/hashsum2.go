package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"
)

const N = 10 // число воркеров

func ScanDir(ctx context.Context, dir string, paths chan<- string) {
	defer close(paths)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if d != nil && d.IsDir() {
				log.Printf("access denied %s: %v\n", path, err)
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		select {
		case paths <- path:
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	})
	if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		log.Printf("scan dir %s: %v", dir, err)
	}
}

type result struct {
	path string
	sum  string
	err  error
}

func main() {
	pathsDir := os.Args[1:]
	if len(pathsDir) == 0 {
		log.Fatal("error: empty dir")
	}
	if len(pathsDir) > 1 {
		log.Fatal("error: more than one dir")
	}

	paths := make(chan string)
	results := make(chan result)

	var wg sync.WaitGroup
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	go ScanDir(ctx, pathsDir[0], paths)

	for range N {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range paths {
				select {
				case results <- hashFile(path):
				case <-ctx.Done():
					return
				}
			}

		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for data := range results {
		if data.err != nil {
			log.Printf("!error!\tpath: %s\terr: %s\n", data.path, data.err)
		} else {
			fmt.Printf("path: %s\thash: %s\n", data.path, data.sum)
		}
	}
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
		return result{path: path, err: err}
	}
	sum := h.Sum(nil)

	return result{path: path, sum: fmt.Sprintf("%x", sum), err: err}
}
