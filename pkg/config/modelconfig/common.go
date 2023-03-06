package modelconfig

import (
	"github.com/MusicDin/kubitect/pkg/utils/validation"
)

// Uint8 is intentionally set to int to avoid panic if value is set
// outside the uint8 size.
//
// For example, if LB priority is set to -1, raising a custom error
// is not possible since go will panic when converting -1 to uint8.
type Uint8 int

func (u Uint8) Validate() error {
	return validation.Var(u, validation.Min(0), validation.Max(255))
}

type GB int

func (size GB) Validate() error {
	return validation.Var(size, validation.Min(1))
}

type VCpu int

func (s VCpu) Validate() error {
	return validation.Var(s, validation.Min(1))
}

type Port int

func (p Port) Validate() error {
	return validation.Var(p, validation.Min(1), validation.Max(65535))
}

type IP string

func (ip IP) Validate() error {
	return validation.Var(ip, validation.IP())
}

type IPv4 string

func (ip IPv4) Validate() error {
	return validation.Var(ip, validation.IPv4())
}

type CIDRv4 string

func (cidr CIDRv4) Validate() error {
	return validation.Var(cidr, validation.CIDRv4())
}

type MAC string

func (mac MAC) Validate() error {
	return validation.Var(mac, validation.MAC())
}

type User string

func (u User) Validate() error {
	return validation.Var(u, validation.MinLen(1), validation.AlphaNumericHypUS())
}

type File string

func (f File) Validate() error {
	return validation.Var(f, validation.FileExists())
}

type URL string

func (u URL) Validate() error {
	return validation.Var(u, validation.URL())
}

type Taint string

func (t Taint) Validate() error {
	return validation.Var(t, validation.Min(1))
}

type Labels map[string]string

func (l Labels) Validate() error {
	return validation.Var(l, validation.Required()) // TODO: Validate MAP
}

type DataDisk struct {
	Name string `yaml:"name" opt:",id"`
	Pool string `yaml:"pool"`
	Size GB     `yaml:"size"`
}

func (d DataDisk) Validate() error {
	return validation.Struct(&d,
		validation.Field(&d.Name, validation.NotEmpty(), validation.AlphaNumericHyp()),
		validation.Field(&d.Pool, validation.OmitEmpty(), validation.Skip().When(d.Pool == "main"), validation.Custom(VALID_POOL)),
		validation.Field(&d.Size, validation.NotEmpty()),
	)
}

type Version string

func (ver Version) Validate() error {
	return validation.Var(ver, validation.VSemVer())
}

type MasterVersion string

func (ver MasterVersion) Validate() error {
	return validation.Var(ver, validation.Skip().When(ver == "master"), validation.VSemVer())
}
