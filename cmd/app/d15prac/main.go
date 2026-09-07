package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

func main() {

	err := SafeRun(func() error { return nil }) // nil
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("1 OK")
	}
	err = SafeRun(func() error { return io.EOF }) // io.EOF, не обёрнута
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("2 OK")
	}
	err = SafeRun(func() error { panic("boom") }) // ошибка с "boom"
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("3 OK")
	}
	err = SafeRun(func() error { var s []int; _ = s[5]; return nil }) // runtime error
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("4 OK")
	}

	err = SafeRun(func() error {
		go func() { panic("in goroutine") }()
		time.Sleep(100 * time.Millisecond)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("5 OK")
	}
}

func oneFileOpen(s string) error {
	f, err := os.Open(s)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func SafeRun(f func() error) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if pErr, ok := p.(error); ok {
				err = fmt.Errorf("panic: %w", pErr) // сохраняем цепочку
			} else {
				err = fmt.Errorf("panic: %v", p)
			}
			//fmt.Fprintf(os.Stderr, "panic: %v\n%s\n", p, debug.Stack())
		}
	}()
	return f()
}
