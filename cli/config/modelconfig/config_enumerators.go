package modelconfig

import (
	"github.com/go-ozzo/ozzo-validation/v4"
)

type ConnectionType string

func (ct ConnectionType) validate() error {
	return validation.Validate(&ct, validation.In(ConnectionTypeList))
}

var ConnectionTypeList = [...]ConnectionType{
	local_connection,
	remote_connection,
}

const (
	local_connection  ConnectionType = "local"
	remote_connection                = "remote"
)

type OperatingSystem string

func (os OperatingSystem) validate() error {
	return validation.Validate(&os, validation.In(OperatingSystemList))
}

var OperatingSystemList = [...]OperatingSystem{
	ubuntu,
	ubuntu22,
	ubuntu20,
	debian,
	debian11,
	custom,
}

const (
	ubuntu   OperatingSystem = "ubuntu"
	ubuntu20                 = "ubuntu20"
	ubuntu22                 = "ubuntu22"
	debian                   = "debian"
	debian11                 = "debian11"
	custom                   = "custom"
)

type NetworkMode string

func (nm NetworkMode) Validate() error {
	return validation.Validate(&nm, validation.In(NetworkModeList))
}

var NetworkModeList = [...]NetworkMode{
	nat_network_mode,
	remote_network_mode,
}

const (
	nat_network_mode    NetworkMode = "nat"
	remote_network_mode             = "remote"
)

type PortForwardTarget string

func (pft PortForwardTarget) Validate() error {
	return validation.Validate(&pft, validation.In(PortForwardTargetList))
}

var PortForwardTargetList = [...]PortForwardTarget{
	workers,
	masters,
	all,
}

const (
	workers PortForwardTarget = "workers"
	masters                   = "masters"
	all                       = "all"
)
