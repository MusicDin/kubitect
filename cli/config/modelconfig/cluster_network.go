package modelconfig

import v "cli/validation"

type NetworkMode string

const (
	NAT    NetworkMode = "nat"
	ROUTE  NetworkMode = "route"
	BRIDGE NetworkMode = "bridge"
)

func (mode NetworkMode) Validate() error {
	return v.Var(mode, v.OneOf(NAT, ROUTE, BRIDGE))
}

type Network struct {
	CIDR    *CIDRv4        `yaml:"cidr"`
	Gateway *IPv4          `yaml:"gateway"`
	Mode    *NetworkMode   `yaml:"mode"`
	Bridge  *NetworkBridge `yaml:"bridge"`
}

func (n Network) Validate() error {

	//v.Required().Errorf("Property '{.Field}' is required when network mode is set to %s.", BRIDGE),
	return v.Struct(&n,
		v.Field(&n.CIDR, v.Required()),
		v.Field(&n.Gateway),
		v.Field(&n.Mode),
		v.Field(&n.Bridge, v.Required().When(n.Mode != nil && *n.Mode == BRIDGE).Errorf("Field '{.Field}' is required when network mode is set to '%v'.", BRIDGE)),
	)
}

type NetworkBridge string

func (br NetworkBridge) Validate() error {
	return v.Var(br,
		v.AlphaNumeric(),
		v.MaxLen(16),
	)
}
