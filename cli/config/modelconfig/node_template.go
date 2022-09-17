package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type NodeTemplate struct {
	UpdateOnBoot *bool     `yaml:"updateOnBoot"`
	User         *UserName `yaml:"user"`

	DNS []IP `yaml:"dns"`
	OS  struct {
		Distro           *OperatingSystem       `yaml:"distro"`
		NetworkInterface *NetworkInterface      `yaml:"networkInterface"`
		Source           *OperatingSystemSource `yaml:"source"`
	} `yaml:"os"`

	SSH struct {
		AddToKnownHosts bool    `yaml:"addToKnownHosts"`
		PrivateKeyPath  *string `yaml:"privateKeyPath"`
	} `yaml:"ssh"`
}

func (n NodeTemplate) Validate() error {
	return validation.ValidateStruct(&n,
		validation.Field(&n.DNS), // TODO: isValidIp for each?
		validation.Field(n.User),
		validation.Field(n.UpdateOnBoot),
		validation.Field(n.OS.Distro),
		validation.Field(n.OS.NetworkInterface), // TODO: depends on Distro
		validation.Field(n.OS.Source),           // TODO: depends on Distro
		validation.Field(n.SSH.AddToKnownHosts, validation.When(n.SSH.PrivateKeyPath != nil, validation.By(PathExists))),
	)
}

type NetworkInterface string

func (n NetworkInterface) Validate() error {
	return validation.Validate(&n, validation.Min(1))
}

type OperatingSystemSource string

func (s OperatingSystemSource) Validate() error {
	return validation.Validate(&s, validation.Min(1))
}

type UserName string

func (n UserName) Validate() error {
	return validation.Validate(&n, StringNotEmptyAlphaNumeric...)
}
