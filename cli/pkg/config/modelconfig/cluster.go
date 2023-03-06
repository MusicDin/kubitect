package modelconfig

import (
	"github.com/MusicDin/kubitect/cli/pkg/utils/validation"
)

type Cluster struct {
	Name         string       `yaml:"name"`
	Network      Network      `yaml:"network"`
	NodeTemplate NodeTemplate `yaml:"nodeTemplate"`
	Nodes        Nodes        `yaml:"nodes"`
}

func (c Cluster) Validate() error {
	return validation.Struct(&c,
		validation.Field(&c.Name, validation.NotEmpty(), validation.AlphaNumericHyp()),
		validation.Field(&c.Network),
		validation.Field(&c.Nodes, c.uniqueIpValidator(), c.uniqueMacValidator()),
		validation.Field(&c.NodeTemplate),
	)
}

// uniqueIpValidator returns a validator that triggers an error if multiple nodes
// are assigned the same IP address.
func (c Cluster) uniqueIpValidator() validation.Validator {
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

	return validation.Fail().Errorf("IP address of each node instance (including VIP) must be unique. (duplicates: %v)", duplicates)
}

// uniqueMacValidator returns a validator that triggers an error if multiple nodes
// are assigned the same MAC address.
func (c Cluster) uniqueMacValidator() validation.Validator {
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

	return validation.Fail().Errorf("MAC address of each node instance must be unique. (duplicates: %v)", duplicates)
}
