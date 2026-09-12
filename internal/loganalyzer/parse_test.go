package loganalyzer

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseLine(t *testing.T) {
	var tests = []struct {
		name    string
		input   string
		want    Entry
		wantErr bool
	}{
		{
			name:  "correct",
			input: `192.168.1.1 - - [01/Aug/2026:00:00:00 +0300] "GET /api/users HTTP/1.1" 200 1234 0.045`,
			want: Entry{
				IP:       "192.168.1.1",
				Time:     time.Date(2026, 8, 1, 0, 0, 0, 0, time.FixedZone("", 3*60*60)),
				Method:   "GET",
				URL:      "/api/users",
				Status:   200,
				Size:     1234,
				Duration: 0.045,
			},
		},
		{
			name:    "nonNumericStatus",
			input:   `192.168.1.1 - - [01/Aug/2026:00:00:00 +0300] "GET / HTTP/1.1" abc 1234 0.045`,
			want:    Entry{},
			wantErr: true,
		},
		{
			name:    "nonNumericSize",
			input:   `192.168.1.1 - - [01/Aug/2026:00:00:00 +0300] "GET / HTTP/1.1" 200 abc 0.045`,
			want:    Entry{},
			wantErr: true,
		},
		{
			name:    "nonNumericDuration",
			input:   `192.168.1.1 - - [01/Aug/2026:00:00:00 +0300] "GET / HTTP/1.1" 200 1234 abc`,
			want:    Entry{},
			wantErr: true,
		},
		{
			name:  "correct redirect",
			input: `10.0.0.42 - - [01/Aug/2026:00:00:01 +0300] "GET /login HTTP/1.1" 301 0 0.012`,
			want: Entry{
				IP:       "10.0.0.42",
				Time:     time.Date(2026, 8, 1, 0, 0, 1, 0, time.FixedZone("", 3*60*60)),
				Method:   "GET",
				URL:      "/login",
				Status:   301,
				Size:     0,
				Duration: 0.012,
			},
		},
		{
			name:    "notCorrect",
			input:   "192.168.1.1 - - [01/Aug/2026:00:00:00",
			want:    Entry{},
			wantErr: true,
		},
		{
			name:    "empty",
			input:   "",
			want:    Entry{},
			wantErr: true,
		},
		{
			name:    "brokeDate",
			input:   `192.168.1.1 - - [99/Xyz/2026:00:00:00 +0300] "GET / HTTP/1.1" 200 1234 0.045`,
			want:    Entry{},
			wantErr: true,
		},
		{
			name:    "extraSpaces",
			input:   `192.168.1.1 - -  [01/Aug/2026:00:00:00 +0300] "GET / HTTP/1.1" 200 1234 0.045`,
			want:    Entry{},
			wantErr: true,
		},
		{
			name:    "brokeURL",
			input:   `192.168.1.1 - - [01/Aug/2026:00:00:00 +0300] "GET /a b HTTP/1.1" 200 1234 0.045`,
			want:    Entry{},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseLine(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("ParseLine() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}

func TestParseErrorChain(t *testing.T) {
	const brokenDate = `192.168.1.1 - - [99/Xyz/2026:00:00:00 +0300] "GET / HTTP/1.1" 200 1234 0.045`

	_, inner := ParseLine(brokenDate)
	if inner == nil {
		t.Fatal("ParseLine() = nil error, want error")
	}

	pe := &ParseError{Line: 42, Text: "…", Err: inner}
	wrapped := fmt.Errorf("process file: %w", pe)

	// 1. errors.As находит *ParseError через внешнюю обёртку
	var target *ParseError
	if !errors.As(wrapped, &target) {
		t.Fatalf("errors.As(wrapped, *ParseError) = false, want true")
	}
	// 2. найденный ParseError содержит Line == 42
	if target.Line != 42 {
		t.Errorf("Line = %d, want 42", target.Line)
	}
	// 3. errors.Is доходит до inner
	if !errors.Is(wrapped, inner) {
		t.Errorf("errors.Is(wrapped, inner) = false, want true")
	}

	plain := errors.New("something else")
	var err *ParseError
	if errors.As(plain, &err) {
		t.Error("errors.As found ParseError in unrelated error")
	}
}
