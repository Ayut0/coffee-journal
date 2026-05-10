package domain_test

import (
	"testing"

	"github.com/Ayut0/coffee-journal/api/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser_ValidInput(t *testing.T) {
	email := "test@example.com"
	name := "Yuto"

	u, err := domain.NewUser(email, name)
	require.NoError(t, err)
	require.NotNil(t, u)

	assert.Equal(t, email, u.Email())
	assert.Equal(t, name, u.Name())
	assert.NotEqual(t, [16]byte{}, u.ID()) // UUID is non-zero
	assert.Nil(t, u.PasswordHash())
	assert.Nil(t, u.DeletedAt())
	assert.False(t, u.CreatedAt().IsZero())
	assert.False(t, u.UpdatedAt().IsZero())
}

func TestNewUser_EmptyEmail(t *testing.T) {
	_, err := domain.NewUser("", "Yuto")
	assert.Error(t, err)
}

func TestNewUser_EmptyName(t *testing.T) {
	_, err := domain.NewUser("test@example.com", "")
	assert.Error(t, err)
}
