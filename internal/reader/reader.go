package reader

import "io"

type StringReader struct {
	s string
	i int // номер след. байта
}

func (s *StringReader) Read(p []byte) (int, error) {
	if s.i >= len(s.s) {
		return 0, io.EOF
	}
	n := copy(p, s.s[s.i:])
	s.i += n
	return n, nil
}

func NewReader(s string) *StringReader {
	return &StringReader{s: s}
}

type LimitRead struct {
	r io.Reader
	n int64
}

func (r *LimitRead) Read(p []byte) (int, error) {
	if r.n <= 0 {
		return 0, io.EOF
	}
	if r.n < int64(len(p)) {
		p = p[:r.n]
	}
	n, err := r.r.Read(p)
	r.n -= int64(n)
	return n, err

}

func LimitReader(r io.Reader, n int64) io.Reader {
	var res LimitRead
	res.r = r
	res.n = n
	return &res
}
