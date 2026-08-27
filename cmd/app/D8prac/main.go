package main

import (
	"fmt"

	"golang.org/x/net/html"

	"github.com/LuckDuckTCS/go-learning/internal/reader"
)

const page = `<html>
<head><title>Тест</title></head>
<body>
<h1>Заголовок</h1>
<p>Текст со <a href="/one">первой ссылкой</a>.</p>
<div><a href="https://go.dev">Go</a></div>
<a name="anchor">якорь без href</a>
</body>
</html>`

func main() {
	doc, err := html.Parse(reader.NewReader(page))
	if err == nil {
		links := visit(nil, doc)
		fmt.Printf("%q", links)
	}

}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}
