package modelconfig

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"os"
)

const (
	MinStringLength = 1
)

var (
	ErrPathDoesNotExist = validation.NewError("path_does_not_exist", "path does not exist")
)

var StringNotEmptyAlphaNumeric = []validation.Rule{
	validation.Length(MinStringLength, 0),
	is.Alphanumeric,
}

func PathExists(value interface{}) error {
	path, _ := value.(string)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err != nil {
			// Internal error.
			return err
		}
		// Path exists.
		return nil
	}
	return ErrPathDoesNotExist
}
