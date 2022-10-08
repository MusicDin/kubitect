package modelconfig

import v "cli/validation"

type Cluster struct {
	Name         *string       `yaml:"name"`
	Network      *Network      `yaml:"network"`
	NodeTemplate *NodeTemplate `yaml:"nodeTemplate"`
	Nodes        *Nodes        `yaml:"nodes"`
}

func (c Cluster) Validate() error {
	return v.Struct(&c,
		v.Field(&c.Name, v.Required(), v.AlphaNumericHyp()),
		v.Field(&c.Network),
		v.Field(&c.Nodes),
		v.Field(&c.NodeTemplate),
	)
}
