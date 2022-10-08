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
	Bridge  *NetworkBridge  `yaml:"bridge"`
	CIDR    *NetworkCIDR    `yaml:"cidr"`
	Gateway *NetworkGateway `yaml:"gateway"`
	Mode    *NetworkMode    `yaml:"mode"`
}

func (n Network) Validate() error {
	return v.Struct(&n,
		v.Field(&n.CIDR, v.Required()),
		v.Field(&n.Bridge, v.OmitEmpty().When(*n.Mode != BRIDGE)),
		v.Field(&n.Gateway),
		v.Field(&n.Mode, v.OmitEmpty()),
	)
}

type NetworkGateway string

func (gw NetworkGateway) Validate() error {
	return v.Var(gw, v.OmitEmpty(), v.IPv4())
}

type NetworkCIDR string

func (cidr NetworkCIDR) Validate() error {
	return v.Var(cidr, v.Required(), v.CIDRv4())
}

type NetworkBridge string

func (cidr NetworkBridge) Validate() error {
	return v.Var(cidr,
		v.Required().Errorf("Property '{.Field}' is required when network mode is set to %s.", BRIDGE),
		v.AlphaNumeric(),
	)
}
