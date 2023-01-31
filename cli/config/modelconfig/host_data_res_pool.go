package modelconfig

import v "cli/utils/validation"

type DataResourcePool struct {
	Name *string `yaml:"name" opt:",id"`
	Path *string `yaml:"path"`
}

func (rp DataResourcePool) Validate() error {
	return v.Struct(&rp,
		v.Field(&rp.Name, v.Required(), v.AlphaNumericHyp()),
		v.Field(&rp.Path), // TODO: Valid file path. File does not need to exist.
	)
}
