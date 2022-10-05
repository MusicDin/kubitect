package modelconfig

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

type ConnectionType string

func (ct ConnectionType) validate() error {
	if err := validation.In(ConnectionTypeList...).Validate(ct); err != nil {
		return err
	}
	return nil
}

var ConnectionTypeList = []interface{}{
	local_connection,
	remote_connection,
}

const (
	local_connection  ConnectionType = "local"
	remote_connection ConnectionType = "remote"
)

type Distro string

func (os Distro) validate() error {
	if err := validation.In(OperatingSystemList...).Validate(os); err != nil {
		return err
	}
	return nil
}

var OperatingSystemList = []interface{}{
	ubuntu,
	ubuntu22,
	ubuntu20,
	debian,
	debian11,
	custom,
}

const (
	ubuntu   Distro = "ubuntu"
	ubuntu20 Distro = "ubuntu20"
	ubuntu22 Distro = "ubuntu22"
	debian   Distro = "debian"
	debian11 Distro = "debian11"
	custom   Distro = "custom"
)

type NetworkMode string

func (nm NetworkMode) Validate() error {
	if err := validation.In(NetworkModeList...).Validate(nm); err != nil {
		return err
	}
	return nil
}

var NetworkModeList = []interface{}{
	nat_network_mode,
	remote_network_mode,
	bridge_network_mode,
}

var (
	nat_network_mode    NetworkMode = "nat"
	remote_network_mode NetworkMode = "remote"
	bridge_network_mode NetworkMode = "bridge"
)

type PortForwardTarget string

func (pft PortForwardTarget) Validate() error {

	if err := validation.In(PortForwardTargetList...).Validate(pft); err != nil {
		return err
	}
	return nil
}

var PortForwardTargetList = []interface{}{
	workers,
	masters,
	all,
}

const (
	workers PortForwardTarget = "workers"
	masters PortForwardTarget = "masters"
	all     PortForwardTarget = "all"
)
