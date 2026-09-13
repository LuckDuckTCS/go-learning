package loganalyzer

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_avg(t *testing.T) {
	var tests = []struct {
		name  string
		input []float64
		want  float64
	}{
		{
			name:  "nil",
			input: nil,
			want:  0,
		},
		{
			name:  "empty",
			input: []float64{},
			want:  0,
		},
		{
			name:  "oneElem",
			input: []float64{5},
			want:  5,
		},
		{
			name:  "correct",
			input: []float64{1, 2, 3},
			want:  2,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := avg(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("avg() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_p95(t *testing.T) {

	hundred := make([]float64, 100)
	for i := range hundred {
		hundred[i] = float64(i + 1)
	}
	hundredR := make([]float64, 100)
	for i := range hundredR {
		hundredR[i] = float64(100 - i)
	}

	var tests = []struct {
		name  string
		input []float64
		want  float64
	}{
		{
			name:  "nil",
			input: nil,
			want:  0,
		},
		{
			name:  "empty",
			input: []float64{},
			want:  0,
		},
		{
			name:  "oneElem",
			input: []float64{5},
			want:  5,
		},
		{
			name:  "tenElem",
			input: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			want:  10,
		},
		{
			name:  "hundredElem",
			input: hundred,
			want:  96,
		},
		{
			name:  "hundredElemReverse",
			input: hundredR,
			want:  96,
		},
		{
			name:  "randomElem",
			input: []float64{3, 1, 4, 1, 5, 9, 2, 6, 5, 3},
			want:  9,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := p95(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("p95() mismatch (-want +got):\n%s", diff)
			}
		})
	}

	// отдельно проверим отсортировался ли массив (в моём алгоритме должен)
	// срез принадлежит агрегатору, его сортировка безвредна
	input := []float64{3, 1, 2}
	p95(input)
	if !slices.IsSorted(input) {
		t.Error("p95() input slice is not sorted, want - sorted")
	}
}

func Test_topN(t *testing.T) {
	var tests = []struct {
		name      string
		inputData map[string]int
		inputN    int
		want      []Pair
	}{
		{
			name:      "nil map",
			inputData: nil,
			inputN:    0,
			want:      []Pair{},
		},
		{
			name:      "empty map",
			inputData: map[string]int{},
			inputN:    5,
			want:      []Pair{},
		},
		{
			name:      "top 2 of 3",
			inputData: map[string]int{"/": 73, "/login": 20, "/about": 7},
			inputN:    2,
			want:      []Pair{{Key: "/", Count: 73}, {Key: "/login", Count: 20}},
		},
		{
			name:      "n zero means all",
			inputData: map[string]int{"/": 73, "/login": 20, "/about": 7},
			inputN:    0,
			want: []Pair{
				{Key: "/", Count: 73},
				{Key: "/login", Count: 20},
				{Key: "/about", Count: 7},
			},
		},
		{
			name:      "n greater than len",
			inputData: map[string]int{"/": 73, "/login": 20, "/about": 7},
			inputN:    100,
			want: []Pair{
				{Key: "/", Count: 73},
				{Key: "/login", Count: 20},
				{Key: "/about", Count: 7},
			},
		},
		{
			name:      "n equals len",
			inputData: map[string]int{"/": 73, "/login": 20, "/about": 7},
			inputN:    3,
			want: []Pair{
				{Key: "/", Count: 73},
				{Key: "/login", Count: 20},
				{Key: "/about", Count: 7},
			},
		},
		{
			name:      "equal counts sorted by key",
			inputData: map[string]int{"/b": 5, "/a": 5, "/c": 5},
			inputN:    3,
			want: []Pair{
				{Key: "/a", Count: 5},
				{Key: "/b", Count: 5},
				{Key: "/c", Count: 5},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := topN(tc.inputData, tc.inputN)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("topN() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAggregatorTotalBroken(t *testing.T) {

}

func TestAggregatorNoEntries(t *testing.T) {
	a := NewAggregator()
	report := a.Report(0)
	if report.Total != 0 {
		t.Errorf("wrong Total: %d, want: 0", report.Total)
	}
	if report.Broken != 0 {
		t.Errorf("wrong Broken: %d, want: 0", report.Broken)
	}
	if report.AvgDur != 0 {
		t.Errorf("wrong AvgDur: %g, want: 0", report.AvgDur)
	}
	if report.P95Dur != 0 {
		t.Errorf("wrong P95Dur: %g, want: 0", report.P95Dur)
	}
	if report.ByStatus == nil {
		t.Error("wrong ByStatus: nil, want: map[int]int{}")
	}
	if report.TopURLs == nil {
		t.Error("wrong TopURLs: nil, want: []Pair{}")
	}
	if report.TopIPs == nil {
		t.Error("wrong TopIPs: nil, want: []Pair{}")
	}
	if report.Examples == nil {
		t.Error("wrong Examples: nil, want: []string{}")
	}

	a1 := NewAggregator()
	a1.AddBroken(errors.New("something"))
	report = a1.Report(0)
	if report.Total != 1 {
		t.Errorf("(addBroken) wrong Total: %d, want: 1", report.Total)
	}
	if report.Broken != 1 {
		t.Errorf("(addBroken) wrong Broken: %d, want: 1", report.Broken)
	}
	if report.AvgDur != 0 {
		t.Errorf("(addBroken) wrong AvgDur: %g, want: 0", report.AvgDur)
	}
	if report.P95Dur != 0 {
		t.Errorf("(addBroken) wrong P95Dur: %g, want: 0", report.P95Dur)
	}

}

// testEntry возвращает валидную запись
func testEntry(url, ip string, status int, dur float64) Entry {
	return Entry{
		IP:       ip,
		Time:     time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
		Method:   "GET",
		URL:      url,
		Status:   status,
		Size:     100,
		Duration: dur,
	}
}

func TestAggregatorTotals(t *testing.T) {
	a := NewAggregator()

	a.Add(testEntry("/api/users", "10.0.0.1", 200, 0.1))
	a.Add(testEntry("/login", "10.0.0.2", 404, 0.2))
	a.Add(testEntry("/", "10.0.0.1", 500, 0.3))

	a.AddBroken(errors.New("broken line 1"))
	a.AddBroken(errors.New("broken line 2"))

	r := a.Report(0)

	if r.Total != 5 {
		t.Errorf("wrong Total: %d, want: 5", r.Total)
	}
	if r.Broken != 2 {
		t.Errorf("wrong Broken: %d, want: 2", r.Broken)
	}
}

func TestAggregatorExamplesLimit(t *testing.T) {
	a := NewAggregator()

	const total = maxExamples + 20

	for i := range total {
		a.AddBroken(fmt.Errorf("error %d", i))
	}

	r := a.Report(0)

	if r.Broken != total {
		t.Errorf("wrong Broken: %d, want: %d", r.Broken, total)
	}
	// len(r.Examples) == maxExamples
	if len(r.Examples) != maxExamples {
		t.Errorf("wrong len(Examples): %d, want: %d", len(r.Examples), maxExamples)
	}
	if len(r.Examples) > 0 {
		if !strings.Contains(r.Examples[0], "error 0") {
			t.Errorf("Examples[0] = %q, want first error", r.Examples[0])
		}
	}
}

func TestAggregatorReport(t *testing.T) {
	a := NewAggregator()

	a.Add(testEntry("/api/users", "10.0.0.1", 200, 0.1))
	a.Add(testEntry("/api/users", "10.0.0.1", 200, 0.2))
	a.Add(testEntry("/login", "127.0.0.53", 404, 0.3))
	a.Add(testEntry("/", "127.0.0.10", 500, 0.4))

	r := a.Report(0)

	// ByStatus == map[int]int{200: 2, 404: 1, 500: 1}
	byStatus := map[int]int{200: 2, 404: 1, 500: 1}
	if !maps.Equal(r.ByStatus, byStatus) {
		t.Errorf("wrong Report.ByStatus: %v, want: %v", r.ByStatus, byStatus)
	}
	// TopURLs[0] == Pair{Key: "/api/users", Count: 2}
	want := []Pair{
		{Key: "/api/users", Count: 2},
		{Key: "/", Count: 1},
		{Key: "/login", Count: 1},
	}
	if diff := cmp.Diff(want, r.TopURLs); diff != "" {
		t.Errorf("TopURLS mismatch (-want +got):\n%s", diff)
	}
	// AvgDur == 0.25
	if diff := cmp.Diff(0.25, r.AvgDur, cmpopts.EquateApprox(0.0, 1e-9)); diff != "" {
		t.Errorf("AvgDur mismatch (-want +got):\n%s", diff)
	}
	// P95Dur == 0.4
	if diff := cmp.Diff(0.4, r.P95Dur, cmpopts.EquateApprox(0.0, 1e-9)); diff != "" {
		t.Errorf("P95Dur mismatch (-want +got):\n%s", diff)
	}
}
