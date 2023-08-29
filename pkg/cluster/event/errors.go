package event

import "fmt"

// ValidationError represents an error that occurs during the validation
// of the rule related object. It contains a descriptive message about
// the nature of the validation error and the object that caused the error.
type ValidationError struct {
	message string
}

func NewValidationError(o any, errFmt string, errArgs ...any) ValidationError {
	errMsg := fmt.Sprintf(errFmt, errArgs...)

	switch v := o.(type) {
	case RulePathSegment:
		return ValidationError{fmt.Sprintf("Rule path segment %q: %s", v.path, errMsg)}
	case RulePath:
		return ValidationError{fmt.Sprintf("Rule path %q: %s", v.path, errMsg)}
	case Rule:
		return ValidationError{fmt.Sprintf("Rule %q: %s", v.MatchPath.path, errMsg)}
	default:
		return ValidationError{errMsg}
	}
}

func (e ValidationError) Error() string {
	return e.message
}
