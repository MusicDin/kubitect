package provisioner

type Provisioner interface {
	Plan() (bool, error)
	Apply() error
	Destroy() error
}
