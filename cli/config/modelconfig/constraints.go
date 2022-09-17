package modelconfig

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	MinStringLength = 1
)

var StringNotEmptyAlphaNumeric = []validation.Rule{
	validation.Length(MinStringLength, 0),
	is.Alphanumeric,
}

var (
	ErrPathDoesNotExist = validation.NewError("path_does_not_exist", "path does not exist")
)
