package modelinfra

import (
	"cli/config/modelconfig"
	v "cli/utils/validation"
)

type Config struct {
	Nodes modelconfig.Nodes `yaml:"nodes"`
}

func (c Config) Validate() error {
	return v.Struct(&c,
		v.Field(&c.Nodes, c.uniqueIpValidator(), c.uniqueMacValidator(), v.Skip()),
	)
}

// uniqueIpValidator returns a validator that triggers an error if multiple nodes
// are assigned the same IP address.
func (c Config) uniqueIpValidator() v.Validator {
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
		return v.None
	}

	return v.Fail().Errorf("Duplicate IPs detected in the provisioned infrastructure. (duplicates: %v)", duplicates)
}

// uniqueMacValidator returns a validator that triggers an error if multiple nodes
// are assigned the same MAC address.
func (c Config) uniqueMacValidator() v.Validator {
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
		return v.None
	}

	return v.Fail().Errorf("Duplicate MAC addresses detected in the provisioned infrastructure. (duplicates: %v)", duplicates)
}
