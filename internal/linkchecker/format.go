package linkchecker

import (
	"encoding/json"
	"fmt"
	"io"
)

type Formatter interface {
	Format(w io.Writer, r Report) error
}

type Report struct {
	Total     int
	OK        int
	Redirects int
	Broken    int
	Results   []Result
}

type resultJSON struct {
	URL      string `json:"url"`
	Status   int    `json:"status"`
	Location string `json:"location,omitempty"`
	Duration int64  `json:"duration_ms"`
	Error    string `json:"error,omitempty"`
}

type reportJSON struct {
	Total     int          `json:"total"`
	OK        int          `json:"ok"`
	Redirects int          `json:"redirects"`
	Broken    int          `json:"broken"`
	Results   []resultJSON `json:"results"`
}

type TextFormatter struct{}
type JSONFormatter struct{}

var _ Formatter = TextFormatter{}
var _ Formatter = JSONFormatter{}

func (TextFormatter) Format(w io.Writer, rep Report) error {
	// сводка: всего, ОК, редиректы, битые
	_, err := fmt.Fprintf(w, "\tReport:\nTotal:\t%d\nOK:\t%d\nRedirects:\t%d\nBroken:\t%d\n\n", rep.Total, rep.OK, rep.Redirects, rep.Broken)
	if err != nil {
		return fmt.Errorf("text formatter: %w", err)
	}
	// список редиректов: URL → Location
	_, err = fmt.Fprintf(w, "\tRedirects:\n")
	if err != nil {
		return fmt.Errorf("text formatter: %w", err)
	}
	for _, value := range rep.Results {
		if value.Status >= 300 && value.Status < 400 {
			_, err = fmt.Fprintf(w, "%s ---> %s\n", value.URL, value.Location)
			if err != nil {
				return fmt.Errorf("text formatter: %w", err)
			}
		}
	}
	// список битых: URL, статус или ошибка
	_, err = fmt.Fprintf(w, "\n\tBroken:\n")
	if err != nil {
		return fmt.Errorf("text formatter: %w", err)
	}
	for _, value := range rep.Results {
		if value.Err != nil || value.Status >= 400 {
			_, err = fmt.Fprintf(w, "%s\tstatus:%d error:%v\n", value.URL, value.Status, value.Err)
			if err != nil {
				return fmt.Errorf("text formatter: %w", err)
			}
		}
	}
	return nil
}
func (JSONFormatter) Format(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(toJSON(r))
}

func BuildReport(results []Result) Report {
	output := Report{Total: len(results), Results: results}
	for _, value := range results {
		if value.Err != nil || value.Status >= 400 {
			output.Broken++
			continue
		}
		if value.Status >= 200 && value.Status < 300 {
			output.OK++
			continue
		}
		if value.Status >= 300 && value.Status < 400 {
			output.Redirects++
			continue
		}
	}
	return output
}

func toJSON(r Report) reportJSON {
	out := reportJSON{
		Total:     r.Total,
		OK:        r.OK,
		Redirects: r.Redirects,
		Broken:    r.Broken,
		Results:   make([]resultJSON, 0, len(r.Results)),
	}
	for _, res := range r.Results {
		// заполнить resultJSON, преобразовав Duration и Err
		ms := res.Duration.Milliseconds()
		err := ""
		if res.Err != nil {
			err = res.Err.Error()
		}
		temp := resultJSON{URL: res.URL, Status: res.Status, Location: res.Location, Duration: ms, Error: err}
		out.Results = append(out.Results, temp)
	}
	return out
}
