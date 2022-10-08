package modelconfig

import (
	v "cli/validation"
)

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
		// v.Field(&d.Pool), // TODO: Cross validate poolName - pool with that name must exist (if not "main")
		v.Field(&d.Size),
	)
}

type Version string

func (ver Version) Validate() error {
	return v.Var(ver, v.VSemVer())
}
