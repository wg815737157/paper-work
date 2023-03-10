package util

import "math"

type (
	ConnectionUtil interface {
		Max() interface{}
		Min() interface{}
		Sum() interface{}
	}
	FloatConnection []float64
)

func (fc FloatConnection) Max() interface{} {
	max := float64(0)
	for _, v := range fc {
		if v > max {
			max = v
		}
	}
	return max
}
func (fc FloatConnection) Min() interface{} {
	min := math.MaxFloat64
	for _, v := range fc {
		if v < min {
			min = v
		}
	}
	return min
}
func (fc FloatConnection) Sum() interface{} {
	sum := float64(0)
	for _, v := range fc {
		sum += v
	}
	return sum
}

func NewFloatConnection(f []float64) ConnectionUtil {
	return FloatConnection(f)
}
