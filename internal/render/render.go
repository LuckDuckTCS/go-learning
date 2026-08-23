// Package render - выводит результат
package render

import (
	"fmt"
	"io"

	"github.com/LuckDuckTCS/go-learning/internal/stat"
)

func Text(w io.Writer, s stat.Summary) error {
	_, err := fmt.Fprintf(w, "Count: %d\nMin: %f\nMax: %f\nAvg: %f\nSum: %f", s.Count, s.Min, s.Max, s.Avg, s.Sum)
	return err
}
