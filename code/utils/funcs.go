package utils

import (
	"math"
	"sort"
)

type PFunc func([]float64) float64

func Max(data []float64) float64 {
	max := math.SmallestNonzeroFloat32
	for _, n := range data {
		if n > max {
			max = n
		}
	}
	return max
}

func Min(data []float64) float64 {
	min := math.MaxFloat32
	for _, n := range data {
		if n < min {
			min = n
		}
	}
	return min
}

func Sum(data []float64) float64 {
	var sum float64
	for _, n := range data {
		sum += n
	}
	return sum
}

func Avg(data []float64) float64 {
	return Sum(data) / float64(len(data))
}

func Stdev(data []float64) float64 {
	a := Avg(data)
	var ret float64
	for _, n := range data {
		n -= a
		n *= n
		ret += n
	}
	ret /= float64(len(data))
	return math.Sqrt(ret)
}

func PercentFloat32(n int) PFunc {
	return func(data []float64) float64 {
		tmp := make([]float64, len(data))
		copy(tmp, data)
		sort.Slice(tmp, func(i, j int) bool {
			return tmp[i] < tmp[j]
		})
		n := int(float32(n)/100.) * len(tmp)
		return tmp[n]
	}
}
