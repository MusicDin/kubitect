package modelconfig

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"os"
	"regexp"
)

const (
	MinStringLength = 1
)

const (
	AlphaNumericMinusRegex = "^[a-zA-Z0-9-_]+$"
)

var (
	ErrAlphaNumericMinusRegex = validation.NewError("alpha_numeric_minus_string", "string should contain only alpha numeric and minus character")
	ErrPathDoesNotExist       = validation.NewError("path_does_not_exist", "path does not exist")
)

var StringNotEmptyAlphaNumericMinus = []validation.Rule{
	validation.Match(regexp.MustCompile(AlphaNumericMinusRegex)).ErrorObject(ErrAlphaNumericMinusRegex),
	validation.Length(MinStringLength, 0),
}

func PathExists(value interface{}) error {
	path, _ := value.(string)
	if !GlobalSettings.CheckValidPath {
		return nil
	}
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
