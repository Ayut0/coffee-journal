package domain

import "fmt"

// Aftertaste is a value object representing the lingering quality of a coffee's finish.
// It has no identity — two Aftertastes with the same value are equal.
type Aftertaste string

const (
	AftertasteShort    Aftertaste = "Short"
	AftertasteClean    Aftertaste = "Clean"
	AftertasteLingering Aftertaste = "Lingering"
)

// ParseAftertaste validates and creates an Aftertaste from a string.
// This is the only way to create one — enforcing validity at birth.
func ParseAftertaste(s string) (Aftertaste, error) {
	switch Aftertaste(s) {
	case AftertasteShort, AftertasteClean, AftertasteLingering:
		return Aftertaste(s), nil
	default:
		return "", fmt.Errorf("invalid aftertaste %q: must be Short, Clean, or Lingering", s)
	}
}

func (a Aftertaste) String() string {
	return string(a)
}
