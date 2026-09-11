package loganalyzer

import (
	"cmp"
	"slices"
)

const maxExamples = 10

type Aggregator struct {
	total     int
	broken    int
	byStatus  map[int]int
	byURL     map[string]int
	byIP      map[string]int
	durations []float64
	errs      []error
}

func NewAggregator() *Aggregator {
	var result Aggregator
	result.byStatus = make(map[int]int)
	result.byURL = make(map[string]int)
	result.byIP = make(map[string]int)
	return &result
}

func (a *Aggregator) Add(e Entry) {
	a.total++
	a.byStatus[e.Status]++
	a.byURL[e.URL]++
	a.byIP[e.IP]++
	a.durations = append(a.durations, e.Duration)
}

func (a *Aggregator) AddBroken(err error) {
	a.total++
	a.broken++
	if len(a.errs) < maxExamples {
		a.errs = append(a.errs, err)
	}

}

func (a *Aggregator) Report(top int) Report {
	result := Report{}
	result.Total = a.total
	result.Broken = a.broken
	result.ByStatus = a.byStatus
	result.TopURLs = topN(a.byURL, top)
	result.TopIPs = topN(a.byIP, top)
	result.AvgDur = avg(a.durations)
	result.P95Dur = p95(a.durations)

	examples := make([]string, 0, len(a.errs))
	for _, err := range a.errs {
		examples = append(examples, err.Error())
	}
	result.Examples = examples
	return result
}

func topN(m map[string]int, n int) []Pair {
	// из карты в срез
	result := make([]Pair, 0, len(m))
	for key, value := range m {
		result = append(result, Pair{Key: key, Count: value})
	}
	slices.SortFunc(result, func(a, b Pair) int {
		if c := cmp.Compare(b.Count, a.Count); c != 0 {
			return c
		}
		return cmp.Compare(a.Key, b.Key)
	})
	if n > 0 && n < len(result) {
		result = result[:n]
	}
	return result
}

func avg(f []float64) float64 {
	sum := 0.0
	if len(f) == 0 {
		return 0
	}
	for _, dur := range f {
		sum += dur
	}
	return sum / float64(len(f))
}

func p95(f []float64) float64 {
	if len(f) == 0 {
		return 0
	}
	slices.Sort(f)
	idx := int(float64(len(f)) * 0.95)
	if idx >= len(f) {
		idx = len(f) - 1
	}
	return f[idx]
}
