package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})

	go func() {
		io.Copy(os.Stdout, conn) // note: ignoring errors
		log.Println("done")
		done <- struct{}{} // сигнал главной горутине
	}()

	mustCopy(conn, os.Stdin)
	tcp, ok := conn.(*net.TCPConn)
	if ok {
		if err := tcp.CloseWrite(); err != nil {
			conn.Close()
		}
	} else {
		conn.Close()
	}
	<-done // ожидание завершения фоновой горутины
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
