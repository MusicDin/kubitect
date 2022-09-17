package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type NodeTemplate struct {
	UpdateOnBoot *bool
	User         *UserName

	DNS []IP
	OS  struct {
		Distro           *OperatingSystem
		NetworkInterface *NetworkInterface
		Source           *OperatingSystemSource
	}

	SSH struct {
		AddToKnownHosts bool
		PrivateKeyPath  *string
	}
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
