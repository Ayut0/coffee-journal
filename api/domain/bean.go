package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Bean struct {
	id uuid.UUID
	userID uuid.UUID
	name string
	roaster string
	origin string
	roastLevel RoastLevel
	process *Process
	altitude *Altitude
	harvestSeason *string
	packagePhotoURL *string
	isPublic bool
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func NewBean(
	userID uuid.UUID,
	name string,
	roaster string,
	origin string,
	roastLevel RoastLevel,
) (*Bean, error) {
if userID == uuid.Nil {
	return nil, fmt.Errorf("userID is required")
}

if name == "" {
	return nil, fmt.Errorf("name is required")
}

if roaster == "" {
	return nil, fmt.Errorf("roaster is required")
}

if origin == "" {
	return nil, fmt.Errorf("origin is required")
}

return &Bean{
	id: uuid.New(),
	userID: userID,
	name: name,
	roaster: roaster,
	origin: origin,
	roastLevel: roastLevel,
	createdAt: time.Now(),
	updatedAt: time.Now(),
	deletedAt: nil,
}, nil
}

func (b *Bean) ID() uuid.UUID            { return b.id }
func (b *Bean) UserID() uuid.UUID        { return b.userID }
func (b *Bean) Name() string             { return b.name }
func (b *Bean) Roaster() string          { return b.roaster }
func (b *Bean) Origin() string           { return b.origin }
func (b *Bean) RoastLevel() RoastLevel   { return b.roastLevel }
func (b *Bean) Process() *Process        { return b.process }
func (b *Bean) Altitude() *Altitude      { return b.altitude }
func (b *Bean) HarvestSeason() *string   { return b.harvestSeason }
func (b *Bean) PackagePhotoURL() *string { return b.packagePhotoURL }
func (b *Bean) IsPublic() bool           { return b.isPublic }
func (b *Bean) CreatedAt() time.Time     { return b.createdAt }
func (b *Bean) UpdatedAt() time.Time     { return b.updatedAt }
func (b *Bean) DeletedAt() *time.Time    { return b.deletedAt }

type BeanRepository interface {
	Create(ctx context.Context, bean *Bean) error
	GetByID(ctx context.Context, id uuid.UUID) (*Bean, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*Bean, error)
	Update(ctx context.Context, bean *Bean) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}