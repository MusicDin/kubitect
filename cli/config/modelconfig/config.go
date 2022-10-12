package modelconfig

import (
	v "cli/validation"
	"strings"
)

const (
	IP_IN_CIDR = "ipInCidr"
	VALID_HOST = "validHost"
)

type Config struct {
	Hosts      *[]Host     `yaml:"hosts"`
	Cluster    *Cluster    `yaml:"cluster"`
	Kubernetes *Kubernetes `yaml:"kubernetes"`
	Addons     *Addons     `yaml:"addons"`
}

func (c Config) Validate() error {
	defer v.ClearCustomValidators()

	c.regIPInCIDRValidator()
	c.regHostNameValidator()

	return v.Struct(&c,
		v.Field(&c.Hosts, v.MinLen(1).Error("At least {.Param} {.Field} must be configured.")),
		v.Field(&c.Cluster, v.Required().Error("Configuration must contain '{.Field}' section.")),
		v.Field(&c.Kubernetes, v.Required().Error("Configuration must contain '{.Field}' section.")),
		v.Field(&c.Addons, v.OmitEmpty()),
	)
}

// regIPInCIDRValidator registers a custom validator that checks whether
// an IP address is within the configured network CIDR.
func (c Config) regIPInCIDRValidator() {
	if c.Cluster != nil && c.Cluster.Network != nil && c.Cluster.Network.CIDR != nil {
		cidr := *c.Cluster.Network.CIDR
		v.RegisterCustomValidator(IP_IN_CIDR, v.IPInRange(string(cidr)))
	}
}

// regHostNameValidator registers a custom validator that checks whether
// a host with a given name has been configured.
func (c Config) regHostNameValidator() {
	if c.Hosts == nil {
		return
	}

	var hostNames []string

	for _, h := range *c.Hosts {
		if h.Name != nil {
			hostNames = append(hostNames, *h.Name)
		}
	}

	v.RegisterCustomValidator(VALID_HOST, v.OneOf(hostNames).Errorf("Field '{.Field}' must be a valid name of the configured host: [%s]", strings.Join(hostNames, "|")))
}
