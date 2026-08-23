// Package stat - вычисляем максимум, минимум, среднее значение
package stat

import (
	"fmt"
)

type Summary struct {
	Count int
	Sum   float64
	Min   float64
	Max   float64
	Avg   float64
}

func Describe(nums []float64) (Summary, error) {
	result := Summary{}
	if len(nums) == 0 {
		return result, fmt.Errorf("пустой набор")
	}
	result.Min = nums[0]
	result.Max = nums[0]
	result.Avg = nums[0]
	result.Count = len(nums)
	for _, num := range nums {
		if num < result.Min {
			result.Min = num
		}
		if num > result.Max {
			result.Max = num
		}
		result.Sum += num
	}
	result.Avg = result.Sum / float64(result.Count)
	return result, nil
}
