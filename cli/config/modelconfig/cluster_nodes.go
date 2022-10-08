package modelconfig

import v "cli/validation"

type Nodes struct {
	LoadBalancer *LB     `yaml:"loadBalancer"`
	Master       *Worker `yaml:"master"`
	Worker       *Worker `yaml:"worker"`
}

func (n Nodes) Validate() error {
	return v.Struct(&n,
		v.Field(&n.LoadBalancer),
		v.Field(&n.Master),
		v.Field(&n.Worker),
	)
}
