package cmd

import (
	"cli/env"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockAppContext(t *testing.T, opts ...AppOptions) *AppContext {
	var o AppOptions

	if len(opts) > 0 {
		o = opts[0]
	}

	return o.AppContext()
}

func TestVerifyRequirements(t *testing.T) {
	assert.NoError(t, MockAppContext(t).VerifyRequirements())
}

func TestVerifyRequirements_Missing(t *testing.T) {
	tmp := env.ProjectRequiredApps
	env.ProjectRequiredApps = append(env.ProjectRequiredApps, "invalid-app")

	err := MockAppContext(t).VerifyRequirements()
	assert.EqualError(t, err, "Some requirements are not met: [invalid-app]")

	env.ProjectRequiredApps = tmp
}
