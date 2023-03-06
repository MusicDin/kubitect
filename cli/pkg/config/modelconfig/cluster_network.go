package modelconfig

import (
	"github.com/MusicDin/kubitect/cli/pkg/utils/defaults"
	"github.com/MusicDin/kubitect/cli/pkg/utils/validation"
)

type NetworkMode string

const (
	NAT    NetworkMode = "nat"
	ROUTE  NetworkMode = "route"
	BRIDGE NetworkMode = "bridge"
)

func (mode NetworkMode) Validate() error {
	return validation.Var(mode, validation.OneOf(NAT, ROUTE, BRIDGE))
}

type Network struct {
	CIDR    CIDRv4        `yaml:"cidr"`
	Gateway *IPv4         `yaml:"gateway,omitempty"`
	Mode    NetworkMode   `yaml:"mode"`
	Bridge  NetworkBridge `yaml:"bridge,omitempty"`
}

func (n Network) Validate() error {
	return validation.Struct(&n,
		validation.Field(&n.CIDR, validation.NotEmpty()),
		validation.Field(&n.Gateway),
		validation.Field(&n.Mode),
		validation.Field(&n.Bridge, validation.NotEmpty().When(n.Mode == BRIDGE).Errorf("Field '{.Field}' is required when network mode is set to '%v'.", BRIDGE)),
	)
}

func (n *Network) SetDefaults() {
	n.Mode = defaults.Default(n.Mode, NAT)
}

type NetworkBridge string

func (br NetworkBridge) Validate() error {
	return validation.Var(br,
		validation.OmitEmpty(),
		validation.AlphaNumeric(),
		validation.MaxLen(16),
	)
}
