package ansible

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/ui"

	"github.com/stretchr/testify/assert"
)

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
		ExtraVars: map[string]string{"key": "value"},
	}

	assert.ErrorContains(t, a.Exec(pb), "ansible-playbook (pb.yaml): Binary file")
}
