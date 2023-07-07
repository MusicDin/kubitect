package config

import (
	"github.com/MusicDin/kubitect/pkg/utils/defaults"
	v "github.com/MusicDin/kubitect/pkg/utils/validation"
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
	Mode     NetworkMode   `yaml:"mode"`
	Bridge   NetworkBridge `yaml:"bridge,omitempty"`
	CIDR     CIDRv4        `yaml:"cidr"`
	CIDR6    CIDRv6        `yaml:"cidr6,omitempty"`
	Gateway  *IPv4         `yaml:"gateway,omitempty"`
	Gateway6 *IPv6         `yaml:"gateway6,omitempty"`
}

func (n Network) Validate() error {
	return v.Struct(&n,
		v.Field(&n.Mode),
		v.Field(&n.Bridge, v.NotEmpty().When(n.Mode == BRIDGE).Errorf("Field '{.Field}' is required when network mode is set to '%v'.", BRIDGE)),
		v.Field(&n.CIDR, v.NotEmpty()),
		v.Field(&n.CIDR6, v.OmitEmpty()),
		v.Field(&n.Gateway),
		v.Field(&n.Gateway6),
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
