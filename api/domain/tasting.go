package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Tasting is the root entity representing a single evaluation session for a coffee bean.
// All fields are private; access is exclusively through typed getters.
type Tasting struct {
	id         uuid.UUID
	userID     uuid.UUID
	beanID     uuid.UUID
	flavorTags []string
	brewMethod string
	grindSize  string
	acidity    int
	aroma      int
	body       int
	sweetness  *int
	overall    int
	aftertaste Aftertaste
	noteText   *string
	createdAt  time.Time
	updatedAt  time.Time
	deletedAt  *time.Time
}

// NewTasting creates and validates a new Tasting entity.
// Optional fields sweetness and noteText must be set separately via SetSweetness / SetNoteText.
func NewTasting(
	userID uuid.UUID,
	beanID uuid.UUID,
	flavorTags []string,
	brewMethod string,
	grindSize string,
	acidity int,
	aroma int,
	body int,
	overall int,
	aftertaste Aftertaste,
) (*Tasting, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("userID is required")
	}

	if beanID == uuid.Nil {
		return nil, fmt.Errorf("beanID is required")
	}

	if brewMethod == "" {
		return nil, fmt.Errorf("brewMethod is required")
	}

	if grindSize == "" {
		return nil, fmt.Errorf("grindSize is required")
	}

	if err := validateScore("acidity", acidity); err != nil {
		return nil, err
	}

	if err := validateScore("aroma", aroma); err != nil {
		return nil, err
	}

	if err := validateScore("body", body); err != nil {
		return nil, err
	}

	if err := validateScore("overall", overall); err != nil {
		return nil, err
	}

	// Ensure flavorTags is never nil — an empty slice is the canonical zero value.
	if flavorTags == nil {
		flavorTags = []string{}
	}

	now := time.Now()
	return &Tasting{
		id:         uuid.New(),
		userID:     userID,
		beanID:     beanID,
		flavorTags: flavorTags,
		brewMethod: brewMethod,
		grindSize:  grindSize,
		acidity:    acidity,
		aroma:      aroma,
		body:       body,
		overall:    overall,
		aftertaste: aftertaste,
		createdAt:  now,
		updatedAt:  now,
		deletedAt:  nil,
	}, nil
}

// validateScore checks that a score is in the valid 1–5 range.
func validateScore(field string, value int) error {
	if value < 1 || value > 5 {
		return fmt.Errorf("%s must be between 1 and 5, got %d", field, value)
	}
	return nil
}

// Getters — the only way to read Tasting state from outside the domain package.

func (t *Tasting) ID() uuid.UUID          { return t.id }
func (t *Tasting) UserID() uuid.UUID      { return t.userID }
func (t *Tasting) BeanID() uuid.UUID      { return t.beanID }
func (t *Tasting) FlavorTags() []string   { return t.flavorTags }
func (t *Tasting) BrewMethod() string     { return t.brewMethod }
func (t *Tasting) GrindSize() string      { return t.grindSize }
func (t *Tasting) Acidity() int           { return t.acidity }
func (t *Tasting) Aroma() int             { return t.aroma }
func (t *Tasting) Body() int              { return t.body }
func (t *Tasting) Sweetness() *int        { return t.sweetness }
func (t *Tasting) Overall() int           { return t.overall }
func (t *Tasting) Aftertaste() Aftertaste { return t.aftertaste }
func (t *Tasting) NoteText() *string      { return t.noteText }
func (t *Tasting) CreatedAt() time.Time   { return t.createdAt }
func (t *Tasting) UpdatedAt() time.Time   { return t.updatedAt }
func (t *Tasting) DeletedAt() *time.Time  { return t.deletedAt }

// Setters for optional fields.

// SetSweetness sets the nullable sweetness score.
// Pass a pointer to an int (1–5), or nil to clear it.
func (t *Tasting) SetSweetness(v *int) error {
	if v != nil {
		if err := validateScore("sweetness", *v); err != nil {
			return err
		}
	}
	t.sweetness = v
	t.updatedAt = time.Now()
	return nil
}

// SetNoteText sets the nullable free-form note text.
func (t *Tasting) SetNoteText(v *string) {
	t.noteText = v
	t.updatedAt = time.Now()
}

// TastingRepository defines the persistence contract for Tasting aggregates.
// All implementations must live in the repository layer — never in the domain.
type TastingRepository interface {
	Create(ctx context.Context, tasting *Tasting) error
	GetByID(ctx context.Context, id uuid.UUID) (*Tasting, error)
	ListByBean(ctx context.Context, beanID uuid.UUID) ([]*Tasting, error)
	// Timeline returns tastings for a user ordered by created_at descending,
	// using cursor-based pagination. Pass nil cursor for the first page.
	Timeline(ctx context.Context, userID uuid.UUID, cursor *time.Time, limit int) ([]*Tasting, error)
	Update(ctx context.Context, tasting *Tasting) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
}
