package modelconfig

type OperatingSystem string

const (
	ubuntu   OperatingSystem = "ubuntu"
	ubuntu20                 = "ubuntu20"
	ubuntu22                 = "ubuntu22"
	debian                   = "debian"
	debian11                 = "debian11"
	custom                   = "custom"
)

type UserName string
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
