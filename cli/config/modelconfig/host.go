package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Host struct {
	Connection           *Connection
	DataResourcePools    []DataResourcePool
	Default              bool
	Name                 *HostName
	MainResourcePoolPath *ResourcePath
}

func (h Host) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(h.Name, validation.Required),
		validation.Field(h.Connection),
		validation.Field(&h.DataResourcePools),
		validation.Field(h.MainResourcePoolPath),
	)
}
