package validation

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	TAG_NAME      = "validate"
	TAG_OPTION_ID = "id"
)

var (
	tag           = TAG_NAME
	fieldNameTags = []string{"json", "yaml"}
)

var (
	ErrorMissingStructPointer = fmt.Errorf("validation.Struct: Struct pointer is missing!")
	ErrorInvalidStructPointer = fmt.Errorf("validation.Struct: First argument must be a pointer to a struct!")
	ErrorStructFieldNotFound  = fmt.Errorf("validation.Struct: Struct field not found int the struct!")
)

// Structure is validatable if it contains Validate method that returns an
// error.
type Validatable interface {
	Validate() error
}

// FieldValidator contains a pointer to a struct field and corresponding
// validators.
type FieldValidator struct {
	fieldPtr   interface{}
	validators []Validator
}

// Var validates a variable against the provided validators. Validation will
// dive deeper, if variable is validatable struct, map or slice.
func Var(value interface{}, validators ...Validator) error {
	initialize()

	errs := make(ValidationErrors, 0)

	for _, v := range validators {
		err, skip := v.validate(value)

		if err != nil {
			errs.append(err)
			break
		}

		if skip {
			break
		}
	}

	rv := reflect.ValueOf(value)
	dv := getDeepValue(rv)

	if value != nil && dv.Kind() != reflect.Invalid {
		ri := dv.Interface()

		switch dv.Kind() {
		case reflect.Slice, reflect.Array:
			errs.append(validateSlice(dv))
		case reflect.Map:
			errs.append(validateMap(dv))
		default:
			if dv.Kind() == reflect.Struct || rv.Kind() == reflect.Pointer {
				if val, ok := ri.(Validatable); ok {
					errs.append(val.Validate())
				}
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// validateSlice triggers validation of each validatable child element.
func validateSlice(rv reflect.Value) ValidationErrors {
	errs := make(ValidationErrors, 0)

	for i := 0; i < rv.Len(); i++ {
		vi := rv.Index(i).Interface()
		id := tagOptionId(tag, rv.Index(i))

		if val, ok := vi.(Validatable); ok {
			if id == nil {
				id = i
			}

			errs.subAppend(val.Validate(), i, id)
		}
	}

	return errs
}

// validateMap triggers validation of each validatable child element.
func validateMap(rv reflect.Value) ValidationErrors {
	errs := make(ValidationErrors, 0)

	for _, key := range rv.MapKeys() {
		vi := rv.MapIndex(key).Interface()

		if val, ok := vi.(Validatable); ok {
			errs.subAppend(val.Validate(), key, key)
		}
	}

	return errs
}

// Struct validates every given field against corresponding validators.
//
// It panics if a provided value is not a pointer to a struct or if a field
// cannot be found within the structure.
func Struct(sPtr interface{}, validators ...FieldValidator) error {

	if sPtr == nil {
		panic(ErrorMissingStructPointer)
	}

	rv := reflect.ValueOf(sPtr)
	rs := getDeepValue(rv)

	if rv.Kind() != reflect.Pointer || rs.Kind() != reflect.Struct {
		panic(ErrorInvalidStructPointer)
	}

	errs := make(ValidationErrors, 0)

	for _, v := range validators {
		rf := reflect.ValueOf(v.fieldPtr)
		sf := findStructField(rs, rf)

		if sf == nil {
			panic(ErrorStructFieldNotFound)
			//continue
		}

		err := Var(rf.Interface(), v.validators...)
		errs.subAppend(err, sf.Name, fieldName(*sf))
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// findStructField looks for a field in the given struct.
// If found, the field info is returned. Otherwise, nil is returned.
func findStructField(s reflect.Value, f reflect.Value) *reflect.StructField {
	for i := 0; i < s.NumField(); i++ {
		ft := s.Type().Field(i)

		if !s.Field(i).CanAddr() {
			continue
		}

		if f.Pointer() == s.Field(i).Addr().Pointer() {
			if ft.Type == f.Elem().Type() {
				return &ft
			}
		}
	}

	return nil
}

// Field returns new FieldValidator.
func Field(fPtr interface{}, validators ...Validator) FieldValidator {
	if fPtr == nil || reflect.ValueOf(fPtr).Kind() != reflect.Pointer {
		panic(ErrorFieldIsNotPointer)
	}

	return FieldValidator{
		fieldPtr:   fPtr,
		validators: validators,
	}
}

// fieldName extracts field name from labels provided in fieldNameTags.
// If name is not specified in any of the given tags, struct field name
// is returned.
func fieldName(fld reflect.StructField) string {
	name := fld.Name

	tags := append([]string{tag}, fieldNameTags...)

	for _, tag := range tags {
		tName := strings.SplitN(fld.Tag.Get(tag), ",", 2)[0]

		if len(tName) > 0 {
			name = tName
			break
		}
	}

	if name == "-" {
		return ""
	}

	return name
}

// SetTagName sets primary name tag from which validation options are
// extracted.
func SetTagName(tagName string) {
	tag = tagName
}

// tagOptionId returns value of the field with id label.
func tagOptionId(tagName string, v reflect.Value) interface{} {
	if v.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < v.NumField(); i++ {
		if hasTagOption(tagName, v.Type().Field(i), TAG_OPTION_ID) {
			return v.Field(i).Interface()
		}
	}

	return nil
}

// hasTagOption checks whether specified tag contains the given option.
func hasTagOption(tagName string, field reflect.StructField, option string) bool {
	tag := field.Tag.Get(tagName)
	options := strings.Split(tag, ",")

	if len(options) < 2 {
		return false
	}

	for _, o := range options[1:] {
		o = strings.TrimSpace(o)
		o = strings.ToLower(o)

		if o == option {
			return true
		}
	}

	return false
}
