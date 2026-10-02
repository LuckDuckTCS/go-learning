package linkchecker

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func Process(ctx context.Context, r io.Reader, client *http.Client, workers, attempts int) ([]Result, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(workers)
	results := make(chan Result, workers/2)
	errorChan := make(chan error, 1)

	go func() {
		defer close(results)
		for scanner.Scan() {
			url := scanner.Text()
			g.Go(func() error {
				res := CheckWithRetry(ctx, client, url, attempts)
				select {
				case results <- res:
				case <-ctx.Done():
					return ctx.Err()
				}
				return nil
			})
		}
		scanErr := scanner.Err()
		errorChan <- errors.Join(scanErr, g.Wait())
	}()

	out := make([]Result, 0)
	for result := range results {
		out = append(out, result)
	}
	if err := <-errorChan; err != nil {
		return out, fmt.Errorf("checking links: %w", err)
	}
	return out, nil
}

func ProcessFile(ctx context.Context, path string, client *http.Client, workers, attempts int) (_ []Result, err error) {
	f, err := os.Open(path)
	if err != nil {
		return []Result{}, fmt.Errorf("open file %s: %w", path, err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close input %s: %w", path, closeErr)
		}
	}()
	result, err := Process(ctx, f, client, workers, attempts)
	if err != nil {
		return result, fmt.Errorf("process %s: %w", path, err)
	}
	return result, nil
}
