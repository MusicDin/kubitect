package terraform

import (
	"cli/config/modelconfig"
	"cli/utils/defaults"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func clsPath(t *testing.T) string {
	clsPath, err := filepath.Abs("../../../../")
	assert.NoError(t, err)

	return filepath.Clean(clsPath)
}

func MockPKey(t *testing.T) modelconfig.File {
	path := filepath.Join(t.TempDir(), ".ssh/id_rsa")

	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	assert.NoError(t, err)

	_, err = os.Create(path)
	assert.NoError(t, err)

	return modelconfig.File(path)
}

func MockLocalHost(t *testing.T, name string, def bool) modelconfig.Host {
	host := modelconfig.Host{
		Name:    name,
		Default: def,
		Connection: modelconfig.Connection{
			Type: modelconfig.LOCAL,
		},
	}

	assert.NoError(t, defaults.Assign(&host).Validate())
	return host
}

func MockRemoteHost(t *testing.T, name string, def bool, verify bool) modelconfig.Host {
	host := modelconfig.Host{
		Name:    name,
		Default: def,
		Connection: modelconfig.Connection{
			Type: modelconfig.REMOTE,
			IP:   modelconfig.IPv4("192.168.113.42"),
			User: modelconfig.User("mocked-user"),
			SSH: modelconfig.ConnectionSSH{
				Keyfile: MockPKey(t),
				Verify:  verify,
			},
		},
	}

	assert.NoError(t, defaults.Assign(&host).Validate())
	return host
}

func TestHostUri_Empty(t *testing.T) {
	uri, err := hostUri(modelconfig.Host{})
	assert.NoError(t, err)
	assert.Equal(t, "qemu:///system", uri)
}

func TestHostUri_Local(t *testing.T) {
	h := MockLocalHost(t, "local", false)

	uri, err := hostUri(h)
	assert.NoError(t, err)
	assert.Equal(t, "qemu:///system", uri)
}

func TestHostUri_Remote(t *testing.T) {
	h := MockRemoteHost(t, "remote", false, false)
	pkey := h.Connection.SSH.Keyfile
	expected := fmt.Sprintf("qemu+ssh://mocked-user@192.168.113.42:22/system?keyfile=%s&no_verify=1", pkey)

	uri, err := hostUri(h)
	assert.NoError(t, err)
	assert.Equal(t, expected, uri)
}

func TestHostUri_Remote_Verify(t *testing.T) {
	h := MockRemoteHost(t, "remote", false, true)
	pkey := h.Connection.SSH.Keyfile
	expected := fmt.Sprintf("qemu+ssh://mocked-user@192.168.113.42:22/system?keyfile=%s", pkey)

	uri, err := hostUri(h)
	assert.NoError(t, err)
	assert.Equal(t, expected, uri)
}

func TestHostUri_NoHomeVar(t *testing.T) {
	home := os.Getenv("HOME")
	defer func() { os.Setenv("HOME", home) }()
	assert.NoError(t, os.Setenv("HOME", ""))

	h := MockRemoteHost(t, "remote", false, false)
	_, err := hostUri(h)
	assert.EqualError(t, err, "$HOME is not defined")
}

func TestIsDefault(t *testing.T) {
	lh := MockLocalHost(t, "local", false)
	rh := MockRemoteHost(t, "remote", true, false)

	hosts := []modelconfig.Host{lh, rh}

	def, err := defaultHost(hosts)
	assert.NoError(t, err)
	assert.Equal(t, rh, def)
}

func TestDefaultHost_NoDefaultHostSet(t *testing.T) {
	lh := MockLocalHost(t, "local", false)
	rh := MockRemoteHost(t, "remote", false, false)

	hosts := []modelconfig.Host{lh, rh}

	def, err := defaultHost(hosts)
	assert.NoError(t, err)
	assert.Equal(t, lh, def)
}

func TestDefaultHost_EmptyList(t *testing.T) {
	_, err := defaultHost([]modelconfig.Host{})
	assert.EqualError(t, err, "defaultHost: hosts list is empty")
}
