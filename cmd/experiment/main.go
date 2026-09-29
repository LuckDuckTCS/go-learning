package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const N = 1000000

func main() {
	chan1 := make(chan int, 100)
	chan2 := make(chan int, 100)
	var wg sync.WaitGroup
	start := time.Now()
	go func() { chan1 <- 0 }()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			message := <-chan1
			chan2 <- message
		}
		close(chan2)
		<-chan1 // читаем последнюю отправку третьей горутины
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for message := range chan2 {
			chan1 <- message
		}
	}()

	wg.Wait()
	fmt.Fprintln(os.Stderr, time.Since(start))
}
