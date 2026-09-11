package loganalyzer

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
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

func ParseLine(s string) (Entry, error) {
	m := lineRe.FindStringSubmatch(s)
	if m == nil {
		return Entry{}, errors.New("line does not match log format")
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
