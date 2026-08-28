package main

import (
	"fmt"
	"sort"

	"github.com/LuckDuckTCS/go-learning/internal/sortml"
)

func main() {
	fmt.Println(sortml.IsPalindrome(sort.IntSlice{1, 2, 3, 2, 1}))
	fmt.Println(sortml.IsPalindrome(sort.StringSlice{"a", "b", "a"}))
}
