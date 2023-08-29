package event

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/MusicDin/kubitect/pkg/models/config"
	"github.com/stretchr/testify/assert"
)

func TestRulePath_Validation(t *testing.T) {
	tests := []struct {
		RulePath string
		IsValid  bool
	}{
		// Valid rule paths.
		{IsValid: true, RulePath: "a"},
		{IsValid: true, RulePath: " a "},
		{IsValid: true, RulePath: "a.b.c"},
		{IsValid: true, RulePath: "{a}"},
		{IsValid: true, RulePath: "{a,b}"},
		{IsValid: true, RulePath: "*"},
		{IsValid: true, RulePath: "*.*"},
		{IsValid: true, RulePath: "@"},
		{IsValid: true, RulePath: "@*"},
		{IsValid: true, RulePath: "@a"},
		{IsValid: true, RulePath: "@ a"},
		{IsValid: true, RulePath: "@{a}"},
		{IsValid: true, RulePath: "@{a,b}"},
		{IsValid: true, RulePath: "a.*.b.@"},
		{IsValid: true, RulePath: "a.@.b.*"},
		{IsValid: true, RulePath: "a.@{b}"},
		{IsValid: true, RulePath: "a!"},
		{IsValid: true, RulePath: "{a,b}!"},
		{IsValid: true, RulePath: "a.{a,b}!"},
		{IsValid: true, RulePath: " @ { a , b } . * ! "},
		// Invalid rule paths.
		{IsValid: false, RulePath: ""},
		{IsValid: false, RulePath: ","},
		{IsValid: false, RulePath: "."},
		{IsValid: false, RulePath: "!"},
		{IsValid: false, RulePath: "{"},
		{IsValid: false, RulePath: "}"},
		{IsValid: false, RulePath: "}{"},
		{IsValid: false, RulePath: "{}"},
		{IsValid: false, RulePath: "{ }"},
		{IsValid: false, RulePath: "{,}"},
		{IsValid: false, RulePath: "{!}"},
		{IsValid: false, RulePath: "{@}"},
		{IsValid: false, RulePath: "{*}"},
		{IsValid: false, RulePath: "{{}"},
		{IsValid: false, RulePath: "{}}"},
		{IsValid: false, RulePath: "{{}}"},
		{IsValid: false, RulePath: "{a.b}"},
		{IsValid: false, RulePath: "@.@"},
		{IsValid: false, RulePath: "@@"},
		{IsValid: false, RulePath: "**"},
		{IsValid: false, RulePath: "*@"},
		{IsValid: false, RulePath: "*a"},
		{IsValid: false, RulePath: "a."},
		{IsValid: false, RulePath: "*a"},
		{IsValid: false, RulePath: "*{a}"},
		{IsValid: false, RulePath: "@*a"},
		{IsValid: false, RulePath: "@*{a}"},
		{IsValid: false, RulePath: "a!!"},
		{IsValid: false, RulePath: "a!.b"},
	}

	for _, test := range tests {
		err := NewRulePath(test.RulePath).Validate()

		if err == nil && !test.IsValid {
			assert.Fail(t, fmt.Sprintf("Rule path %q should NOT be valid!", test.RulePath))
		}

		if err != nil && test.IsValid {
			assert.Fail(t, fmt.Sprintf("Rule path %q should be valid!", test.RulePath), err)
		}
	}
}

func TestRulePath_WildcardCount(t *testing.T) {
	tests := []struct {
		RulePath string
		Count    int
	}{
		{Count: 0, RulePath: ""},
		{Count: 0, RulePath: "a.b.c"},
		{Count: 0, RulePath: "a.@{b}"},
		{Count: 1, RulePath: "@"},
		{Count: 1, RulePath: "*"},
		{Count: 1, RulePath: "@*"},
		{Count: 2, RulePath: "*.*"},
		{Count: 2, RulePath: "a.*.b.@"},
		{Count: 2, RulePath: "a.@.b.*"},
		{Count: 2, RulePath: "@.@*"},
		// Invalid wildcards.
		{Count: 0, RulePath: "@@"},
		{Count: 0, RulePath: "*@"},
		{Count: 0, RulePath: "**"},
		{Count: 0, RulePath: "@*a"},
		{Count: 0, RulePath: "@*{a}"},
	}

	for _, test := range tests {
		count := NewRulePath(test.RulePath).WildcardCount()
		if count != test.Count {
			assert.Fail(t, fmt.Sprintf("Invalid wildcard count for path %q.\nExpect: %d\nResult: %d", test.RulePath, test.Count, count))
		}
	}
}

func TestRulePath_IsAnchor(t *testing.T) {
	tests := []struct {
		RulePath string
		IsAnchor bool
	}{
		{IsAnchor: true, RulePath: "@"},
		{IsAnchor: true, RulePath: "a.@"},
		{IsAnchor: true, RulePath: "@{a,b}"},
		{IsAnchor: true, RulePath: "a.@{b}"},
		{IsAnchor: true, RulePath: "a.@{a,b}"},
		{IsAnchor: true, RulePath: "a.@{a,b}.@.*.@"},
		{IsAnchor: true, RulePath: "a.@*"},
		{IsAnchor: true, RulePath: "@@"},
		{IsAnchor: false, RulePath: ""},
		{IsAnchor: false, RulePath: "a.b.c"},
		{IsAnchor: false, RulePath: "*"},
		{IsAnchor: false, RulePath: "*@"},
		{IsAnchor: false, RulePath: "**"},
	}

	for _, test := range tests {
		isAnchor := NewRulePath(test.RulePath).IsAnchorPath()
		if isAnchor != test.IsAnchor {
			assert.Fail(t, fmt.Sprintf("Is path %q anchored?\nExpect: %v\nResult: %v", test.RulePath, test.IsAnchor, isAnchor))
		}
	}
}

func TestRulePath_FindAnchorPath(t *testing.T) {
	tests := []struct {
		RulePath   string
		ChangePath string
		Expect     string
	}{
		{RulePath: "", ChangePath: "@", Expect: "@"},
		{RulePath: "", ChangePath: "a", Expect: "a"},
		{RulePath: "", ChangePath: "a.b", Expect: "a.b"},
		{RulePath: "@", ChangePath: "", Expect: ""},
		{RulePath: "@", ChangePath: "a", Expect: "a"},
		{RulePath: "@", ChangePath: "a.b", Expect: "a"},
		{RulePath: "a.@", ChangePath: "a", Expect: "a"},
		{RulePath: "a.@", ChangePath: "a.b", Expect: "a.b"},
		{RulePath: "a.@", ChangePath: "a.b.c", Expect: "a.b"},
		{RulePath: "a.@.c", ChangePath: "a.b.c", Expect: "a.b"},
		{RulePath: "*.@.*", ChangePath: "a.b.c", Expect: "a.b"},
		{RulePath: "*.@b.*", ChangePath: "a.b.c", Expect: "a.b"},
		{RulePath: "*.@{b}.*", ChangePath: "a.b.c", Expect: "a.b"},
		{RulePath: "*.@{a}.*", ChangePath: "a.b.c", Expect: "a.b.c"},
		{RulePath: "*.@{a,b}.*", ChangePath: "a.b.c", Expect: "a.b"},
	}

	for _, test := range tests {
		err := fmt.Sprintf("Rule path %q should match change path %q.", test.RulePath, test.ChangePath)
		assert.Equal(t, test.Expect, NewRulePath(test.RulePath).FindAnchorPath(test.ChangePath), err)
	}
}

func TestRulePath_Matches(t *testing.T) {
	tests := []struct {
		RulePath   string
		ChangePath string
		IsMatch    bool
	}{
		// Classic matching.
		{IsMatch: true, RulePath: "", ChangePath: ""},
		{IsMatch: true, RulePath: "", ChangePath: "."},
		{IsMatch: false, RulePath: "", ChangePath: "a"},
		{IsMatch: true, RulePath: "a", ChangePath: "a"},
		{IsMatch: true, RulePath: "a", ChangePath: "a.b"},
		{IsMatch: true, RulePath: "a.b", ChangePath: "a.b"},
		{IsMatch: true, RulePath: "a.b", ChangePath: "a.b.c"},
		{IsMatch: false, RulePath: "a", ChangePath: ""},
		{IsMatch: false, RulePath: "a", ChangePath: "."},
		{IsMatch: false, RulePath: "a.b", ChangePath: "a"},
		// Path options.
		{IsMatch: true, RulePath: "{}", ChangePath: ""},
		{IsMatch: false, RulePath: "{}", ChangePath: "a"},
		{IsMatch: true, RulePath: "{*}", ChangePath: "*"},
		{IsMatch: true, RulePath: "{@}", ChangePath: "@"},
		{IsMatch: false, RulePath: "{*}", ChangePath: "x"},
		{IsMatch: false, RulePath: "{@}", ChangePath: "x"},
		{IsMatch: true, RulePath: "{a,b}", ChangePath: "a"},
		{IsMatch: true, RulePath: "{a,b}", ChangePath: "b"},
		{IsMatch: false, RulePath: "{a,b}", ChangePath: ""},
		{IsMatch: true, RulePath: "{a,   }", ChangePath: ""},
		{IsMatch: true, RulePath: "{a,b}", ChangePath: "a"},
		{IsMatch: true, RulePath: "{a, b}", ChangePath: "a"},
		{IsMatch: true, RulePath: "{a, b }", ChangePath: "a"},
		{IsMatch: true, RulePath: "{ a , b }", ChangePath: "a"},
		// Wildcard.
		{IsMatch: true, RulePath: "*", ChangePath: ""},
		{IsMatch: true, RulePath: "*", ChangePath: "."},
		{IsMatch: true, RulePath: "*", ChangePath: "a"},
		{IsMatch: true, RulePath: "*", ChangePath: "a.b"},
		{IsMatch: false, RulePath: "*.b", ChangePath: "a"},
		{IsMatch: true, RulePath: "*.b", ChangePath: "a.b"},
		{IsMatch: true, RulePath: "*.b", ChangePath: "a.b.c"},
		{IsMatch: false, RulePath: "a.*", ChangePath: "a"},
		{IsMatch: true, RulePath: "a.*", ChangePath: "a."},
		{IsMatch: true, RulePath: "a.*", ChangePath: "a.b"},
		{IsMatch: true, RulePath: "a.*", ChangePath: "a.b.c"},
		{IsMatch: false, RulePath: "a.*{x}", ChangePath: "a.b.c"},
		{IsMatch: false, RulePath: "a.*{x}", ChangePath: "a.x.c"},
		// Anchor.
		{IsMatch: true, RulePath: "@", ChangePath: ""},
		{IsMatch: true, RulePath: "@", ChangePath: "."},
		{IsMatch: true, RulePath: "@", ChangePath: "a"},
		{IsMatch: true, RulePath: "@", ChangePath: "a.b"},
		{IsMatch: false, RulePath: "@.b", ChangePath: "a"},
		{IsMatch: true, RulePath: "@.b", ChangePath: "a.b"},
		{IsMatch: true, RulePath: "@.b", ChangePath: "a.b.c"},
		{IsMatch: false, RulePath: "a.@", ChangePath: "a"},
		{IsMatch: true, RulePath: "a.@", ChangePath: "a."},
		{IsMatch: true, RulePath: "a.@", ChangePath: "a.b"},
		{IsMatch: true, RulePath: "a.@", ChangePath: "a.b.c"},
		{IsMatch: true, RulePath: "a.@{a,b}", ChangePath: "a.a.c"},
		{IsMatch: true, RulePath: "a.@{a,b}", ChangePath: "a.b.c"},
		{IsMatch: false, RulePath: "a.@{x}", ChangePath: "a.b.c"},
		// Exact.
		{IsMatch: true, RulePath: "!", ChangePath: ""},
		{IsMatch: true, RulePath: "a!", ChangePath: "a"},
		{IsMatch: true, RulePath: "a.b!", ChangePath: "a.b"},
		{IsMatch: true, RulePath: "@a.b!", ChangePath: "a.b"},
		{IsMatch: true, RulePath: "@{a,b}.b!", ChangePath: "a.b"},
		{IsMatch: false, RulePath: "!", ChangePath: "a"},
		{IsMatch: false, RulePath: "a!", ChangePath: "a.b"},
		{IsMatch: false, RulePath: "a!.b", ChangePath: "a"},
		{IsMatch: false, RulePath: "a!.b", ChangePath: "a.b"},
		// Invalid path matches that should be caught by validation.
		// Just so that we are aware of them.
		{IsMatch: true, RulePath: "a!.b", ChangePath: "a!.b"},
		{IsMatch: true, RulePath: "a.*{x}", ChangePath: "a.*{x}"},
		{IsMatch: true, RulePath: "a.@*x", ChangePath: "a.@*x"},
	}

	for _, test := range tests {
		isMatch := NewRulePath(test.RulePath).Matches(test.ChangePath)

		if isMatch && !test.IsMatch {
			assert.Fail(t, fmt.Sprintf("Rule path %q and change path %q should NOT match!", test.RulePath, test.ChangePath))
		}

		if !isMatch && test.IsMatch {
			assert.Fail(t, fmt.Sprintf("Rule path %q and change path %q should match!", test.RulePath, test.ChangePath))
		}
	}
}

// TestRuleLists verifies that paths of the each rule are present in
// the Config structure.
func TestRuleLists(t *testing.T) {
	var rules []Rule

	rules = append(rules, ModifyRules...)
	rules = append(rules, ScaleRules...)
	rules = append(rules, UpgradeRules...)

	for _, r := range rules {
		err := r.Validate()
		if err != nil {
			assert.Fail(t, err.Error())
			continue
		}

		cfgType := reflect.TypeOf(config.Config{})
		nameTags := []string{"yaml"}

		if !isValidRulePath(cfgType, r.MatchPath.segments, nameTags) {
			assert.Fail(t, fmt.Sprintf("Rule path '%s' does not represent any '%s' object!", r.MatchPath.Path(), cfgType))
		}
	}
}

// isValidRulePath returns true if the rule's path can be fully resolved for
// the type 't', otherwise, it returns false. Path validity is
// determined by recursively traversing into the structure fields based on
// the rule's path segments.
func isValidRulePath(t reflect.Type, segments []RulePathSegment, nameTags []string) bool {
	if len(segments) == 0 {
		return true
	}

	if len(segments) == 1 && segments[0].IsWildcard() {
		return true
	}

	x := reflect.New(t)

	if x.Kind() == reflect.Pointer {
		x = reflect.Indirect(x)
	}

	if x.Kind() == reflect.Struct {
		for i := 0; i < x.NumField(); i++ {
			field := x.Type().Field(i)
			fieldName := fieldNameFromTag(field, nameTags)

			if !segments[0].Matches(fieldName) {
				continue
			}

			if len(segments) == 1 {
				return true
			}

			if !field.IsExported() {
				return false
			}

			return isValidRulePath(field.Type, segments[1:], nameTags)
		}
	}

	if x.Kind() == reflect.Array || x.Kind() == reflect.Slice {
		if segments[0].IsWildcard() {
			return isValidRulePath(x.Type().Elem(), segments[1:], nameTags)
		}
	}

	return false
}

// fieldNameFromTag returns the field's name extracted from any given tag (tag's
// content before first comma). Tags are checked in the provided order. If slice
// of tags is empty, name of the field's structure is returned.
func fieldNameFromTag(field reflect.StructField, tags []string) string {
	if len(tags) == 0 {
		return field.Name
	}

	for _, tag := range tags {
		tagName := strings.SplitN(field.Tag.Get(tag), ",", 2)[0]
		if len(tagName) > 0 {
			return tagName
		}
	}

	return ""
}
