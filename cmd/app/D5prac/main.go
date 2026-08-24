package main

import (
	"fmt"

	"github.com/LuckDuckTCS/go-learning/internal/shapes"
)

func main() {
	r := shapes.Rectangle{H: 10, W: 12}
	area := r.Area
	fmt.Printf("%f\n", area())
	r.H = 20
	fmt.Printf("%f\n", area())
}
