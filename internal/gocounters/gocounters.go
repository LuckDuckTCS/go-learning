package gocounters

import (
	"sync"
)

type MutexCounter struct {
	mu sync.Mutex
	n  int
}

func (c *MutexCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}

func (c *MutexCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}

// вариант 2

type ChanCounter struct {
	inc   chan struct{}
	value chan int
	done  chan struct{}
}

func NewChanCounter() *ChanCounter {
	result := &ChanCounter{inc: make(chan struct{}), value: make(chan int), done: make(chan struct{})}
	go result.teller()
	return result
}

func (c *ChanCounter) Inc() {
	c.inc <- struct{}{}
}

func (c *ChanCounter) Value() int {
	return <-c.value
}

func (c *ChanCounter) Close() {
	close(c.done)
}

func (c *ChanCounter) teller() {
	var val int
	for {
		select {
		case <-c.inc:
			val++
		case c.value <- val:
		case <-c.done:
			return
		}
	}
}
