package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type Action string

const (
	UNKNOWN   Action = ""
	SKIP      Action = "SKIP"
	FAIL      Action = "FAIL"
	OMITEMPTY Action = "OMITEMPTY"
	NOT_EMPTY Action = "NOT_EMPTY"
)

// Validator represents a validation rule.
type Validator struct {
	Tags   string
	Err    string
	ignore bool
	action Action
}

// None is an empty validator that does nothing (is skipped).
// It is useful for custom validators.
var None Validator = Validator{}

var customValidators = make(map[string]Validator)

// RegisterCustomValidator registers custom validator with custom key.
func RegisterCustomValidator(key string, v Validator) {
	customValidators[key] = v
}

// RemoveCustomValidator removes custom validator from a list.
func RemoveCustomValidator(key string) {
	delete(customValidators, key)
}

// RegisterCustomValidator registers custom validator with custom key.
func ClearCustomValidators() {
	customValidators = make(map[string]Validator)
}

// initialize creates new singleton validator if its value is nil.
func initialize() {
	if validate != nil {
		return
	}

	validate = validator.New()
	validate.RegisterTagNameFunc(fieldName)
	validate.RegisterValidation("extra_alphanumhyp", extra_AlphaNumericHyphen)
	validate.RegisterValidation("extra_alphanumhypus", extra_AlphaNumericHyphenUnderscore)
	validate.RegisterValidation("extra_vsemver", extra_VSemVer)
	validate.RegisterValidation("extra_ipinrange", extra_IPInRange)
	validate.RegisterValidation("extra_uniquefield", extra_UniqueField)
}

// validate validates the provided value against the validator.
// It returns encountered validation errors and boolean indicating
// whether to skip further validation of a field or not.
func (v *Validator) validate(value interface{}) (ValidationErrors, bool) {
	initialize()

	if v.ignore {
		return nil, false
	}

	switch v.action {
	case SKIP:
		return nil, true
	case FAIL:
		return v.ToError(), false
	case OMITEMPTY:
		return nil, isEmpty(value)
	case NOT_EMPTY:
		if isEmpty(value) {
			return v.ToError(), false
		} else {
			return nil, false
		}
	}

	errs := validate.Var(value, v.Tags)
	es := toValidationErrors(errs)

	if len(v.Err) > 0 {
		for i := range es {
			es[i].Err = v.Err
		}
	}

	return es, false
}

func (v Validator) ToError() ValidationErrors {
	return ValidationErrors{
		{
			Tag:       v.Tags,
			ActualTag: v.Tags,
			Err:       v.Err,
		},
	}
}

// Error overwrites the validator's default error message with the user-defined error
// and returns the modified validator.
func (v Validator) Error(err string) Validator {
	v.Err = err
	return v
}

// Errorf overwrites the validator's default error message with the formatted user-defined
// error and returns the modified validator.
func (v Validator) Errorf(err string, opt ...interface{}) Validator {
	v.Err = fmt.Sprintf(err, opt...)
	return v
}

// When allows validator to be applied only when the given condition is met.
func (v Validator) When(condition bool) Validator {
	v.ignore = !condition
	return v
}

// Custom returns custom validator registered with the given key.
func Custom(key string) Validator {
	return customValidators[key]
}

// Tags returns a new validator with the given tags. It is a generic validator that
// allows use of any validation rule from 'github.com/go-playground/validator' library.
func Tags(tags string) Validator {
	return Validator{
		Tags: tags,
	}
}

// OmitEmpty prevents further validation of the field, if the field is empty.
func OmitEmpty() Validator {
	return Validator{
		Tags:   "omitempty",
		action: OMITEMPTY,
	}
}

// Skip prevents further validation of the field.
func Skip() Validator {
	return Validator{
		Tags:   "-",
		action: SKIP,
	}
}

// Fail triggers validation error.
func Fail() Validator {
	return Validator{
		Tags:   "fail",
		Err:    "Field '{.Field}' (forcefully) failed validation.",
		action: FAIL,
	}
}

// Required validator verifies that value is provided.
func Required() Validator {
	return Validator{
		Tags: "required",
		Err:  "Field '{.Field}' is required.",
	}
}

// NotEmpty validator checks whether value is not blank or with zero length.
func NotEmpty() Validator {
	return Validator{
		Tags:   "notempty",
		Err:    "Field '{.Field}' is required and cannot be empty.",
		action: NOT_EMPTY,
	}
}

// Unique validator checks whether all elements within array, slice or map are unique.
func Unique() Validator {
	return Validator{
		Tags: "unique",
		Err:  "All elements within '{.Field}' must be unique.",
	}
}

// Unique validator checks whether the given field is unique for all elements within
// a slice of struct.
func UniqueField(field string) Validator {
	return Validator{
		Tags: fmt.Sprintf("extra_uniquefield=%s", field),
		Err:  "Field '{.Param}' must be unique for each element in '{.Field}'.",
	}
}

// Min checks whether the field value is greater than or equal to the specified value.
// In case of strings, slices, arrays and maps the length is checked.
func Min(value int) Validator {
	return Validator{
		Tags: fmt.Sprintf("min=%d", value),
		Err:  "Minimum value for field '{.Field}' is {.Param} (actual: {.Value}).",
	}
}

// Max checks whether the field value is less than or equal to the specified value.
// In case of strings, slices, arrays and maps the length is checked.
func Max(value int) Validator {
	return Validator{
		Tags: fmt.Sprintf("max=%d", value),
		Err:  "Maximum value for field '{.Field}' is {.Param} (actual: {.Value}).",
	}
}

// Len checks if the field length matches the specified value.
func Len(value int) Validator {
	return Validator{
		Tags: fmt.Sprintf("len=%d", value),
		Err:  "Length of the field '{.Field}' must be {.Param} (actual: {.Value}).",
	}
}

// MinLen checks whether the field length is greater than or equal to the specified value.
func MinLen(value int) Validator {
	return Min(value).Error("Minimum length of the field '{.Field}' is {.Param} (actual: {.Value})")
}

// MaxLen checks whether the field length is less than or equal to the specified value.
func MaxLen(value int) Validator {
	return Max(value).Error("Maximum length of the field '{.Field}' is {.Param} (actual: {.Value})")
}

// IP checks whether the field value is a valid IP address.
func IP() Validator {
	return Validator{
		Tags: "ip",
		Err:  "Field '{.Field}' must be a valid IP address (actual: {.Value}).",
	}
}

// IPv4 checks whether the field value is a valid v4 IP address.
func IPv4() Validator {
	return Validator{
		Tags: "ipv4",
		Err:  "Field '{.Field}' must be a valid IPv4 address (actual: {.Value}).",
	}
}

// IPv6 checks whether the field value is a valid v6 IP address.
func IPv6() Validator {
	return Validator{
		Tags: "ipv6",
		Err:  "Field '{.Field}' must be a valid IPv6 address (actual: {.Value}).",
	}
}

// CIDR checks whether the field value is a valid CIDR address.
func CIDR() Validator {
	return Validator{
		Tags: "cidr",
		Err:  "Field '{.Field}' must be a valid CIDR address (actual: {.Value}).",
	}
}

// CIDRv4 checks whether the field value is a valid v4 CIDR address.
func CIDRv4() Validator {
	return Validator{
		Tags: "cidrv4",
		Err:  "Field '{.Field}' must be a valid CIDRv4 address (actual: {.Value}).",
	}
}

// CIDRv6 checks whether the field value is a valid v6 CIDR address.
func CIDRv6() Validator {
	return Validator{
		Tags: "cidrv6",
		Err:  "Field '{.Field}' must be a valid CIDRv6 address (actual: {.Value}).",
	}
}

// IPInRange checks whether the field value is contained within the specified CIDR.
func IPInRange(cidr string) Validator {
	return Validator{
		Tags: fmt.Sprintf("extra_ipinrange=%s", cidr),
		Err:  fmt.Sprintf("Field '{.Field}' must be a valid IP address within '{.Param}' subnet. (actual: {.Value})"),
	}
}

// MAC checks whether the field value is a valid MAC address.
func MAC() Validator {
	return Validator{
		Tags: "mac",
		Err:  "Field '{.Field}' must be a valid MAC address (actual: {.Value}).",
	}
}

// OneOf checks whether the field value equals one of the specified values.
// If no value is provided, the validation always fails.
func OneOf[T any](values ...T) Validator {
	var s []string

	for _, v := range values {
		s = append(s, toString(v))
	}

	oneOf := strings.Join(s, " ")
	valid := strings.Join(s, "|")

	return Validator{
		Tags: fmt.Sprintf("oneof=%s", oneOf),
		Err:  fmt.Sprintf("Field '{.Field}' must be one of the following values: [%s] (actual: {.Value}).", valid),
	}
}

// Alpha checks whether the field contains only ASCII alpha characters.
func Alpha() Validator {
	return Validator{
		Tags: "alpha",
		Err:  "Field '{.Field}' can contain only alpha characters (a-Z). (actual: {.Value})",
	}
}

// Numeric checks whether the field contains only numeric characters.
// Validation fails for integers and floats.
func Numeric() Validator {
	return Validator{
		Tags: "numeric",
		Err:  "Field '{.Field}' can contain only numeric characters (0-9). (actual: {.Value})",
	}
}

// AlphaNumeric checks whether the field contains only alphanumeric characters.
// Validation fails for integers and floats.
func AlphaNumeric() Validator {
	return Validator{
		Tags: "alphanum",
		Err:  "Field '{.Field}' can contain only alphanumeric characters. (actual: {.Value})",
	}
}

// AlphaNumericDash checks whether the field contains only alphanumeric characters
// (a-Z0-9) and hyphen (-). Validation fails non string values.
func AlphaNumericHyp() Validator {
	return Validator{
		Tags: "extra_alphanumhyp",
		Err:  "Field '{.Field}' can contain only alphanumeric characters and hyphen. (actual: {.Value})",
	}
}

// AlphaNumericDash checks whether the field contains only alphanumeric characters
// (a-Z0-9), hyphen (-) and underscore (_). Validation fails non string values.
func AlphaNumericHypUS() Validator {
	return Validator{
		Tags: "extra_alphanumhypus",
		Err:  "Field '{.Field}' can contain only alphanumeric characters, hyphen and underscore. (actual: {.Value})",
	}
}

// Lowercase checks whether the field contains only lowercase characters.
func Lowercase() Validator {
	return Validator{
		Tags: "lowercase",
		Err:  "Field '{.Field}' can contain only lowercase characters. (actual: {.Value})",
	}
}

// Uppercase checks whether the field contains only uppercase characters.
func Uppercase() Validator {
	return Validator{
		Tags: "uppercase",
		Err:  "Field '{.Field}' can contain only uppercase characters. (actual: {.Value})",
	}
}

// File checks whether the field is a valid file path and whether file exists.
func FileExists() Validator {
	return Validator{
		Tags: "file",
		Err:  "Field '{.Field}' must be a valid file path that points to an existing file. (actual: {.Value})",
	}
}

// URL checks whether the field is a valid URL.
func URL() Validator {
	return Validator{
		Tags: "url",
		Err:  "Field '{.Field}' must be a valid URL. (actual: {.Value})",
	}
}

// SemVer checks whether the field is a valid semantic version.
func SemVer() Validator {
	return Validator{
		Tags: "semver",
		Err:  "Field '{.Field}' must be a valid semantic version (e.g. 1.2.3). (actual: {.Value})",
	}
}

// VSemVer checks whether the field is a valid semantic version and is prefixed with 'v'.
func VSemVer() Validator {
	return Validator{
		Tags: "extra_vsemver",
		Err:  "Field '{.Field}' must be a valid semantic version prefixed with 'v' (e.g. v1.2.3). (actual: {.Value})",
	}
}
