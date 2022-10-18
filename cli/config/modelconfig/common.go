package modelconfig

import (
	v "cli/validation"
)

// Uint8 is intentionally set to int to avoid panic if value is set
// outside the uint8 size.
//
// For example, if LB priority is set to -1, raising a custom error
// is not possible since go will panic when converting -1 to uint8.
type Uint8 int

func (u Uint8) Validate() error {
	return v.Var(u, v.Min(0), v.Max(255))
}

type GB int

func (size GB) Validate() error {
	return v.Var(size, v.Min(1))
}

type VCpu int

func (s VCpu) Validate() error {
	return v.Var(s, v.Min(1))
}

type Port int

func (p Port) Validate() error {
	return v.Var(p, v.Min(1), v.Max(65535))
}

type IP string

func (ip IP) Validate() error {
	return v.Var(ip, v.IP())
}

type IPv4 string

func (ip IPv4) Validate() error {
	return v.Var(ip, v.IPv4())
}

type CIDRv4 string

func (cidr CIDRv4) Validate() error {
	return v.Var(cidr, v.CIDRv4())
}

type MAC string

func (mac MAC) Validate() error {
	return v.Var(mac, v.MAC())
}

type User string

func (u User) Validate() error {
	return v.Var(u, v.MinLen(1), v.AlphaNumericHypUS())
}

type File string

func (f File) Validate() error {
	return v.Var(f, v.FileExists())
}

type URL string

func (u URL) Validate() error {
	return v.Var(u, v.URL())
}

type Taint string

func (t Taint) Validate() error {
	return v.Var(t, v.Min(1))
}

type Labels map[string]string

func (l Labels) Validate() error {
	return v.Var(l, v.Required()) // TODO: Validate MAP
}

type DataDisk struct {
	Name *string `yaml:"name" opt:",id"`
	Pool *string `yaml:"pool"`
	Size *GB     `yaml:"size"`
}

func (d DataDisk) Validate() error {
	return v.Struct(&d,
		v.Field(&d.Name, v.Required(), v.AlphaNumericHyp()),
		v.Field(&d.Pool, v.OmitEmpty(), v.Skip().When(d.Pool != nil && *d.Pool == "main"), v.Custom(VALID_POOL)),
		v.Field(&d.Size, v.Required()),
	)
}

type Version string

func (ver Version) Validate() error {
	return v.Var(ver, v.VSemVer())
}

type MasterVersion string

func (ver MasterVersion) Validate() error {
	return v.Var(ver, v.Skip().When(ver == "master"), v.VSemVer())
}
