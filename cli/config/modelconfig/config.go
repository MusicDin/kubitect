package modelconfig

import (
	v "cli/validation"
	"strings"
)

// Keys of custom validators
const (
	IP_IN_CIDR = "ipInCidr"
	VALID_HOST = "validHost"
	VALID_POOL = "validPool"
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

// regHostNameValidator registers a custom validator that checks whether
// a host with a given name has been configured.
func getPoolNameValidator(hostName *string) v.Validator {

	c, ok := v.TopParent().(*Config)
	if !ok || c == nil {
		return v.Validator{}
	}

	if c.Hosts == nil || len(*c.Hosts) == 0 {
		return v.Validator{}
	}

	// By default, the first host in a list is a default host.
	host := (*c.Hosts)[0]

	for _, h := range *c.Hosts {
		if h.Default != nil && *h.Default {
			host = h
		}

		if hostName == nil || h.Name == nil {
			continue
		}

		if *h.Name == *hostName {
			host = h
			break
		}
	}

	if host.Name == nil {
		// Ignore, because in such case an error is already triggered for a host.
		return v.Validator{}
	}

	pools := host.DataResourcePools

	if pools == nil || len(*pools) == 0 {
		return v.Fail().Errorf("Field '{.Field}' points to a data resource pool, but matching host '%s' has none configured.", host)
	}

	var poolNames []string

	for _, p := range *host.DataResourcePools {
		if p.Name != nil {
			poolNames = append(poolNames, *p.Name)
		}
	}

	return v.OneOf(poolNames).Errorf("Field '{.Field}' must point to one of the pools configured on a matching host '%s': [%s]", *host.Name, strings.Join(poolNames, "|"))
}
