package validation

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	assert.False(t, isEmpty(42))
	assert.False(t, isEmpty("42"))
	assert.False(t, isEmpty([]int{42}))
	assert.False(t, isEmpty(&[]int{42}))
	assert.False(t, isEmpty(map[string]int{"test": 42}))
	assert.False(t, isEmpty(&map[string]string{"test": "42"}))

	assert.True(t, isEmpty(nil))
	assert.True(t, isEmpty(0))
	assert.True(t, isEmpty(""))
	assert.True(t, isEmpty([]int{}))
	assert.True(t, isEmpty(&[]int{}))
	assert.True(t, isEmpty(map[string]int{}))
	assert.True(t, isEmpty(&map[string]string{}))
}

func TestDeepValue(t *testing.T) {
	a := 42
	ra := reflect.ValueOf(a)
	assert.Equal(t, a, getDeepValue(ra).Interface())
}
