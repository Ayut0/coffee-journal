package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// OAuthAccount is a value object representing a linked OAuth provider account.
// Fields are exported because it is passed across layers as plain data.
type OAuthAccount struct {
	Provider       string
	ProviderUserID string
}

// User is the aggregate root representing an application user.
// All fields are private; access via getters only.
type User struct {
	id           uuid.UUID
	email        string
	passwordHash *string
	name         string
	createdAt    time.Time
	updatedAt    time.Time
	deletedAt    *time.Time
}

// NewUser creates a validated User with a new UUID and current timestamps.
// passwordHash is always nil at creation — set separately for password-based auth.
func NewUser(email, name string) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	return &User{
		id:           uuid.New(),
		email:        email,
		passwordHash: nil,
		name:         name,
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
		deletedAt:    nil,
	}, nil
}

func (u *User) ID() uuid.UUID           { return u.id }
func (u *User) Email() string           { return u.email }
func (u *User) PasswordHash() *string   { return u.passwordHash }
func (u *User) Name() string            { return u.name }
func (u *User) CreatedAt() time.Time    { return u.createdAt }
func (u *User) UpdatedAt() time.Time    { return u.updatedAt }
func (u *User) DeletedAt() *time.Time   { return u.deletedAt }

// UserRepository defines the persistence contract for User aggregates.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
}
