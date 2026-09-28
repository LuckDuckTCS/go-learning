// Fetch1 конкурентно загружает URL и возвращает первый полученный ответ.
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, url := range os.Args[1:] {
		go fetch(ctx, url, ch) // запуск горутины
	}
	fmt.Println(<-ch) // первый полученный ответ
	cancel()
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
	fmt.Println(runtime.NumGoroutine())
}

// fetch загружает URL и отправляет в ch время выполнения и размер ответа.

func fetch(ctx context.Context, url string, ch chan<- string) {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		select {
		case ch <- fmt.Sprint(err):
		case <-ctx.Done():
		}
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		select {
		case ch <- fmt.Sprint(err):
		case <-ctx.Done():
		}
		return
	}
	nbytes, err := io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		select {
		case ch <- fmt.Sprintf("while reading %s: %v", url, err):
		case <-ctx.Done():
		}
		return
	}
	secs := time.Since(start).Seconds()
	select {
	case ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url):
	case <-ctx.Done():
	}
}
