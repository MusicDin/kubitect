package modelinfra

import v "cli/lib/validation"

type Config struct {
	Cluster Cluster `yaml:"cluster"`
}

func (c Config) Validate() error {
	return v.Struct(&c,
		v.Field(&c.Cluster, v.NotEmpty().Error("Terraform produced invalid output.")),
	)
}
