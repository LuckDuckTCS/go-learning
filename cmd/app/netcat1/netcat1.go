package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {

	fs := flag.NewFlagSet("netcat1", flag.ExitOnError)
	port := fs.Int("port", 8000, "port")

	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(2)
	}

	conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(*port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(os.Stdout, conn)
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
