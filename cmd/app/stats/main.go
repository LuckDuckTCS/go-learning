package main

import (
	"fmt"
	"io"
	"os"

	"github.com/LuckDuckTCS/go-learning/internal/parse"
	"github.com/LuckDuckTCS/go-learning/internal/render"
	"github.com/LuckDuckTCS/go-learning/internal/stat"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	// parse.Numbers → stat.Describe → render.Text
	nums, err := parse.Numbers(args)
	if err != nil {
		return err
	}
	summ, err2 := stat.Describe(nums)
	if err2 != nil {
		return err
	}
	err3 := render.Text(out, summ)
	if err3 != nil {
		return err
	}
	return nil
}
