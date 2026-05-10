package domain

import "fmt"

// RoastLevel is a value object representing coffee roast intensity.
// It has no identity — two RoastLevels with the same value are equal.
type RoastLevel string

const (
	RoastLight  RoastLevel = "Light"
	RoastMedium RoastLevel = "Medium"
	RoastDark   RoastLevel = "Dark"
)

// NewRoastLevel validates and creates a RoastLevel.
// This is the only way to create one — enforcing validity at birth.
func NewRoastLevel(s string) (RoastLevel, error) {
	switch RoastLevel(s) {
	case RoastLight, RoastMedium, RoastDark:
		return RoastLevel(s), nil
	default:
		return "", fmt.Errorf("invalid roast level %q: must be Light, Medium, or Dark", s)
	}
}

func (r RoastLevel) String() string {
	return string(r)
}
