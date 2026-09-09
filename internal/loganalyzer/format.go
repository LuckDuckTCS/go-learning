package loganalyzer

import (
	"encoding/json"
	"fmt"
	"io"
)

type Report struct {
	Total int `json:"total"`
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
	_, err := fmt.Fprintf(w, "%d\n", r.Total)
	if err != nil {
		return fmt.Errorf("text formatter: %w", err)
	}
	return nil
}

func (CSVFormatter) Format(w io.Writer, r Report) error {

	return nil
}

func (JSONFormatter) Format(w io.Writer, r Report) error {
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		return fmt.Errorf("json formatter: %w", err)
	}
	return nil
}
