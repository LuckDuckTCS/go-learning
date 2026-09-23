package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {

	args := os.Args[1:]

	tzs := make([]string, 0, 3)
	targets := make([]string, 0, 3)
	for i := range args {
		tz, target, f := strings.Cut(args[i], "=")
		if !f {
			log.Fatal(errors.New("wrong flag format"))
		}
		tzs = append(tzs, tz)
		targets = append(targets, target)
	}

	var wg sync.WaitGroup
	for i := range targets {
		wg.Add(1)
		go func(tz, target string) {
			defer wg.Done()
			conn, err := net.Dial("tcp", target)
			if err != nil {
				log.Printf("Dial err: %v", err)
				return
			}
			defer conn.Close()
			scanner := bufio.NewScanner(conn)

			for scanner.Scan() {
				line := scanner.Text()
				fmt.Println(tz+":", line)
			}

			if err := scanner.Err(); err != nil {
				log.Printf("scan error: %v", err)
			}

		}(tzs[i], targets[i])
	}
	wg.Wait()
}
