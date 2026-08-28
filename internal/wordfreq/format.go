package wordfreq

import (
	"encoding/json"
	"fmt"
	"io"
)

type Formatter interface {
	Format(w io.Writer, pairs []Pair) error
}

type TextFormatter struct{}
type JSONFormatter struct{}

var _ Formatter = TextFormatter{}
var _ Formatter = JSONFormatter{}

func (TextFormatter) Format(w io.Writer, pairs []Pair) error {
	for i := range pairs {
		_, err := fmt.Fprintf(w, "%s\t%d\n", pairs[i].Word, pairs[i].Count)
		if err != nil {
			return err
		}
	}
	return nil
}

func (JSONFormatter) Format(w io.Writer, pairs []Pair) error {
	err := json.NewEncoder(w).Encode(pairs)
	return err
}
