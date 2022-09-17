package modelconfig

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Cluster struct {
	Name         *ClusterName  `yaml:"name"`
	Network      *Network      `yaml:"network"`
	Nodes        *Nodes        `yaml:"nodes"`
	NodeTemplate *NodeTemplate `yaml:"nodeTemplate"`
}

func (c Cluster) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(c.Network),
		validation.Field(c.Name),
		validation.Field(c.Nodes),
		validation.Field(c.NodeTemplate),
	)
}

type ClusterName string

func (n ClusterName) Validate() error {
	return validation.Validate(&n, StringNotEmptyAlphaNumeric...)
}
