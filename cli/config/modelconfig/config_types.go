package modelconfig

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type MB uint
type Port uint16
type CpuSize uint
type HostName string

type IP string

func (ip IP) Validate() error {
	return validation.Validate(&ip, is.IP)
}

type MAC string

func (mac MAC) Validate() error {
	return validation.Validate(&mac, is.MAC)
}

type LabelKey string
type Label string // TODO: Check if correct type
type Taint string

type SSHKeyPath string

type UserString string

type ResourcePath string

type ForwardPortName string

type LoadBalancerId uint

type DiskName string
type PoolName string
type DiskSize uint

type InstanceId uint
type Priority uint

type Bridge string
type CIDR string
type Gateway string
type Mode string

type UserName string
