package validation

import (
	"fmt"
	"net"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var ErrorExportInterface = fmt.Errorf("validators.extra_UniqueField: Cannot export private field!")
var ErrorFieldNotFound = fmt.Errorf("validators.extra_UniqueField: Field not found!")

// extra_AlphaNumericDash checks whether the field contains only alphanumeric characters
// (a-Z0-9) and hyphen (-).
func extra_AlphaNumericHyphen(fl validator.FieldLevel) bool {
	return regex("^[a-zA-Z0-9-]*$", fl.Field().String())
}

// extra_AlphaNumericDashUnderscore checks whether the field contains only alphanumeric
// characters (a-Z0-9), hyphen (-) and underscore (_).
func extra_AlphaNumericHyphenUnderscore(fl validator.FieldLevel) bool {
	return regex("^[a-zA-Z0-9-_]*$", fl.Field().String())
}

// extra_VSemVer checks whether the field is a valid semantic version
// prefixed with 'v'.
func extra_VSemVer(fl validator.FieldLevel) bool {
	return regex("^(v){1}(\\*|\\d+(\\.\\d+){2})$", fl.Field().String())
}

// extra_IPInRange checks whether the field is a valid IP within provided CIDR
func extra_IPInRange(fl validator.FieldLevel) bool {
	_, subnet, err := net.ParseCIDR(fl.Param())

	if err != nil {
		return false
	}

	ip := net.ParseIP(fl.Field().String())
	return subnet.Contains(ip)
}

// extra_UniqueField returns true if struct field with a given name is unique for
// all slice elements.
func extra_UniqueField(fl validator.FieldLevel) bool {
	fName := string(fl.Param())
	rv := getDeepValue(fl.Field())

	if rv.Kind() != reflect.Slice {
		return true
	}

	var fields []interface{}

	for i := 0; i < rv.Len(); i++ {
		ri := rv.Index(i)

		if ri.Kind() != reflect.Struct {
			return true
		}

		f := structFieldValue(ri, fName)
		fields = append(fields, f)
	}

	for i := 0; i < len(fields); i++ {
		for j := 0; j < len(fields); j++ {
			if i != j && fields[i] == fields[j] {
				return false
			}
		}
	}

	return true
}

// structFieldValue returns value of a struct field with a given name.
func structFieldValue(rs reflect.Value, fName string) interface{} {
	for j := 0; j < rs.NumField(); j++ {
		rf := rs.Field(j)
		rfName := rs.Type().Field(j).Name

		if rfName != fName {
			continue
		}

		v := getDeepValue(rf)

		if v.Kind() == reflect.Invalid {
			return nil
		}

		if !v.CanInterface() {
			panic(ErrorExportInterface)
		}

		return v.Interface()
	}

	panic(ErrorFieldNotFound)
}
