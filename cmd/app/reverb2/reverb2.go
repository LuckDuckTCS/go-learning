// Reverb2 — конкурентный эхо-сервер, обрабатывающий каждое эхо в своей горутине.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	var wg sync.WaitGroup
	lines := make(chan string, 1)
	afkTimer := time.NewTimer(10 * time.Second)
	defer afkTimer.Stop()
	go func() {
		input := bufio.NewScanner(c)
		for input.Scan() {
			lines <- input.Text()
		}
		close(lines)
	}()
loop:
	for {

		select {
		case str, ok := <-lines:
			if !ok {
				break loop
			}
			wg.Add(1)
			afkTimer.Reset(10 * time.Second)
			go echo(c, str, 1*time.Second, &wg)
		case <-afkTimer.C:
			fmt.Fprintf(c, "Вы бездействовали 10 сек\n")
			break loop
		}
	}
	// Ошибки input.Err() игнорируются
	wg.Wait()
	c.Close()
}

//реализация без select
/*func handleConn2(c net.Conn) {
	var wg sync.WaitGroup
	input := bufio.NewScanner(c)

	c.SetReadDeadline(time.Now().Add(10 * time.Second))
	for input.Scan() {
		wg.Add(1)
		go echo(c, input.Text(), 1*time.Second, &wg)
		c.SetReadDeadline(time.Now().Add(10 * time.Second))
	}
	wg.Wait()
	c.Close()
}*/

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
