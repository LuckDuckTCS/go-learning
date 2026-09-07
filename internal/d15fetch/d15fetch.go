package d15fetch

import (
	"io"
	"net/http"
	"os"
	"path"
)

//Упражнение 5.18. Перепишите, не изменяя ее поведение, функцию fetch так,
//чтобы она использовала defer для закрытия записываемого файла.

//Fetch загружает URL и возвращает имя и длину локального файла.

func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()

	n, err = io.Copy(f, resp.Body)

	return local, n, err
}
