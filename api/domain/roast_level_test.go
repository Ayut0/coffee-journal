package domain_test

import (
	"testing"

	"github.com/Ayut0/coffee-journal/api/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRoastLevel_ValidValues(t *testing.T) {
	// Table-driven test — a Go convention for testing multiple inputs
	// against the same logic. Each row is one scenario.
	tests := []struct {
		input string
		want  string
	}{
		{"Light", "Light"},
		{"Medium", "Medium"},
		{"Dark", "Dark"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			rl, err := domain.NewRoastLevel(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, rl.String())
		})
	}
}

func TestNewRoastLevel_InvalidValue(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"unknown value", "Extra Dark"},
		{"lowercase", "light"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewRoastLevel(tt.input)
			assert.Error(t, err)
		})
	}
}

func TestRoastLevel_Equality(t *testing.T) {
	// Value object property: two objects with the same value ARE equal.
	// No identity — just the value matters.
	a, _ := domain.NewRoastLevel("Light")
	b, _ := domain.NewRoastLevel("Light")
	assert.Equal(t, a, b)
}
