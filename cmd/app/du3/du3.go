package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var vFlag = flag.Bool("v", false, "показывать подробный ход выполнения")

type fileSize struct {
	root int // индекс корня
	size int64
}

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	fileSizes := make(chan fileSize)
	var n sync.WaitGroup
	for i, root := range roots {
		n.Add(1)
		go walkDir(root, i, &n, fileSizes)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	var tick <-chan time.Time
	if *vFlag {
		tick = time.Tick(500 * time.Millisecond)
	}

	nfiles := make([]int64, len(roots))
	nbytes := make([]int64, len(roots))
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes закрыт
			}
			nfiles[size.root]++
			nbytes[size.root] += size.size
		case <-tick:
			printDiskUsage(roots, nfiles, nbytes)
		}
	}
	printDiskUsage(roots, nfiles, nbytes) // окончательный итог
}

func printDiskUsage(roots []string, nfiles, nbytes []int64) {
	for i, root := range roots {
		fmt.Printf("%s: %d файлов  %.1f ГБ\n", root, nfiles[i], float64(nbytes[i])/1e9)
	}
}

// walkDir рекурсивно обходит дерево файлов с корнем dir
// и отправляет размер каждого файла в fileSizes.
func walkDir(dir string, root int, n *sync.WaitGroup, fileSizes chan<- fileSize) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, root, n, fileSizes)
		} else {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			fileSizes <- fileSize{root: root, size: info.Size()}
		}
	}
}

var sema = make(chan struct{}, 20) // счётный семафор для ограничения параллелизма

// dirents возвращает записи каталога dir.
func dirents(dir string) []os.DirEntry {
	sema <- struct{}{}        // захват токена
	defer func() { <-sema }() // освобождение токена

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}
