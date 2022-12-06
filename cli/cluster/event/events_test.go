package event

import (
	"cli/config/modelconfig"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventPaths(t *testing.T) {
	var events []Event

	events = append(events, ModifyEvents...)
	events = append(events, ScaleEvents...)
	events = append(events, UpgradeEvents...)

	for _, e := range events {
		for _, p := range e.Paths() {
			validateConfigPath(t, p)
		}
	}
}

func validateConfigPath(t *testing.T, path string) {
	paths := strings.Split(path, ".")
	cType := reflect.TypeOf(modelconfig.Config{})

	pass := typePathExists(cType, paths...)

	if !pass {
		err := fmt.Sprintf("Change event path '%s' does not represent any '%s' object!", path, cType)
		assert.Fail(t, err)
	}
}

func typePathExists(t reflect.Type, path ...string) bool {
	if len(path) == 0 {
		return true
	}

	x := reflect.New(t)

	if x.Kind() == reflect.Pointer {
		x = reflect.Indirect(x)
	}

	switch x.Kind() {
	case reflect.Struct:
		for i := 0; i < x.NumField(); i++ {
			f := x.Type().Field(i)

			if f.Name != path[0] {
				continue
			}

			if len(path) == 1 {
				return true
			}

			if !f.IsExported() {
				return false
			}

			return typePathExists(f.Type, path[1:]...)
		}
	case reflect.Slice, reflect.Array:
		if path[0] == "*" {
			return typePathExists(x.Type().Elem(), path[1:]...)
		}
	}

	return false
}
