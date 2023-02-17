package cluster

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	v "github.com/MusicDin/kubitect/cli/utils/validation"

	"github.com/stretchr/testify/assert"
)

type configMock struct {
	Value string
}

func (c configMock) Validate() error {
	return nil
}

type invalidConfigMock struct {
	configMock
}

func (c invalidConfigMock) Validate() error {
	return v.ValidationErrors{
		{Namespace: "Value", Err: "error"},
	}
}

func WriteConfigMockFile(t *testing.T) string {
	cfg := "value: test"
	cfgPath := path.Join(t.TempDir(), "cfg.yaml")

	err := ioutil.WriteFile(cfgPath, []byte(cfg), os.ModePerm)
	assert.NoError(t, err)

	return cfgPath
}

func TestValidate(t *testing.T) {
	assert.Empty(t, validateConfig(configMock{}), "validateConfig detected a non-existing error")
}

func TestValidate_InvalidConfig(t *testing.T) {
	expect := []error{NewValidationError("error", "Value")}
	assert.Equal(t, expect, validateConfig(invalidConfigMock{}))
}

func TestReadConfig(t *testing.T) {
	cfgPath := WriteConfigMockFile(t)

	cfg, err := readConfig(cfgPath, configMock{})
	assert.NoError(t, err)
	assert.Equal(t, "test", cfg.Value)
}

func TestReadConfig_NotExists(t *testing.T) {
	cfgPath := path.Join(t.TempDir(), "cfg.yaml")
	_, err := readConfig(cfgPath, configMock{})
	assert.ErrorContains(t, err, "does not exist")
}

func TestReadConfig_Invalid(t *testing.T) {
	cfgPath := path.Join(t.TempDir(), "cfg.yaml")
	_, err := readConfig(cfgPath, configMock{})
	assert.ErrorContains(t, err, "does not exist")
}

func TestReadConfigIfExists(t *testing.T) {
	cfgPath := WriteConfigMockFile(t)
	cfg, err := readConfigIfExists(cfgPath, configMock{})
	assert.NoError(t, err)
	assert.Equal(t, "test", cfg.Value)
}

func TestReadConfigIfExists_NotExists(t *testing.T) {
	cfgPath := path.Join(t.TempDir(), "cfg.yaml")
	cfg, err := readConfigIfExists(cfgPath, configMock{})
	assert.NoError(t, err)
	assert.Nil(t, cfg)
}
