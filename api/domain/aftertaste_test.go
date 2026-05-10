package domain_test

import (
	"testing"

	"github.com/Ayut0/coffee-journal/api/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAftertaste_ValidValues(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Short", "Short"},
		{"Clean", "Clean"},
		{"Lingering", "Lingering"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			a, err := domain.ParseAftertaste(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, a.String())
		})
	}
}

func TestParseAftertaste_InvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"unknown value", "Bitter"},
		{"lowercase", "short"},
		{"mixed case", "SHORT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.ParseAftertaste(tt.input)
			assert.Error(t, err)
		})
	}
}

func TestAftertaste_String(t *testing.T) {
	assert.Equal(t, "Short", domain.AftertasteShort.String())
	assert.Equal(t, "Clean", domain.AftertasteClean.String())
	assert.Equal(t, "Lingering", domain.AftertasteLingering.String())
}

func TestAftertaste_Equality(t *testing.T) {
	// Value object property: two objects with the same value ARE equal.
	a, _ := domain.ParseAftertaste("Clean")
	b, _ := domain.ParseAftertaste("Clean")
	assert.Equal(t, a, b)
}
