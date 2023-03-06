package modelinfra

import (
	"github.com/MusicDin/kubitect/cli/pkg/config/modelconfig"
	"github.com/MusicDin/kubitect/cli/pkg/utils/validation"
)

type Config struct {
	Nodes modelconfig.Nodes `yaml:"nodes"`
}

func (c Config) Validate() error {
	return validation.Struct(&c,
		validation.Field(&c.Nodes, c.uniqueIpValidator(), c.uniqueMacValidator(), validation.Skip()),
	)
}

// uniqueIpValidator returns a validator that triggers an error if multiple nodes
// are assigned the same IP address.
func (c Config) uniqueIpValidator() validation.Validator {
	var duplicates []string

	ips := c.Nodes.IPs()

	for i := 0; i < len(ips); i++ {
		for j := i + 1; j < len(ips); j++ {
			if ips[i] == ips[j] {
				duplicates = append(duplicates, ips[i])
			}
		}
	}

	if len(duplicates) == 0 {
		return validation.None
	}

	return validation.Fail().Errorf("Duplicate IPs detected in the provisioned infrastructure. (duplicates: %v)", duplicates)
}

// uniqueMacValidator returns a validator that triggers an error if multiple nodes
// are assigned the same MAC address.
func (c Config) uniqueMacValidator() validation.Validator {
	var duplicates []string

	macs := c.Nodes.MACs()

	for i := 0; i < len(macs); i++ {
		for j := i + 1; j < len(macs); j++ {
			if macs[i] == macs[j] {
				duplicates = append(duplicates, macs[i])
			}
		}
	}

	if len(duplicates) == 0 {
		return validation.None
	}

	return validation.Fail().Errorf("Duplicate MAC addresses detected in the provisioned infrastructure. (duplicates: %v)", duplicates)
}
