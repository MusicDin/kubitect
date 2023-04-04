package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLongDesc(t *testing.T) {
	assert.Equal(t, "", LongDesc(" \n \n "))
	assert.Equal(t, "a", LongDesc("\n a \n"))
	assert.Equal(t, "a b", LongDesc("\n a b \n"))
	assert.Equal(t, "a\nb", LongDesc("\n a \n b \n"))
}

func TestExample(t *testing.T) {
	assert.Equal(t, "  \n  \n", Example(" \n \n "))
	assert.Equal(t, "  a\n", Example("\n a \n"))
	assert.Equal(t, "  a b\n", Example("\n a b \n"))
	assert.Equal(t, "  a\n  b\n", Example("\n a \n b \n"))
}

func TestPresetName(t *testing.T) {
	assert.Equal(t, ".", presetName(""))
	assert.Equal(t, ".yml", presetName(".yml"))
	assert.Equal(t, ".yml", presetName("test/.yml"))
	assert.Equal(t, "test", presetName("test"))
	assert.Equal(t, "test", presetName("test/test"))
	assert.Equal(t, "test", presetName("test.yml"))
	assert.Equal(t, "test", presetName("test/test.yml"))
	assert.Equal(t, "test.test", presetName("test.test.yml"))
}
