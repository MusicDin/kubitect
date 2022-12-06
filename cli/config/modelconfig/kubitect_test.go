package modelconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubitect_Empty(t *testing.T) {
	url := URL("https://github.com/musicdin/kubitect")
	version := MasterVersion("master")

	k := Kubitect{
		Url:     &url,
		Version: &version,
	}

	assert.NoError(t, k.Validate())
	assert.NoError(t, Kubitect{}.Validate())
}
