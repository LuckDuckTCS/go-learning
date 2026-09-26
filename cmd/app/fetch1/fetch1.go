// Fetch1 конкурентно загружает URL и возвращает первый полученный ответ.
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	done := make(chan struct{})
	for _, url := range os.Args[1:] {
		go fetch(url, ch, done) // запуск горутины
	}
	fmt.Println(<-ch) // первый полученный ответ
	close(done)
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

// fetch загружает URL и отправляет в ch время выполнения и размер ответа.
// добавил отмену как того хочет упр. 8.11 (устаревшим форматом)
func fetch(url string, ch chan<- string, done <-chan struct{}) {
	start := time.Now()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		select {
		case ch <- fmt.Sprint(err):
		case <-done:
		}
		return
	}
	req.Cancel = done
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		select {
		case ch <- fmt.Sprint(err):
		case <-done:
		}
		return
	}
	nbytes, err := io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		select {
		case ch <- fmt.Sprintf("while reading %s: %v", url, err):
		case <-done:
		}
		return
	}
	secs := time.Since(start).Seconds()
	select {
	case ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url):
	case <-done:
	}
}
