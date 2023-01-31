package provisioner

type Provisioner interface {
	Init() error
	Plan() (bool, error)
	Apply() error
	Destroy() error
}
