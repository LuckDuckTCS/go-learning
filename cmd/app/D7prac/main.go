package main

import (
	"fmt"
	"os"

	"github.com/LuckDuckTCS/go-learning/internal/counters"
)

func main() {
	var c counters.ByteCounter
	var Name = "Bob"
	fmt.Fprintf(&c, "Hello, %s", Name)
	fmt.Println(c)

	d, n := counters.CountingWriter(os.Stdout)

	fmt.Fprintf(d, "Hers, %s\n", Name)
	fmt.Fprintf(d, "Hello, %s\n", Name)
	fmt.Println(*n)
}
