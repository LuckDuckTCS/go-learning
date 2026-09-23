// Clock2 — конкурентный TCP-сервер, передающий время.
package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {

	fs := flag.NewFlagSet("clock2", flag.ExitOnError)
	port := fs.Int("port", 8000, "port")

	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(2)
	}

	var err error
	loc := time.Local
	if tz := os.Getenv("TZ"); tz != "" {
		loc, err = time.LoadLocation(tz)
		if err != nil {
			log.Fatal(err)
		}
	}

	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(*port))
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // например, сбой соединения
			continue
		}
		go handleConn(conn, loc) // обработка соединений одновременно
	}
}

func handleConn(c net.Conn, loc *time.Location) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().In(loc).Format("15:04:05\n"))
		if err != nil {
			return // например, клиент отключился
		}
		time.Sleep(1 * time.Second)
	}
}
