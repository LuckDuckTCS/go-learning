package linkchecker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	var calls atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		switch r.URL.Path {
		case "/ok":
			w.WriteHeader(http.StatusOK)
		case "/notfound":
			w.WriteHeader(http.StatusNotFound)
		case "/moved":
			w.Header().Set("Location", "/ok")
			w.WriteHeader(http.StatusMovedPermanently)
		}
	}))
	defer srv.Close()

	input := strings.NewReader(
		srv.URL + "/ok\n" +
			srv.URL + "/ok\n" +
			srv.URL + "/notfound\n" +
			srv.URL + "/moved\n" +
			"http://127.0.0.1:1\n", // закрытый порт — сетевая ошибка
	)

	client := &http.Client{Timeout: 2 * time.Second}

	results, err := Process(context.Background(), input, client, 3, 1)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	// проверить количество результатов
	if len(results) != 5 {
		t.Errorf("len(results) = %d, want = 5", len(results))
	}
	// посчитать категории через BuildReport и сравнить
	rep := BuildReport(results)
	if rep.OK+rep.Redirects+rep.Broken != rep.Total {
		t.Errorf("categories don't sum to total: %d+%d+%d != %d",
			rep.OK, rep.Redirects, rep.Broken, rep.Total)
	}
}
