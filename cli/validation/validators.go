package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

var ErrorFieldIsNotPointer = fmt.Errorf("validators.Field: First argument must be a pointer to a struct field!")

type Action string

const (
	UNKNOWN   = ""
	OMITEMPTY = "OMITEMPTY"
	SKIP      = "SKIP"
)

// Validator represents a validation rule.
type Validator struct {
	Tags   string
	Err    string
	ignore bool
	action Action
}

// initialize creates new singleton validator if its value is nil.
func initialize() {
	if validate != nil {
		return
	}

	validate = validator.New()
	validate.RegisterTagNameFunc(fieldName)
	validate.RegisterValidation("custom_alphanumhyp", custom_AlphaNumericHyphen)
	validate.RegisterValidation("custom_alphanumhypus", custom_AlphaNumericHyphenUnderscore)
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
	case OMITEMPTY:
		return nil, isEmpty(value)
	case SKIP:
		return nil, true
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

// Error overrides the default error of the validator and returns modified validator.
func (v Validator) Error(err string) Validator {
	v.Err = err
	return v
}

// When allows validator to be applied only when the given condition is met.
func (v Validator) When(condition bool) Validator {
	v.ignore = !condition
	return v
}

// custom_AlphaNumericDash checks whether the field contains only alphanumeric characters
// (a-Z0-9) and hyphen (-).
func custom_AlphaNumericHyphen(fl validator.FieldLevel) bool {
	return regex("^[a-zA-Z0-9-]*$", fl.Field().String())
}

// custom_AlphaNumericDashUnderscore checks whether the field contains only alphanumeric
// characters (a-Z0-9), hyphen (-) and underscore (_).
func custom_AlphaNumericHyphenUnderscore(fl validator.FieldLevel) bool {
	return regex("^[a-zA-Z0-9-_]*$", fl.Field().String())
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

// Required validator verifies that value is provided.
func Required() Validator {
	return Validator{
		Tags: "required",
		Err:  "Property '{.Namespace}' is required.",
	}
}

// Min checks whether the field value is greater than or equal to the specified value.
// In case of strings, slices, arrays and maps the length is checked.
func Min(value int) Validator {
	return Validator{
		Tags: fmt.Sprintf("min=%d", value),
		Err:  "Minimum value for property '{.Namespace}' is {.Param} (actual: {.Value}).",
	}
}

// Max checks whether the field value is less than or equal to the specified value.
// In case of strings, slices, arrays and maps the length is checked.
func Max(value int) Validator {
	return Validator{
		Tags: fmt.Sprintf("max=%d", value),
		Err:  "Maximum value for property '{.Namespace}' is {.Param} (actual: {.Value}).",
	}
}

// Len checks if the field length matches the specified value.
func Len(value int) Validator {
	return Validator{
		Tags: fmt.Sprintf("len=%d", value),
		Err:  "Length of '{.Namespace}' must be {.Param} (actual: {.Value}).",
	}
}

// MinLen checks whether the field length is greater than or equal to the specified value.
func MinLen(value int) Validator {
	return Min(value).Error("Minimum length of '{.StructField}' is {.Param} (actual: {.Value})")
}

// MaxLen checks whether the field length is less than or equal to the specified value.
func MaxLen(value int) Validator {
	return Max(value).Error("Maximum length of '{.StructField}' is {.Param} (actual: {.Value})")
}

// IP checks whether the field value is a valid IP address.
func IP() Validator {
	return Validator{
		Tags: "ip",
		Err:  "Property '{.Field}' must be a valid IP address (actual: {.Value}).",
	}
}

// IPv4 checks whether the field value is a valid v4 IP address.
func IPv4() Validator {
	return Validator{
		Tags: "ipv4",
		Err:  "Property '{.Field}' must be a valid IPv4 address (actual: {.Value}).",
	}
}

// IPv6 checks whether the field value is a valid v6 IP address.
func IPv6() Validator {
	return Validator{
		Tags: "ipv6",
		Err:  "Property '{.Field}' must be a valid IPv6 address (actual: {.Value}).",
	}
}

// MAC checks whether the field value is a valid MAC address.
func MAC() Validator {
	return Validator{
		Tags: "mac",
		Err:  "Property '{.Field}' must be a valid MAC address (actual: {.Value}).",
	}
}

// OneOf checks whether the field value equals one of the specified values.
// If no value is provided, the validation always fails.
func OneOf(values ...interface{}) Validator {
	var s []string

	for _, v := range values {
		s = append(s, toString(v))
	}

	oneOf := strings.Join(s, " ")
	valid := strings.Join(s, "|")

	return Validator{
		Tags: fmt.Sprintf("oneof=%s", oneOf),
		Err:  fmt.Sprintf("Property '{.Field}' must be one of the following values: [%s] (actual: {.Value}).", valid),
	}
}

// Alpha checks whether the field contains only ASCII alpha characters.
func Alpha() Validator {
	return Validator{
		Tags: "alpha",
		Err:  "Property '{.Field}' can contain only alpha characters (a-Z). (actual: {.Value})",
	}
}

// Numeric checks whether the field contains only numeric characters.
// Validation fails for integers and floats.
func Numeric() Validator {
	return Validator{
		Tags: "numeric",
		Err:  "Property '{.Field}' can contain only numeric characters (0-9). (actual: {.Value})",
	}
}

// AlphaNumeric checks whether the field contains only alphanumeric characters.
// Validation fails for integers and floats.
func AlphaNumeric() Validator {
	return Validator{
		Tags: "alphanum",
		Err:  "Property '{.Field}' can contain only alphanumeric characters. (actual: {.Value})",
	}
}

// AlphaNumericDash checks whether the field contains only alphanumeric characters
// (a-Z0-9) and hyphen (-). Validation fails non string values.
func AlphaNumericHyp() Validator {
	return Validator{
		Tags: "custom_alphanumhyp",
		Err:  "Property '{.Field}' can contain only alphanumeric characters and hyphen. (actual: {.Value})",
	}
}

// AlphaNumericDash checks whether the field contains only alphanumeric characters
// (a-Z0-9), hyphen (-) and underscore (_). Validation fails non string values.
func AlphaNumericHypUS() Validator {
	return Validator{
		Tags: "custom_alphanumhypus",
		Err:  "Property '{.Field}' can contain only alphanumeric characters, hyphen and underscore. (actual: {.Value})",
	}
}

// Lowercase checks whether the field contains only lowercase characters.
func Lowercase() Validator {
	return Validator{
		Tags: "lowercase",
		Err:  "Property '{.Field}' can contain only lowercase characters. (actual: {.Value})",
	}
}

// Uppercase checks whether the field contains only uppercase characters.
func Uppercase() Validator {
	return Validator{
		Tags: "uppercase",
		Err:  "Property '{.Field}' can contain only uppercase characters. (actual: {.Value})",
	}
}
