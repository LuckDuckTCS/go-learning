package main

import (
	"fmt"
	"log"
	"os"

	"github.com/LuckDuckTCS/go-learning/internal/links"
)

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

type batch struct {
	urls  []string
	depth int
}

const maxDepth = 0

func main() {
	worklist := make(chan batch)
	var n int // количество ожидающих отправок в worklist

	// Начало с аргументов командной строки
	n++
	go func() {
		urlF := batch{urls: os.Args[1:], depth: 0}
		worklist <- urlF
	}()

	// Конкурентный обход веб-страниц
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		if list.depth > maxDepth {
			continue
		}
		for _, link := range list.urls {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string, depth int) {
					newLinks := crawl(link)
					worklist <- batch{urls: newLinks, depth: depth + 1}
				}(link, list.depth)
			}
		}
	}
}
