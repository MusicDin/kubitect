package modelconfig

import (
	"cli/utils/defaults"
	v "cli/utils/validation"
)

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
	CIDR    CIDRv4        `yaml:"cidr"`
	Gateway *IPv4         `yaml:"gateway"`
	Mode    NetworkMode   `yaml:"mode"`
	Bridge  NetworkBridge `yaml:"bridge"`
}

func (n Network) Validate() error {
	return v.Struct(&n,
		v.Field(&n.CIDR, v.NotEmpty()),
		v.Field(&n.Gateway),
		v.Field(&n.Mode),
		v.Field(&n.Bridge, v.NotEmpty().When(n.Mode == BRIDGE).Errorf("Field '{.Field}' is required when network mode is set to '%v'.", BRIDGE)),
	)
}

func (n *Network) SetDefaults() {
	n.Mode = defaults.Default(n.Mode, NAT)
}

type NetworkBridge string

func (br NetworkBridge) Validate() error {
	return v.Var(br,
		v.OmitEmpty(),
		v.AlphaNumeric(),
		v.MaxLen(16),
	)
}
