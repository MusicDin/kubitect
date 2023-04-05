package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddonRook(t *testing.T) {
	rook := Rook{
		Version:      Version("v1.9.9"),
		NodeSelector: Labels{"rook": "true"},
	}

	assert.NoError(t, rook.Validate())
	assert.NoError(t, Rook{}.Validate())
}
