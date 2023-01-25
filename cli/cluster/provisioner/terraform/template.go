package terraform

import (
	"cli/config/modelconfig"
	"cli/lib/template"
	"fmt"
	"os"
	"path"
	"strings"
)

type MainTemplate struct {
	Hosts []modelconfig.Host
}

func NewMainTemplate(hosts []modelconfig.Host) MainTemplate {
	return MainTemplate{
		Hosts: hosts,
	}
}

func (t MainTemplate) Name() string {
	return "main.tf"
}

func (t MainTemplate) Functions() map[string]interface{} {
	return map[string]interface{}{
		"hostUri":             hostUri,
		"defaultHost":         defaultHost,
		"hostMainResPoolPath": hostMainResPoolPath,
	}
}

// Write creates main.tf file from template.
func (t MainTemplate) Write(clusterPath string) error {
	srcPath := path.Join(clusterPath, "terraform", "main.tf.tpl")
	dstPath := path.Join(clusterPath, "terraform", "main.tf")

	return template.WriteFrom(t, srcPath, dstPath)
}

// defaultHost returns default host from a given list of hosts.
func defaultHost(hosts []modelconfig.Host) (modelconfig.Host, error) {
	if hosts == nil || len(hosts) == 0 {
		return modelconfig.Host{}, fmt.Errorf("defaultHost: hosts list is empty")
	}

	for _, h := range hosts {
		if h.Default != nil && *h.Default {
			return h, nil
		}
	}

	return hosts[0], nil
}

// hostMainResPoolPath returns main resource pool path (MRPP) of the host.
// If MRPP is nil, a default MRRP is returned.
func hostMainResPoolPath(host modelconfig.Host) string {
	if host.MainResourcePoolPath != nil {
		return *host.MainResourcePoolPath
	}

	return "/var/lib/libvirt/images/"
}

// hostUri returns URI of a given host.
func hostUri(host modelconfig.Host) (string, error) {
	typ := host.Connection.Type

	if typ == nil || *typ == modelconfig.LOCALHOST || *typ == modelconfig.LOCAL {
		return "qemu:///system", nil
	}

	ip := *host.Connection.IP
	user := *host.Connection.User
	pkey := "~/.ssh/id_rsa"
	port := 22
	verify := "&no_verify=1"

	if host.Connection.SSH.Port != nil {
		port = int(*host.Connection.SSH.Port)
	}

	if host.Connection.SSH.Keyfile != nil {
		pkey = fmt.Sprintf("%v", *host.Connection.SSH.Keyfile)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	pkey = strings.Replace(pkey, "~", homeDir, 1)

	if host.Connection.SSH.Verify != nil && *host.Connection.SSH.Verify {
		verify = ""
	}

	uri := fmt.Sprintf("qemu+ssh://%s@%s:%d/system?keyfile=%s%s", user, ip, port, pkey, verify)
	return uri, nil
}
