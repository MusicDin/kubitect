package modelconfig

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Cluster struct {
	Name         *ClusterName  `yaml:"name,omitempty"`
	Network      *Network      `yaml:"network,omitempty"`
	NodeTemplate *NodeTemplate `yaml:"nodeTemplate,omitempty"`
	Nodes        *Nodes        `yaml:"nodes,omitempty"`
}

func (c Cluster) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name),
		validation.Field(&c.Network),
		validation.Field(&c.Nodes),
		validation.Field(&c.NodeTemplate),
	)
}

type ClusterName string

func (n ClusterName) Validate() error {
	return validation.Validate(string(n), StringNotEmptyAlphaNumericMinus...)
}
