package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
