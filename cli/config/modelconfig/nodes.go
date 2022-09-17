package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Nodes struct {
	LoadBalancer *LoadBalancer
	Master       *Worker
	Worker       *Worker
}

func (n Nodes) Validate() error {
	return validation.ValidateStruct(&n,
		validation.Field(n.LoadBalancer),
		validation.Field(n.Master),
		validation.Field(n.Worker),
	)
}
