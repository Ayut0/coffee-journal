package domain_test

import (
	"testing"

	"github.com/Ayut0/coffee-journal/api/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helpers used across tasting tests
var (
	validUserID = uuid.New()
	validBeanID = uuid.New()
)

func validTastingArgs() (uuid.UUID, uuid.UUID, []string, string, string, int, int, int, int, domain.Aftertaste) {
	return validUserID, validBeanID, []string{"chocolate", "citrus"}, "V60", "Medium", 3, 4, 3, 4, domain.AftertasteClean
}

func TestNewTasting_ValidInput(t *testing.T) {
	userID, beanID, tags, brew, grind, acidity, aroma, body, overall, aftertaste := validTastingArgs()

	tasting, err := domain.NewTasting(userID, beanID, tags, brew, grind, acidity, aroma, body, overall, aftertaste)
	require.NoError(t, err)

	assert.NotEqual(t, uuid.Nil, tasting.ID())
	assert.Equal(t, userID, tasting.UserID())
	assert.Equal(t, beanID, tasting.BeanID())
	assert.Equal(t, tags, tasting.FlavorTags())
	assert.Equal(t, brew, tasting.BrewMethod())
	assert.Equal(t, grind, tasting.GrindSize())
	assert.Equal(t, acidity, tasting.Acidity())
	assert.Equal(t, aroma, tasting.Aroma())
	assert.Equal(t, body, tasting.Body())
	assert.Equal(t, overall, tasting.Overall())
	assert.Equal(t, aftertaste, tasting.Aftertaste())
	assert.Nil(t, tasting.Sweetness())
	assert.Nil(t, tasting.NoteText())
	assert.Nil(t, tasting.DeletedAt())
	assert.False(t, tasting.CreatedAt().IsZero())
	assert.False(t, tasting.UpdatedAt().IsZero())
}

func TestNewTasting_NilFlavorTagsBecomesEmptySlice(t *testing.T) {
	userID, beanID, _, brew, grind, acidity, aroma, body, overall, aftertaste := validTastingArgs()

	tasting, err := domain.NewTasting(userID, beanID, nil, brew, grind, acidity, aroma, body, overall, aftertaste)
	require.NoError(t, err)
	assert.NotNil(t, tasting.FlavorTags())
	assert.Empty(t, tasting.FlavorTags())
}

func TestNewTasting_EmptyBrewMethod(t *testing.T) {
	userID, beanID, tags, _, grind, acidity, aroma, body, overall, aftertaste := validTastingArgs()

	_, err := domain.NewTasting(userID, beanID, tags, "", grind, acidity, aroma, body, overall, aftertaste)
	assert.Error(t, err)
}

func TestNewTasting_EmptyGrindSize(t *testing.T) {
	userID, beanID, tags, brew, _, acidity, aroma, body, overall, aftertaste := validTastingArgs()

	_, err := domain.NewTasting(userID, beanID, tags, brew, "", acidity, aroma, body, overall, aftertaste)
	assert.Error(t, err)
}

func TestNewTasting_NilUserID(t *testing.T) {
	_, beanID, tags, brew, grind, acidity, aroma, body, overall, aftertaste := validTastingArgs()

	_, err := domain.NewTasting(uuid.Nil, beanID, tags, brew, grind, acidity, aroma, body, overall, aftertaste)
	assert.Error(t, err)
}

func TestNewTasting_NilBeanID(t *testing.T) {
	userID, _, tags, brew, grind, acidity, aroma, body, overall, aftertaste := validTastingArgs()

	_, err := domain.NewTasting(userID, uuid.Nil, tags, brew, grind, acidity, aroma, body, overall, aftertaste)
	assert.Error(t, err)
}

func TestNewTasting_InvalidScore(t *testing.T) {
	tests := []struct {
		name             string
		acidity, aroma, body, overall int
	}{
		{"acidity too low", 0, 3, 3, 3},
		{"acidity too high", 6, 3, 3, 3},
		{"aroma too low", 3, 0, 3, 3},
		{"aroma too high", 3, 6, 3, 3},
		{"body too low", 3, 3, 0, 3},
		{"body too high", 3, 3, 6, 3},
		{"overall too low", 3, 3, 3, 0},
		{"overall too high", 3, 3, 3, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, beanID, tags, brew, grind, _, _, _, _, aftertaste := validTastingArgs()
			_, err := domain.NewTasting(userID, beanID, tags, brew, grind, tt.acidity, tt.aroma, tt.body, tt.overall, aftertaste)
			assert.Error(t, err)
		})
	}
}

func TestTasting_SetSweetness(t *testing.T) {
	userID, beanID, tags, brew, grind, acidity, aroma, body, overall, aftertaste := validTastingArgs()
	tasting, err := domain.NewTasting(userID, beanID, tags, brew, grind, acidity, aroma, body, overall, aftertaste)
	require.NoError(t, err)

	v := 4
	err = tasting.SetSweetness(&v)
	require.NoError(t, err)
	require.NotNil(t, tasting.Sweetness())
	assert.Equal(t, 4, *tasting.Sweetness())
}

func TestTasting_SetSweetness_Invalid(t *testing.T) {
	userID, beanID, tags, brew, grind, acidity, aroma, body, overall, aftertaste := validTastingArgs()
	tasting, err := domain.NewTasting(userID, beanID, tags, brew, grind, acidity, aroma, body, overall, aftertaste)
	require.NoError(t, err)

	v := 6
	err = tasting.SetSweetness(&v)
	assert.Error(t, err)
}

func TestTasting_SetNoteText(t *testing.T) {
	userID, beanID, tags, brew, grind, acidity, aroma, body, overall, aftertaste := validTastingArgs()
	tasting, err := domain.NewTasting(userID, beanID, tags, brew, grind, acidity, aroma, body, overall, aftertaste)
	require.NoError(t, err)

	note := "Bright and juicy with a long finish."
	tasting.SetNoteText(&note)
	require.NotNil(t, tasting.NoteText())
	assert.Equal(t, note, *tasting.NoteText())
}
