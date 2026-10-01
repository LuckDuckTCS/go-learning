package main

import (
	"cmp"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"time"

	"github.com/LuckDuckTCS/go-learning/internal/linkchecker"
)

type InvalidFlagError struct {
	Flag  string
	Value string
}

func (e *InvalidFlagError) Error() string {
	return fmt.Sprintf("invalid value %q for -%s", e.Value, e.Flag)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := run(ctx, os.Args[1:], os.Stdout); err != nil {
		var ife *InvalidFlagError
		fmt.Fprintf(os.Stderr, "linkchecker: %v\n", err)
		if errors.As(err, &ife) {
			fmt.Fprintf(os.Stderr, "hint: run -help to see valid values for -%s\n", ife.Flag)
		}
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, out io.Writer) (err error) {
	fs := flag.NewFlagSet("linkchecker", flag.ContinueOnError)
	format := fs.String("format", "text", "output format: text|json")
	workers := fs.Int("workers", 10, "number of url check goroutines")
	attempts := fs.Int("attempts", 5, "number of attempts to connect to the server")
	timeout := fs.Int("timeout", 5, "connect timeout (seconds)")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}
	// валидация workers
	if *workers <= 0 {
		return &InvalidFlagError{Flag: "workers", Value: strconv.Itoa(*workers)}
	}
	// валидация attempts
	if *attempts <= 0 {
		return &InvalidFlagError{Flag: "attempts", Value: strconv.Itoa(*attempts)}
	}
	// валидация timeout
	if *timeout <= 0 {
		return &InvalidFlagError{Flag: "timeout", Value: strconv.Itoa(*timeout)}
	}

	// формат
	/// ...
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: *workers,
		IdleConnTimeout:     90 * time.Second,

		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(*timeout) * time.Second,
	}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // не следовать, вернуть 301 как есть
	}

	// вызов форматтеров
	var outFormat linkchecker.Formatter
	switch *format {
	case "text":
		outFormat = linkchecker.TextFormatter{}
	case "json":
		outFormat = linkchecker.JSONFormatter{}
	default:
		return fmt.Errorf("unknown format %q (want text|json)", *format)
	}

	rest := fs.Args()
	var results []linkchecker.Result

	if len(rest) > 1 {
		return fmt.Errorf("more than one input file")
	}
	if len(rest) == 0 {
		results, err = linkchecker.Process(ctx, os.Stdin, client, *workers, *attempts) // stdin напрямую
	} else {
		// определить: файл или директория
		info, statErr := os.Stat(rest[0])
		if statErr != nil {
			return fmt.Errorf("stat input: %w", statErr)
		}
		if info.IsDir() {
			return fmt.Errorf("want file, got directory")
		}
		results, err = linkchecker.ProcessFile(ctx, rest[0], client, *workers, *attempts)
	}

	// отсортируем срез
	slices.SortFunc(results, func(a, b linkchecker.Result) int {
		if a.URL != b.URL {
			return cmp.Compare(a.URL, b.URL)
		}

		return cmp.Compare(a.Duration, b.Duration)
	})

	if len(results) == 0 {
		return err
	}
	rep := linkchecker.BuildReport(results)
	formatErr := outFormat.Format(out, rep)
	return errors.Join(formatErr, err)
}
