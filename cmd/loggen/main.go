package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	ips      = []string{"192.168.1.1", "10.0.0.42", "172.16.254.1", "8.8.8.8", "1.1.1.1", "192.168.0.100", "10.10.10.10", "172.31.0.5", "8.8.4.4", "77.88.8.8", "192.168.88.1", "10.0.1.15", "172.20.10.2", "45.33.32.156", "93.184.216.34", "185.199.108.153", "192.168.31.1", "10.8.0.1", "172.17.0.1", "208.67.222.222"} // штук 20-30
	urls     = []string{"/", "/api/users", "/login", "/logout", "/register", "/api/v1/status", "/api/auth/token", "/dashboard", "/profile", "/settings", "/api/posts/latest", "/api/products/1024", "/api/cart/checkout", "/about", "/contact", "/pricing", "/api/v1/metrics", "/admin/panel", "/api/comments/create", "/static/css/main.css"}
	methods  = []string{"GET", "GET", "GET", "POST", "PUT", "DELETE"}
	statuses = []int{200, 200, 200, 200, 404, 500, 301}
	start    = time.Date(2026, 8, 1, 0, 0, 0, 0, time.Local)
)

type InvalidFlagError struct {
	Flag  string
	Value string
}

func (e *InvalidFlagError) Error() string {
	return fmt.Sprintf("invalid value %q for -%s", e.Value, e.Flag)
}

func main() {

	if err := run(os.Args[1:]); err != nil {
		var ife *InvalidFlagError
		fmt.Fprintf(os.Stderr, "loggen: %v\n", err)
		if errors.As(err, &ife) {
			fmt.Fprintf(os.Stderr, "hint: run -help to see valid values for -%s\n", ife.Flag)
		}
		os.Exit(1)
	}
}

func run(args []string) (err error) {

	fs := flag.NewFlagSet("loggen", flag.ContinueOnError)
	path := fs.String("out", "testdata/big.log", "path to log file")
	lines := fs.Int("lines", 10, "count of lines")
	broken := fs.Int("broken", 0, "percent of broken lines")
	longLine := fs.Int("long", 0, "instead of log lines, write one line of N bytes")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	if *broken < 0 {
		return &InvalidFlagError{Flag: "broken", Value: strconv.Itoa(*broken)}
	}
	if *lines < 0 {
		return &InvalidFlagError{Flag: "lines", Value: strconv.Itoa(*lines)}
	}

	if err := os.MkdirAll(filepath.Dir(*path), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	f, err := os.Create(*path)
	if err != nil {
		return fmt.Errorf("create %s: %w", *path, err)
	}

	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close output: %w", closeErr)
		}
	}()

	w := bufio.NewWriter(f)

	defer func() {
		if flushErr := w.Flush(); flushErr != nil && err == nil {
			err = fmt.Errorf("flush output: %w", flushErr)
		}
	}()
	if *longLine > 0 {
		return writeLongLine(w, *longLine)
	}
	return generate(w, *lines, *broken)
}

func generate(w io.Writer, lines, broken int) error {
	for i := 0; i < lines; i++ {
		line := logLine(i)

		if shouldBreak(i, broken) {
			line = breakLine(line)
		}

		if _, err := io.WriteString(w, line); err != nil {
			return fmt.Errorf("write line %d: %w", i, err)
		}
	}
	return nil
}

func logLine(i int) string {
	ip := ips[i%len(ips)]
	url := urls[(i*7)%len(urls)] // умножение на простое число разводит совпадения
	status := statuses[(i*13)%len(statuses)]
	method := methods[i%len(methods)]
	size := 200 + (i*31)%9800

	ts := start.Add(time.Duration(i) * 50 * time.Millisecond)
	tsStr := ts.Format("02/Jan/2006:15:04:05 -0700")

	dur := 0.005 + float64(i%97)*0.001 // от 5 мс до ~100 мс
	if i%1000 == 0 {
		dur = 2.5 // редкие медленные запросы
	}

	return fmt.Sprintf(`%s - - [%s] "%s %s HTTP/1.1" %d %d %.3f`, ip, tsStr, method, url, status, size, dur) + "\n"
}

func shouldBreak(i, pct int) bool {
	if pct > 100 {
		pct = 100
	}
	if pct <= 0 {
		return false
	}
	return i%(100/pct) == 0
}

func breakLine(line string) string {

	return line[:len(line)/2] + "\n"
}

func writeLongLine(w io.Writer, size int) error {
	const chunk = 64 * 1024
	buf := bytes.Repeat([]byte("x"), chunk)

	for written := 0; written < size; written += chunk {
		n := chunk
		if remaining := size - written; remaining < chunk {
			n = remaining
		}
		if _, err := w.Write(buf[:n]); err != nil {
			return fmt.Errorf("write long line: %w", err)
		}
	}

	if _, err := io.WriteString(w, "\n"); err != nil {
		return fmt.Errorf("write newline: %w", err)
	}
	return nil
}
