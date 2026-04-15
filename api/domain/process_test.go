package domain_test

import (
	"testing"

	"github.com/Ayut0/coffee-journal/api/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProcess_ValidValues(t *testing.T){
	tests := []struct{
		input string
		want string
	}{
		{"Washed", "Washed"},
		{"Natural", "Natural"},
		{"Honey", "Honey"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T){
			process, err := domain.NewProcess(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.want, process.String())
		})
	}
}

func TestNewProcess_InvalidValue(t *testing.T) {
	tests := []struct{
		name string
		input string
	}{
		{"empty string", ""},
		{"unknown value", "Extra Washed"},
		{"lowercase", "washed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewProcess(tt.input)
			assert.Error(t, err)
		})
	}
}

func TestProcess_Equality(t *testing.T) {
	a,  _ := domain.NewProcess("Washed")
	b,  _ := domain.NewProcess("Washed")
	assert.Equal(t, a, b)
}