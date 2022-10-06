package validation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomError(t *testing.T) {
	assert.NotEmpty(t, Var("42", Max(0).Error("")))
	assert.Error(t, errors.New("test"), Var("42", Max(0).Error("test")))
}

func TestTags(t *testing.T) {
	assert.NoError(t, Var("42", Tags("")))
	assert.NoError(t, Var(42, Tags("max=42")))
	assert.NoError(t, Var(nil, Tags("omitempty,min=42")))
}

func TestRequired(t *testing.T) {
	assert.Error(t, Var(nil, Required()))
	assert.Error(t, Var("", Required()))
	assert.NoError(t, Var("test", Required()))
}

func TestOmitEmpty(t *testing.T) {
	assert.Error(t, Var(1, OmitEmpty(), Min(42)))
	assert.NoError(t, Var(nil, OmitEmpty(), Min(42)))
	assert.NoError(t, Var(0, OmitEmpty(), Min(42)))
	assert.NoError(t, Var([]string{}, OmitEmpty(), Min(42)))
	assert.NoError(t, Var(&[]string{}, OmitEmpty(), Min(42)))
}

func TestMin(t *testing.T) {
	assert.NoError(t, Var(42, Min(42)))
	assert.Error(t, Var(0, Min(42)))

	assert.Error(t, Var([]int{}, Min(1)))
	assert.NoError(t, Var([]int{42}, Min(1)))
}

func TestMax(t *testing.T) {
	assert.NoError(t, Var(42, Max(42)))
	assert.Error(t, Var(42, Max(0)))

	assert.NoError(t, Var([]int{}, Max(1)))
	assert.Error(t, Var([]int{42, 42}, Max(1)))
}

func TestLen(t *testing.T) {
	assert.NoError(t, Var([]string{"1"}, Len(1)))
	assert.Error(t, Var([]int{42, 42}, Len(1)))
}

func TestMinLen(t *testing.T) {
	assert.NoError(t, Var(42, MinLen(42)))
	assert.Error(t, Var(0, MinLen(42)))

	assert.Error(t, Var([]int{}, MinLen(1)))
	assert.NoError(t, Var([]int{42}, MinLen(1)))
}

func TestMaxLen(t *testing.T) {
	assert.NoError(t, Var(42, MaxLen(42)))
	assert.Error(t, Var(42, MaxLen(0)))

	assert.NoError(t, Var([]int{}, MaxLen(1)))
	assert.Error(t, Var([]int{42, 42}, MaxLen(1)))
}
