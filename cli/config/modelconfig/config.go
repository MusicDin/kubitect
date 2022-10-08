package modelconfig

import v "cli/validation"

type Config struct {
	Hosts      *[]Host     `yaml:"hosts"`
	Cluster    *Cluster    `yaml:"cluster"`
	Kubernetes *Kubernetes `yaml:"kubernetes"`
	Addons     *Addons     `yaml:"addons"`
}

func (c Config) Validate() error {
	return v.Struct(&c,
		v.Field(&c.Hosts, v.MinLen(1).Error("At least {.Param} {.Field} must be configured.")),
		v.Field(&c.Cluster, v.Required().Error("Configuration must contain '{.Field}' section.")),
		v.Field(&c.Kubernetes, v.Required().Error("Configuration must contain '{.Field}' section.")),
		v.Field(&c.Addons, v.OmitEmpty()),
	)
}
