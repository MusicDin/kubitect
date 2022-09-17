package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

const (
	MinHostsLength = 1
	MaxHostsLength = 0
)

type Hosts struct {
	list []Host
}

func (h Hosts) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.list, validation.Length(MinHostsLength, MaxHostsLength)),
	)
}
