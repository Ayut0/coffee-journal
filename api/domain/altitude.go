package domain

import "fmt"

type Altitude struct {
	min int
	max int
}

func NewAltitude(min, max int) (Altitude, error) {
	if min < 0 {
		return Altitude{}, fmt.Errorf("min altitude must be >= 0: %d", min)
	}
	if max < 0 {
		return Altitude{}, fmt.Errorf("max altitude must be >= 0: %d", max)
	}
	if min > max {
		return Altitude{}, fmt.Errorf("min altitude must be <= max: %d > %d", min, max)
	}
	return Altitude{min: min, max: max}, nil
}

func (a Altitude) Min() int {
	return a.min
}

func (a Altitude) Max() int {
	return a.max
}

