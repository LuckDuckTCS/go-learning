package charcount

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
)

func TestCountChars(t *testing.T) {
	var tests = []struct {
		name    string
		input   string
		want    Counts
		wantErr bool
	}{
		{
			name:  "ascii letters",
			input: "aab",
			want: Counts{Runes: map[rune]int{'a': 2, 'b': 1},
				Types:  map[string]int{"Letters": 3},
				UTFLen: [utf8.UTFMax + 1]int{1: 3}},
		},
		{
			name:  "empty input",
			input: "",
			want:  Counts{Runes: map[rune]int{}, Types: map[string]int{}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CountChars(strings.NewReader(tc.input))
			if (err != nil) != tc.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("CountChars() mismatch (-got +want):\n%s", diff)
			}
		})
	}

}
