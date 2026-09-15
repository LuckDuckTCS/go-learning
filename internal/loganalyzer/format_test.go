package loganalyzer

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var update = flag.Bool("update", false, "update golden files")

func TestJSONFormatter(t *testing.T) {
	r := Report{
		Total:    10,
		Broken:   2,
		ByStatus: map[int]int{200: 8, 404: 2},
		TopURLs:  []Pair{{Key: "/", Count: 5}},
		TopIPs:   []Pair{{Key: "10.0.0.1", Count: 3}},
		AvgDur:   0.25,
		P95Dur:   0.4,
		Examples: []string{"line 1: bad"},
	}

	var buf bytes.Buffer
	if err := (JSONFormatter{}).Format(&buf, r); err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	var got Report
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if diff := cmp.Diff(r, got); diff != "" {
		t.Errorf("round-trip mismatch (-want +got):\n%s", diff)
	}
}

func TestCSVFormatter(t *testing.T) {
	r := Report{
		Total:    10,
		Broken:   2,
		ByStatus: map[int]int{200: 8, 404: 2},
		TopURLs:  []Pair{{Key: "/", Count: 5}},
		TopIPs:   []Pair{{Key: "10.0.0.1", Count: 3}},
		AvgDur:   0.25,
		P95Dur:   0.4,
		Examples: []string{"line 1: bad"},
	}

	want := [][]string{
		{"section", "key", "value"},
		{"summary", "total", "10"},
		{"summary", "broken", "2"},
		{"status", "200", "8"},
		{"status", "404", "2"},
		{"url", "/", "5"},
		{"ip", "10.0.0.1", "3"},
		{"duration", "avg_sec", "0.250"},
		{"duration", "p95_sec", "0.400"},
		{"example", "0", "line 1: bad"},
	}

	var buf bytes.Buffer
	if err := (CSVFormatter{}).Format(&buf, r); err != nil {
		t.Fatalf("Format() error = %v", err)
	}
	got, err := csv.NewReader(&buf).ReadAll()
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("round-trip mismatch (-want +got):\n%s", diff)
	}

	t.Run("Examples quotes", func(t *testing.T) {
		r.Examples = []string{`in line 7: bad: "GET /a, /b HTTP/1.1"`}
		var buf bytes.Buffer
		if err := (CSVFormatter{}).Format(&buf, r); err != nil {
			t.Fatalf("Format() error = %v", err)
		}
		got, err := csv.NewReader(&buf).ReadAll()
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}

		lastRow := got[len(got)-1]
		if lastRow[2] != r.Examples[0] {
			t.Errorf("example = %q, want %q", lastRow[2], r.Examples[0])
		}
	})
}

func TestTextFormatter(t *testing.T) {
	r := Report{
		Total:    10,
		Broken:   2,
		ByStatus: map[int]int{200: 8, 404: 2},
		TopURLs:  []Pair{{Key: "/", Count: 5}},
		TopIPs:   []Pair{{Key: "10.0.0.1", Count: 3}},
		AvgDur:   0.25,
		P95Dur:   0.4,
		Examples: []string{"line 1: bad"},
	}

	var buf bytes.Buffer
	if err := (TextFormatter{}).Format(&buf, r); err != nil {
		t.Fatalf("Format() error = %v", err)
	}

	golden := filepath.Join("testdata", "report.golden.txt")

	if *update {
		if err := os.WriteFile(golden, buf.Bytes(), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		t.Logf("golden file updated: %s", golden)
	}

	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if diff := cmp.Diff(string(want), buf.String()); diff != "" {
		t.Errorf("output mismatch (-want +got):\n%s", diff)
	}
}

type failWriter struct{ err error }

func (w failWriter) Write([]byte) (int, error) { return 0, w.err }

func TestFormatterError(t *testing.T) {
	r := Report{
		Total:    10,
		Broken:   2,
		ByStatus: map[int]int{200: 8, 404: 2},
		TopURLs:  []Pair{{Key: "/", Count: 5}},
		TopIPs:   []Pair{{Key: "10.0.0.1", Count: 3}},
		AvgDur:   0.25,
		P95Dur:   0.4,
		Examples: []string{"line 1: bad"},
	}

	wantErr := errors.New("write failed")
	fw := failWriter{err: wantErr}

	formatters := map[string]Formatter{
		"text": TextFormatter{},
		"json": JSONFormatter{},
		"csv":  CSVFormatter{},
	}

	for name, f := range formatters {
		t.Run(name, func(t *testing.T) {
			err := f.Format(fw, r)
			if err == nil {
				t.Fatal("Format() error = nil, want error")
			}
			if !errors.Is(err, wantErr) {
				t.Errorf("Format() error = %v, want %v in chain", err, wantErr)
			}
		})
	}
}
