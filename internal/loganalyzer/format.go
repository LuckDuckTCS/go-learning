package loganalyzer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"text/tabwriter"
)

type Pair struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

type Report struct {
	Total    int         `json:"total"`
	Broken   int         `json:"broken"`
	ByStatus map[int]int `json:"by_status"`
	TopURLs  []Pair      `json:"top_urls"`
	TopIPs   []Pair      `json:"top_ips"`
	AvgDur   float64     `json:"avg_duration_sec"`
	P95Dur   float64     `json:"p95_duration_sec"`
	Examples []string    `json:"broken_examples,omitempty"`
}

type Formatter interface {
	Format(w io.Writer, r Report) error
}

type TextFormatter struct{}
type JSONFormatter struct{}
type CSVFormatter struct{}

var _ Formatter = TextFormatter{}
var _ Formatter = JSONFormatter{}
var _ Formatter = CSVFormatter{}

func (TextFormatter) Format(w io.Writer, r Report) error {
	tw := tabwriter.NewWriter(w, 0, 8, 2, ' ', 0)

	_, err := fmt.Fprintf(tw, "Total lines:\t%d\nBroken:\t%d\n\n", r.Total, r.Broken)
	if err != nil {
		return fmt.Errorf("text formatter: %w", err)
	}

	_, err = fmt.Fprintf(tw, "Status codes:\n")
	if err != nil {
		return fmt.Errorf("text formatter: %w", err)
	}
	statuses := slices.Sorted(maps.Keys(r.ByStatus))
	for _, key := range statuses {
		_, err = fmt.Fprintf(tw, "\t%d\t%d\n", key, r.ByStatus[key])
		if err != nil {
			return fmt.Errorf("text formatter: %w", err)
		}
	}

	_, err = fmt.Fprintf(tw, "\nTop URLs:\n")
	if err != nil {
		return fmt.Errorf("text formatter: %w", err)
	}
	for _, pair := range r.TopURLs {
		_, err = fmt.Fprintf(tw, "\t%s\t%d\n", pair.Key, pair.Count)
		if err != nil {
			return fmt.Errorf("text formatter: %w", err)
		}
	}

	_, err = fmt.Fprintf(tw, "\nTop IPs:\n")
	if err != nil {
		return fmt.Errorf("text formatter: %w", err)
	}
	for _, pair := range r.TopIPs {
		_, err = fmt.Fprintf(tw, "\t%s\t%d\n", pair.Key, pair.Count)
		if err != nil {
			return fmt.Errorf("text formatter: %w", err)
		}
	}

	_, err = fmt.Fprintf(tw, "\nDuration:\n")
	if err != nil {
		return fmt.Errorf("text formatter: %w", err)
	}
	_, err = fmt.Fprintf(tw, "\tAverage:\t%.3f s\n\tp95:\t%.3f s\n\n", r.AvgDur, r.P95Dur)
	if err != nil {
		return fmt.Errorf("text formatter: %w", err)
	}

	if len(r.Examples) > 0 {
		_, err = fmt.Fprintf(tw, "\nExamples:\n")
		if err != nil {
			return fmt.Errorf("text formatter: %w", err)
		}
		for _, str := range r.Examples {
			_, err = fmt.Fprintf(tw, "\t%s\n", str)
			if err != nil {
				return fmt.Errorf("text formatter: %w", err)
			}
		}
	}
	if err := tw.Flush(); err != nil {
		return fmt.Errorf("text formatter: flush: %w", err)
	}
	return nil
}

func (CSVFormatter) Format(w io.Writer, r Report) error {
	cw := csv.NewWriter(w)
	// csv.Writer запоминает первую ошибку внутри; проверяем один раз через Error() в конце
	_ = cw.Write([]string{"section", "key", "value"})

	_ = cw.Write([]string{"summary", "total", strconv.Itoa(r.Total)})

	_ = cw.Write([]string{"summary", "broken", strconv.Itoa(r.Broken)})

	statuses := slices.Sorted(maps.Keys(r.ByStatus))
	for _, key := range statuses {
		_ = cw.Write([]string{"status", strconv.Itoa(key), strconv.Itoa(r.ByStatus[key])})

	}

	for _, pair := range r.TopURLs {
		_ = cw.Write([]string{"url", pair.Key, strconv.Itoa(pair.Count)})
	}

	for _, pair := range r.TopIPs {
		_ = cw.Write([]string{"ip", pair.Key, strconv.Itoa(pair.Count)})
	}

	_ = cw.Write([]string{"duration", "avg_sec", strconv.FormatFloat(r.AvgDur, 'f', 3, 64)})

	_ = cw.Write([]string{"duration", "p95_sec", strconv.FormatFloat(r.P95Dur, 'f', 3, 64)})

	for i, str := range r.Examples {
		_ = cw.Write([]string{"example", strconv.Itoa(i), str})
	}
	cw.Flush()
	return cw.Error()
}

func (JSONFormatter) Format(w io.Writer, r Report) error {
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		return fmt.Errorf("json formatter: %w", err)
	}
	return nil
}
