package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/LuckDuckTCS/go-learning/internal/github"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "issues: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: issues <search terms>")
	}

	result, err := github.SearchIssues(args)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%d issues:\n", result.TotalCount)

	recent, thisYear, old := github.GroupByAge(result.Items, time.Now())
	fmt.Fprintf(out, "---------recent issues---------\n")
	for _, item := range recent {
		fmt.Fprintf(out, "#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
	fmt.Fprintf(out, "---------this year issues---------\n")
	for _, item := range thisYear {
		fmt.Fprintf(out, "#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
	fmt.Fprintf(out, "---------old issues---------\n")
	for _, item := range old {
		fmt.Fprintf(out, "#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
	return nil
}
