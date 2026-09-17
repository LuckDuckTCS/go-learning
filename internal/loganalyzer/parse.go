package loganalyzer

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Entry struct {
	IP       string
	Time     time.Time
	Method   string
	URL      string
	Status   int
	Size     int
	Duration float64
}

type ParseError struct {
	Line int
	Text string
	Err  error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("in line %d: %v: %q", e.Line, e.Err, e.Text)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

var lineRe = regexp.MustCompile(`^(\S+) - - \[([^\]]+)\] "(\S+) (\S+) \S+" (\d+) (\d+) ([\d.]+)$`)

// ParseLine эталонная реализация, разбирает строку лога
func ParseLine(s string) (Entry, error) {
	m := lineRe.FindStringSubmatch(s)
	if m == nil {
		return Entry{}, ErrNoMatch
	}
	var result Entry
	result.IP = m[1]
	t, err := time.Parse("02/Jan/2006:15:04:05 -0700", m[2])
	if err != nil {
		return Entry{}, fmt.Errorf("parse time: %w", err)
	}
	result.Time = t
	result.Method = m[3]
	result.URL = m[4]
	status, err := strconv.Atoi(m[5])
	if err != nil {
		return Entry{}, fmt.Errorf("parse status: %w", err)
	}
	result.Status = status
	size, err := strconv.Atoi(m[6])
	if err != nil {
		return Entry{}, fmt.Errorf("parse size: %w", err)
	}
	result.Size = size

	duration, err := strconv.ParseFloat(m[7], 64)
	if err != nil {
		return Entry{}, fmt.Errorf("parse duration: %w", err)
	}
	result.Duration = duration

	return result, nil
}

var ErrNoMatch = errors.New("line does not match log format")

// ParseLineManual разбирает строку лога без регулярных выражений.
func ParseLineManual(s string) (Entry, error) {
	// 1. IP до первого пробела
	ip, rest, ok := strings.Cut(s, " ")
	if !ok {
		return Entry{}, ErrNoMatch
	}
	if !isToken(ip) {
		return Entry{}, ErrNoMatch
	}
	// 2. пропустить "- - " — проверить, что оно там есть
	rest, ok = strings.CutPrefix(rest, "- - ")
	if !ok {
		return Entry{}, ErrNoMatch
	}
	// 3. время между [ и ]
	rest, ok = strings.CutPrefix(rest, "[")
	if !ok {
		return Entry{}, ErrNoMatch
	}
	timeRaw, rest, ok := strings.Cut(rest, "] ")
	if !ok {
		return Entry{}, ErrNoMatch
	}
	timeT, err := time.Parse("02/Jan/2006:15:04:05 -0700", timeRaw)
	if err != nil {
		return Entry{}, fmt.Errorf("parse time: %w", err)
	}
	// 4. запрос между кавычками: метод, URL, протокол
	rest, ok = strings.CutPrefix(rest, `"`)
	if !ok {
		return Entry{}, ErrNoMatch
	}

	method, rest, ok := strings.Cut(rest, " ")
	if !ok {
		return Entry{}, ErrNoMatch
	}
	if !isToken(method) {
		return Entry{}, ErrNoMatch
	}

	url, rest, ok := strings.Cut(rest, " ")
	if !ok {
		return Entry{}, ErrNoMatch
	}
	if !isToken(url) {
		return Entry{}, ErrNoMatch
	}

	prot, rest, ok := strings.Cut(rest, `" `)
	if !ok {
		return Entry{}, ErrNoMatch
	}
	if !isToken(prot) {
		return Entry{}, ErrNoMatch
	}
	// 5. три числа через пробел: статус, размер, длительность
	status, rest, ok := strings.Cut(rest, " ")
	if !ok {
		return Entry{}, ErrNoMatch
	}
	size, rest, ok := strings.Cut(rest, " ")
	if !ok {
		return Entry{}, ErrNoMatch
	}
	dur := rest

	if !isNumeric(status) {
		return Entry{}, ErrNoMatch
	}
	statusInt, err := strconv.Atoi(status)
	if err != nil {
		return Entry{}, fmt.Errorf("parse status: %w", err)
	}
	if !isNumeric(size) {
		return Entry{}, ErrNoMatch
	}
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		return Entry{}, fmt.Errorf("parse size: %w", err)
	}

	if !isNumeric(dur) {
		return Entry{}, ErrNoMatch
	}
	durFloat, err := strconv.ParseFloat(dur, 64)
	if err != nil {
		return Entry{}, fmt.Errorf("parse duration: %w", err)
	}

	return Entry{
		IP:       ip,
		Time:     timeT,
		Method:   method,
		URL:      url,
		Status:   statusInt,
		Size:     sizeInt,
		Duration: durFloat,
	}, nil
}

func isToken(s string) bool {
	return s != "" && !strings.ContainsAny(s, " \t\n\f\r")
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) && r != '.' {
			return false
		}
	}
	return true
}
