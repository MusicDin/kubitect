package modelconfig

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type MB uint

func (m MB) Validate() error {
	return validation.Validate(int(m), validation.Min(1))
}

type Port uint16

func (p Port) Validate() error {
	return validation.Validate(int(p), validation.Min(1), validation.Max(65535))
}

type CpuSize uint

func (s CpuSize) Validate() error {
	return validation.Validate(int(s), validation.Min(1))
}

type HostName string

func (n HostName) Validate() error {
	return validation.Validate(string(n), StringNotEmptyAlphaNumericMinus...) // Is valid Hostname?
}

type IP string

func (ip IP) Validate() error {
	return validation.Validate(string(ip), is.IP)
}

type MAC string

func (mac MAC) Validate() error {
	return validation.Validate(string(mac), is.MAC)
}

type InstanceId string

func (i InstanceId) Validate() error {
	return validation.Validate(string(i), StringNotEmptyAlphaNumericMinus...)
}

type ResourcePath string

func (r ResourcePath) Validate() error {
	return validation.Validate(string(r))
}
