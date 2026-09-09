package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/LuckDuckTCS/go-learning/internal/tempconv"
)

// tempFlag реализует flag.Value, храня значение в градусах Цельсия.
type tempFlag struct {
	tempconv.Celsius
}

func (f *tempFlag) Set(s string) error {
	var value float64
	var unit string
	fmt.Sscanf(s, "%f%s", &value, &unit)

	switch unit {
	case "C", "°C":
		f.Celsius = tempconv.Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = tempconv.FToC(tempconv.Fahrenheit(value))
		return nil
	case "K", "°K":
		f.Celsius = tempconv.KToC(tempconv.Kelvin(value))
		return nil
	}
	return fmt.Errorf("invalid temperature %q", s)
}

// TempFlag регистрирует флаг температуры и возвращает указатель на значение.
func TempFlag(fs *flag.FlagSet, name string, value tempconv.Celsius, usage string) *tempconv.Celsius {
	f := tempFlag{value}
	fs.Var(&f, name, usage)
	return &f.Celsius
}

func main() {
	fs := flag.NewFlagSet("tempflag", flag.ContinueOnError)
	temp := TempFlag(fs, "temp", 20.0, "the temperature")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(*temp)
}
