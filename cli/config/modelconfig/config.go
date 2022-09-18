package modelconfig

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	MinHostsLength = 1
	MaxHostsLength = 0
)

type Config struct {
	Hosts   *[]Host  `yaml:"hosts,omitempty"`
	Cluster *Cluster `yaml:"cluster,omitempty"`
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Cluster),
		validation.Field(&c.Hosts, validation.Length(MinHostsLength, MaxHostsLength)),
	)
}
