package domain_test

import (
	"testing"

	"github.com/Ayut0/coffee-journal/api/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAltitude_ValidValues(t *testing.T) {
	tests := []struct{
		name string
		min int
		max int
	}{
		{"valid range", 1200, 1800},
		{"same value", 1500, 1500},
		{"zero values", 0, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewAltitude(tt.min, tt.max)
			assert.NoError(t, err)
		})
	}

}

func TestNewAltitude_InvalidValue(t *testing.T) {
	tests := []struct{
		name string
		min int
		max int
	}{
		{"negative min", -100, 100},
		{"negative max", 100, -100},
		{"min greater than max", 2000, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewAltitude(tt.min, tt.max)
			assert.Error(t, err)
		})
	}
}

func TestAltitude_Equality(t *testing.T) {
	a, _ := domain.NewAltitude(1200, 1800)
	b, _ := domain.NewAltitude(1200, 1800)
	assert.Equal(t, a, b)
}

func TestAltitude_Getters(t *testing.T){
	alt, err := domain.NewAltitude(1200, 1800)
	require.NoError(t, err)
	assert.Equal(t, 1200, alt.Min())
	assert.Equal(t, 1800, alt.Max())
}