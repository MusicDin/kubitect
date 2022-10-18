package modelconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint8(t *testing.T) {
	assert.Error(t, Uint8(-1).Validate())
	assert.Error(t, Uint8(256).Validate())
	assert.NoError(t, Uint8(0).Validate())
	assert.NoError(t, Uint8(255).Validate())
}

func TestGB(t *testing.T) {
	assert.Error(t, GB(0).Validate())
	assert.NoError(t, GB(1).Validate())
}

func TestVCpu(t *testing.T) {
	assert.Error(t, VCpu(0).Validate())
	assert.NoError(t, VCpu(1).Validate())
}

func TestPort(t *testing.T) {
	assert.Error(t, Port(0).Validate())
	assert.Error(t, Port(65536).Validate())
	assert.NoError(t, Port(80).Validate())
}

func TestIP(t *testing.T) {
	assert.Error(t, IP("192.168.113.266").Validate())
	assert.NoError(t, IP("192.168.113.20").Validate())
	assert.NoError(t, IP("2001:db8:3333:4444:5555:6666:7777:8888").Validate())
	assert.NoError(t, IP("2001:db8::8888").Validate())
}

func TestIPv4(t *testing.T) {
	assert.Error(t, IPv4("192.168.113.266").Validate())
	assert.NoError(t, IPv4("192.168.113.20").Validate())
	assert.Error(t, IPv4("2001:db8:3333:4444:5555:6666:7777:8888").Validate())
	assert.Error(t, IPv4("2001:db8::8888").Validate())
}

func TestCIDRv4(t *testing.T) {
	assert.Error(t, CIDRv4("192.168.113.0").Validate())
	assert.Error(t, CIDRv4("192.168.113.0/33").Validate())
	assert.NoError(t, CIDRv4("192.168.113.0/20").Validate())
	assert.Error(t, CIDRv4("2001:db8::8888/64").Validate())
}

func TestMAC(t *testing.T) {
	assert.Error(t, MAC("AA:BB::FF").Validate())
	assert.NoError(t, MAC("AA:BB:CC:DD:EE:FF").Validate())
}

func TestUser(t *testing.T) {
	assert.Error(t, User("").Validate())
	assert.Error(t, User(".").Validate())
	assert.Error(t, User(" ").Validate())
	assert.NoError(t, User("user").Validate())
	assert.NoError(t, User("UsEr").Validate())
	assert.NoError(t, User("_UsEr_").Validate())
	assert.NoError(t, User("_-UsEr-_").Validate())
}

func TestFile(t *testing.T) {
	assert.Error(t, File("").Validate())
	assert.Error(t, File("./non-existing-file").Validate())
	assert.NoError(t, File("./common_test.go").Validate())
}

func TestURL(t *testing.T) {
	assert.Error(t, URL("kubitect.io").Validate())
	assert.NoError(t, URL("https://kubitect.io").Validate())
}

func TestTaint(t *testing.T) {
	assert.Error(t, Taint("").Validate())
	assert.NoError(t, Taint("taint").Validate())
}

func TestLabels(t *testing.T) {
	assert.NoError(t, Labels{}.Validate())
}

func TestDataDisk(t *testing.T) {
	str := "test"
	size := GB(5)

	dd := DataDisk{
		Name: &str,
		Pool: &str,
		Size: &(size),
	}

	assert.NoError(t, dd.Validate())
	assert.ErrorContains(t, DataDisk{Name: nil}.Validate(), "Field 'size' is required.")
	assert.ErrorContains(t, DataDisk{Name: nil}.Validate(), "Field 'name' is required.")
}

func TestVersion(t *testing.T) {
	assert.Error(t, Version("1.2.3").Validate())
	assert.Error(t, Version("v1.2").Validate())
	assert.Error(t, Version("v1.2.").Validate())
	assert.NoError(t, Version("v1.2.3").Validate())
	assert.Error(t, Version("master").Validate())
}

func TestMasterVersion(t *testing.T) {
	assert.Error(t, MasterVersion("1.2.3").Validate())
	assert.Error(t, MasterVersion("v1.2").Validate())
	assert.Error(t, MasterVersion("v1.2.").Validate())
	assert.NoError(t, MasterVersion("v1.2.3").Validate())
	assert.NoError(t, MasterVersion("master").Validate())
}
