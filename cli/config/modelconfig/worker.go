package modelconfig

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Worker struct {
	Default   *WorkerDefault
	Instances []Instance
}

func (w Worker) Validate() error {
	return validation.ValidateStruct(&w,
		validation.Field(&w.Instances),
		validation.Field(w.Default),
	)
}
