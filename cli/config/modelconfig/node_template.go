package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type OperatingSystem struct {
	Distro           *Distro                `yaml:"distro,omitempty"`
	NetworkInterface *NetworkInterface      `yaml:"networkInterface,omitempty"`
	Source           *OperatingSystemSource `yaml:"source,omitempty"`
}

func (s OperatingSystem) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Distro),
		validation.Field(&s.NetworkInterface), // TODO: depends on Distro
		validation.Field(&s.Source),           // TODO: depends on Distro
	)
}

type NodeTemplateSSH struct {
	AddToKnownHosts *bool   `yaml:"addToKnownHosts,omitempty"`
	PrivateKeyPath  *string `yaml:"privateKeyPath,omitempty"`
}

func (s NodeTemplateSSH) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.PrivateKeyPath, validation.When(s.PrivateKeyPath != nil, validation.By(PathExists))),
	)
}

type NodeTemplate struct {
	UpdateOnBoot *bool     `yaml:"updateOnBoot,omitempty"`
	User         *UserName `yaml:"user,omitempty"`

	DNS *[]IP            `yaml:"dns,omitempty"`
	OS  *OperatingSystem `yaml:"os"`

	SSH *NodeTemplateSSH `yaml:"ssh,omitempty"`
}

func (n NodeTemplate) Validate() error {
	return validation.ValidateStruct(&n,
		validation.Field(&n.DNS), // TODO: isValidIp for each?
		validation.Field(&n.User),
		validation.Field(&n.UpdateOnBoot),
		validation.Field(&n.SSH),
	)
}

type NetworkInterface string

func (n NetworkInterface) Validate() error {
	return validation.Validate(string(n), validation.Length(MinStringLength, 0))
}

type OperatingSystemSource string

func (s OperatingSystemSource) Validate() error {
	return validation.Validate(string(s), validation.Length(MinStringLength, 0))
}

type UserName string

func (n UserName) Validate() error {
	return validation.Validate(string(n), StringNotEmptyAlphaNumericMinus...)
}
