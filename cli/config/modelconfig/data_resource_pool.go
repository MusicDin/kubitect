package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type DataResourcePool struct {
	Name *DataResourcePoolName `yaml:"name"`
	Path *ResourcePath         `yaml:"path"`
}

func (drp DataResourcePool) Validate() error {
	return validation.ValidateStruct(&drp,
		validation.Field(drp.Name),
		validation.Field(drp.Path),
	)
}

type DataResourcePoolName string

func (d DataResourcePoolName) Validate() error {
	return validation.Validate(&d, StringNotEmptyAlphaNumeric...)
}
