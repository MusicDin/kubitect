package modelconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddonRook(t *testing.T) {
	ver := Version("v1.9.9")

	rook := Rook{
		Version: &ver,
		NodeSelector: &Labels{
			"rook": "true",
		},
	}

	assert.NoError(t, rook.Validate())
	assert.NoError(t, Rook{}.Validate())
}
