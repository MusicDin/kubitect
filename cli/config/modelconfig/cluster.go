package modelconfig

import v "cli/utils/validation"

type Cluster struct {
	Name         *string      `yaml:"name"`
	Network      Network      `yaml:"network"`
	NodeTemplate NodeTemplate `yaml:"nodeTemplate"`
	Nodes        Nodes        `yaml:"nodes"`
}

func (c Cluster) Validate() error {
	return v.Struct(&c,
		v.Field(&c.Name, v.Required(), v.AlphaNumericHyp()),
		v.Field(&c.Network),
		v.Field(&c.Nodes, c.uniqueIpValidator(), c.uniqueMacValidator()),
		v.Field(&c.NodeTemplate),
	)
}

// uniqueIpValidator returns a validator that triggers an error if multiple nodes
// are assigned the same IP address.
func (c Cluster) uniqueIpValidator() v.Validator {
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

	return v.Fail().Errorf("IP address of each node instance (including VIP) must be unique. (duplicates: %v)", duplicates)
}

// uniqueMacValidator returns a validator that triggers an error if multiple nodes
// are assigned the same MAC address.
func (c Cluster) uniqueMacValidator() v.Validator {
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

	return v.Fail().Errorf("MAC address of each node instance must be unique. (duplicates: %v)", duplicates)
}
