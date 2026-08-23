// Package parse - парсер чисел из строки
package parse

import (
	"fmt"
	"strconv"
)

func Numbers(args []string) ([]float64, error) {
	result := []float64{}
	if len(args) == 0 {
		return result, fmt.Errorf("пустой набор")
	}
	for i, numStr := range args {
		numF, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return result, fmt.Errorf("неверный аргумерт %s на позиции %d: %w", numStr, i, err)
		}
		result = append(result, numF)
	}
	return result, nil
}
