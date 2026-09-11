package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/LuckDuckTCS/go-learning/internal/loganalyzer"
)

var lineRe = regexp.MustCompile(`^(\S+) - - \[([^\]]+)\] "(\S+) (\S+) (\S+)" (\d+) (\d+) ([\d.]+)$`)

func main() {
	//line := `192.168.1.1 - - [01/Aug/2026:00:00:00 +0300] "GET /api/users HTTP/1.1" 200 1234 0.045`
	//fmt.Printf("%q\n", lineRe.FindStringSubmatch(line))
	//fmt.Printf("%q\n", lineRe.FindStringSubmatch(`192.168.1.1 - - [01/Aug`))

	_, err := loganalyzer.ParseLine(`192.168.1.1 - - [99/Xyz/2026:00:00:00 +0300] "GET / HTTP/1.1" 200 1 0.1`)
	pe := &loganalyzer.ParseError{Line: 7, Text: "...", Err: err}
	fmt.Println(pe)
	fmt.Println(errors.Is(pe, err)) // true
}
