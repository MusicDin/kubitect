package terraform

import (
	"cli/config/modelconfig"
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
	typ := modelconfig.LOCAL

	host := modelconfig.Host{
		Name:    &name,
		Default: &def,
		Connection: modelconfig.Connection{
			Type: &typ,
		},
	}

	assert.NoError(t, host.Validate())
	return host
}

func MockRemoteHost(t *testing.T, name string, def bool, verify bool) modelconfig.Host {
	typ := modelconfig.REMOTE
	ip := modelconfig.IPv4("192.168.113.42")
	user := modelconfig.User("mocked-user")
	sshPKey := MockPKey(t)
	sshPort := modelconfig.Port(42)

	host := modelconfig.Host{
		Name:    &name,
		Default: &def,
		Connection: modelconfig.Connection{
			Type: &typ,
			IP:   &ip,
			User: &user,
			SSH: modelconfig.ConnectionSSH{
				Keyfile: &sshPKey,
				Port:    &sshPort,
				Verify:  &verify,
			},
		},
	}

	assert.NoError(t, host.Validate())
	return host
}

func TestNewTerraformProvisioner(t *testing.T) {
	hosts := []modelconfig.Host{
		MockLocalHost(t, "test1", false),
		MockLocalHost(t, "test2", true),
		MockRemoteHost(t, "test3", false, false),
	}

	_, err := NewTerraformProvisioner(clsPath(t), "shared/path", true, hosts)
	assert.NoError(t, err)
}

func TestNewTerraformProvisioner_InvalidHosts(t *testing.T) {
	_, err := NewTerraformProvisioner(clsPath(t), "shared/path", true, nil)
	assert.ErrorContains(t, err, "hosts list is empty")
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
	expected := fmt.Sprintf("qemu+ssh://mocked-user@192.168.113.42:42/system?keyfile=%s&no_verify=1", *pkey)

	uri, err := hostUri(h)
	assert.NoError(t, err)
	assert.Equal(t, expected, uri)
}

func TestHostUri_Remote_Verify(t *testing.T) {
	h := MockRemoteHost(t, "remote", false, true)
	pkey := h.Connection.SSH.Keyfile
	expected := fmt.Sprintf("qemu+ssh://mocked-user@192.168.113.42:42/system?keyfile=%s", *pkey)

	uri, err := hostUri(h)
	assert.NoError(t, err)
	assert.Equal(t, expected, uri)
}

func TestHostUri_NoHomeVar(t *testing.T) {
	assert.NoError(t, os.Setenv("HOME", ""))

	h := MockRemoteHost(t, "remote", false, false)
	_, err := hostUri(h)
	assert.EqualError(t, err, "$HOME is not defined")
}

func TestHostMainResPoolPath(t *testing.T) {
	path := "test"

	h := modelconfig.Host{}
	h.MainResourcePoolPath = &path
	assert.Equal(t, "test", hostMainResPoolPath(h))
}

func TestHostMainResPoolPath_Default(t *testing.T) {
	assert.Equal(t, "/var/lib/libvirt/images/", hostMainResPoolPath(modelconfig.Host{}))
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
