package ansible

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/ui"

	"github.com/stretchr/testify/assert"
)

func TestExtraVarsToMap(t *testing.T) {
	vars := []string{"key=value"}
	expect := map[string]string{"key": "value"}

	varsMap, err := extraVarsToMap(vars)
	assert.NoError(t, err)
	assert.Equal(t, expect, varsMap)
}

func TestExtraVarsToMap_Empty(t *testing.T) {
	vars := []string{}
	expect := map[string]string{}

	varsMap, err := extraVarsToMap(vars)
	assert.NoError(t, err)
	assert.Equal(t, expect, varsMap)
}

func TestExtraVarsToMap_Invalid(t *testing.T) {
	vars := []string{"invalid"}

	_, err := extraVarsToMap(vars)
	assert.EqualError(t, err, "extraVarsToMap: variable (invalid) must be in 'key=value' format")
}

func TestAnsible_InvalidPath(t *testing.T) {
	a := NewAnsible(t.TempDir(), "")

	pb := Playbook{}
	assert.EqualError(t, a.Exec(pb), "ansible-playbook: playbook path not set")
}

func TestAnsible_InvalidInventory(t *testing.T) {
	a := NewAnsible(t.TempDir(), "")

	pb := Playbook{
		Path: "pb.yaml",
	}

	assert.EqualError(t, a.Exec(pb), "ansible-playbook (pb.yaml): inventory not set")
}

func TestAnsible_InvalidExtraVar(t *testing.T) {
	a := NewAnsible(t.TempDir(), "")

	pb := Playbook{
		Path:      "pb.yaml",
		Local:     true,
		ExtraVars: []string{"invalid"},
	}

	assert.EqualError(t, a.Exec(pb), "extraVarsToMap: variable (invalid) must be in 'key=value' format")
}

func TestAnsible_InvalidBinPath(t *testing.T) {
	a := NewAnsible(t.TempDir(), "")

	pb := Playbook{
		Path:      "pb.yaml",
		Inventory: "localhost",
	}

	assert.ErrorContains(t, a.Exec(pb), "ansible-playbook (pb.yaml): Binary file")
}

func TestAnsible_InvalidBinPath2(t *testing.T) {
	a := NewAnsible(t.TempDir(), "")

	ui.MockGlobalUi(t, ui.UiOptions{Debug: true, NoColor: true})

	pb := Playbook{
		Path:      "pb.yaml",
		Local:     true,
		ExtraVars: []string{"key=value"},
	}

	assert.ErrorContains(t, a.Exec(pb), "ansible-playbook (pb.yaml): Binary file")
}
