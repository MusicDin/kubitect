package modelconfig

type NodeTemplate struct {
	UpdateOnBoot *bool
	User         *UserName

	DNS []string
	OS  struct {
		Distro *OperatingSystem
	}

	SSH struct {
		AddToKnownHosts bool
		PrivateKeyPath  *string
	}
}
