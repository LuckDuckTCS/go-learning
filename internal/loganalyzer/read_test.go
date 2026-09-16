package loganalyzer

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const goodLine = `192.168.1.1 - - [01/Aug/2026:00:00:00 +0300] "GET /api/users HTTP/1.1" 200 1234 0.045` + "\n"

func TestProcess(t *testing.T) {

	base := time.Date(2026, 8, 1, 12, 0, 0, 0, time.FixedZone("", 3*60*60))

	tests := []struct {
		name       string
		input      string
		from, to   time.Time
		wantTotal  int
		wantBroken int
	}{
		{
			name:      "three good lines",
			input:     goodLine + goodLine + goodLine,
			wantTotal: 3,
		},
		{
			name:       "one broken among good",
			input:      goodLine + "trash\n" + goodLine,
			wantTotal:  3,
			wantBroken: 1,
		},
		{
			name:  "empty input",
			input: "",
		},
		{
			name:       "three broken",
			input:      goodLine[:20] + "\n" + "trash\n" + goodLine[:30] + "\n",
			wantTotal:  3,
			wantBroken: 3,
		},
		{
			name:      "time filter",
			input:     testLine(base) + testLine(base.Add(time.Hour)) + testLine(base.Add(2*time.Hour)),
			from:      base.Add(30 * time.Minute),
			to:        base.Add(90 * time.Minute),
			wantTotal: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a := NewAggregator()
			err := Process(strings.NewReader(tc.input), a, tc.from, tc.to)

			if err != nil {
				t.Fatalf("process() error: %v", err)
			}
			if a.total != tc.wantTotal {
				t.Errorf("total = %d, want %d", a.total, tc.wantTotal)
			}
			if a.broken != tc.wantBroken {
				t.Errorf("broken = %d, want %d", a.broken, tc.wantBroken)
			}
		})
	}
}

func testLine(ts time.Time) string {
	return fmt.Sprintf(`192.168.1.1 - - [%s] "GET /api/users HTTP/1.1" 200 1234 0.045`,
		ts.Format("02/Jan/2006:15:04:05 -0700")) + "\n"
}

func TestProcessBrokenLineNumber(t *testing.T) {
	input := goodLine + "trash\n" + goodLine // битая вторая

	a := NewAggregator()
	if err := Process(strings.NewReader(input), a, time.Time{}, time.Time{}); err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	// len(a.errs) == 1
	if len(a.errs) != 1 {
		t.Fatalf("len slice error = %d, want 1", len(a.errs))
	}
	// errors.As до *ParseError
	var target *ParseError
	if !errors.As(a.errs[0], &target) {
		t.Fatalf("errs[0] = %T, want *ParseError", a.errs[0])
	}
	// pe.Line == 2
	if target.Line != 2 {
		t.Errorf("broke line number = %d, want 2", target.Line)
	}
}

// ветка scanner.Err()

type errReader struct{ err error }

func (r errReader) Read([]byte) (int, error) { return 0, r.err }

func TestProcessScannerErr(t *testing.T) {
	a := NewAggregator()
	wantErr := errors.New("read failed")

	err := Process(errReader{err: wantErr}, a, time.Time{}, time.Time{})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Process() error = %v, want %v in chain", err, wantErr)
	}
}

// файловые тесты

func TestProcessFile(t *testing.T) {
	a := NewAggregator()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.log")
	content := goodLine + "trash\n" + goodLine
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := ProcessFile(path, a, time.Time{}, time.Time{}); err != nil {
		t.Fatalf("ProcessFile() error = %v", err)
	}
	if a.total != 3 {
		t.Errorf("total = %d, want 3", a.total)
	}
	if a.broken != 1 {
		t.Errorf("broken = %d, want 1", a.broken)
	}

	a2 := NewAggregator()
	missing := filepath.Join(tmpDir, "nofile.log")
	if err := ProcessFile(missing, a2, time.Time{}, time.Time{}); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("err = %v, want os.ErrNotExist", err)
	}
}

func TestScanDir(t *testing.T) {
	dir := t.TempDir()
	errDir := os.MkdirAll(filepath.Join(dir, "sub"), 0o755)
	if errDir != nil {
		t.Fatalf("MkdirAll() error = %v", errDir)
	}
	files := []string{"a.log", "b.log", "notes.txt", "sub/c.log", "sub/readme.md"}
	for _, f := range files {
		path := filepath.Join(dir, filepath.FromSlash(f))
		if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
			t.Fatalf("WriteFile(%q) error = %v", f, err)
		}
	}

	got, err := ScanDir(dir)
	if err != nil {
		t.Fatalf("ScanDir() err = %v", err)
	}
	// проверить содержимое got
	var relPaths []string
	for _, p := range got {
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			t.Fatalf("Rel(%q, %q) error = %v", dir, p, err)
		}
		relPaths = append(relPaths, filepath.ToSlash(rel))
	}

	want := []string{"a.log", "b.log", "sub/c.log"}
	if diff := cmp.Diff(want, relPaths, cmpopts.SortSlices(func(a, b string) bool { return a < b })); diff != "" {
		t.Errorf("ScanDir() mismatch (-want +got):\n%s", diff)
	}

	gotEmpty, errEmpty := ScanDir(t.TempDir())
	if len(gotEmpty) != 0 {
		t.Errorf("len ScanDir(t.TempDir()) = %d, want 0", len(gotEmpty))
	}
	if errEmpty != nil {
		t.Errorf("error ScanDir(t.TempDir()) = %v, want nil", errEmpty)
	}
}

func BenchmarkProcess(b *testing.B) {
	data, err := os.ReadFile("testdata/bench.log")
	if err != nil {
		b.Skipf("no bench data: %v", err)
	}

	b.ResetTimer()
	for b.Loop() {
		a := NewAggregator()
		_ = Process(bytes.NewReader(data), a, time.Time{}, time.Time{})
	}
}
