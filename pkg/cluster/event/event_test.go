package event

import (
	"testing"

	"github.com/MusicDin/kubitect/pkg/utils/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mustCompare compares two values and returns the result of a comparison.
// In case of an error, the test is stopped.
func mustCompare(t *testing.T, a any, b any, opts ...cmp.Options) cmp.Result {
	res, err := cmp.Compare(a, b, opts...)
	require.NoError(t, err)
	return *res
}

// mustGenEvents generated events from the comparison result tree.
// In case of an error, the test is stopped.
func mustGenEvents(t *testing.T, a any, b any, rules []Rule) Events {
	events, err := GenerateEvents(mustCompare(t, a, b).Tree(), rules)
	require.NoError(t, err)
	return events
}

func TestEvent(t *testing.T) {
	v1 := map[string]string{"A": "Yes"}
	v2 := map[string]string{"A": "No"}

	r := Rule{MatchPath: NewRulePath("A")}

	events := mustGenEvents(t, v1, v2, []Rule{r})
	require.Len(t, events, 1)
	assert.Equal(t, r, events[0].Rule)
}

func TestEvent_EmptyComparison(t *testing.T) {
	events := mustGenEvents(t, nil, nil, nil)
	require.Len(t, events, 0)
}

func TestEvent_NoDiff(t *testing.T) {
	r := Rule{MatchPath: NewRulePath("*")}

	events := mustGenEvents(t, nil, nil, []Rule{r})
	require.Len(t, events, 0)
}

// Test expects a rule with a longer path to trigger an event.
func TestEvent_RulePrecedenceByPathLength(t *testing.T) {
	type Map map[string]any

	v1 := map[string]Map{"a": {"b": "Yes"}}
	v2 := map[string]Map{"a": {"b": "No"}}

	r1 := Rule{Type: Error, MatchPath: NewRulePath("a")}
	r2 := Rule{Type: Allow, MatchPath: NewRulePath("a.b")}

	events := mustGenEvents(t, v1, v2, []Rule{r1, r2})
	require.Len(t, events, 1)
	assert.Equal(t, r2, events[0].Rule)
}

// Test expects a rule with the longest matching path to trigger an event.
func TestEvent_RulePrecedenceByPathLength_MultiPathRules(t *testing.T) {
	type Map map[string]any

	v1 := map[string]Map{"a": {"b": "Yes"}}
	v2 := map[string]Map{"a": {"b": "No"}}

	r1 := Rule{Type: Error, MatchPath: NewRulePath("a")}
	r2 := Rule{Type: Allow, MatchPath: NewRulePath("a.b")}
	r3 := Rule{Type: Allow, MatchPath: NewRulePath("@")}
	r4 := Rule{Type: Error, MatchPath: NewRulePath("*")}
	r5 := Rule{Type: Allow, MatchPath: NewRulePath("invalid")}

	events := mustGenEvents(t, v1, v2, []Rule{r1, r2, r3, r4, r5})
	require.Len(t, events, 1)
	assert.Equal(t, r2, events[0].Rule)
	assert.Equal(t, "a.b", events[0].Change.Path)
}

// Test expects a rule with a more specific path (path with less asterisks)
// to trigger an event.
func TestEvent_RulePrecedenceByPathAsterisksCount(t *testing.T) {
	type Map map[string]any

	v1 := map[string]Map{"a": {"b": "Yes"}}
	v2 := map[string]Map{"a": {"b": "No"}}

	r1 := Rule{Type: Error, MatchPath: NewRulePath("*.*")}
	r2 := Rule{Type: Warn, MatchPath: NewRulePath("a.*")}
	r3 := Rule{Type: Allow, MatchPath: NewRulePath("a.b")}

	events := mustGenEvents(t, v1, v2, []Rule{r1, r2, r3})
	require.Len(t, events, 1)
	assert.Equal(t, r3, events[0].Rule)
}

// Test expects a rule with with higher priority to take precedence over other
// rules with the same paths.
func TestEvent_RulePrecedenceByPriority(t *testing.T) {
	type Map map[string]any

	v1 := map[string]Map{"a": {"b": "Yes"}}
	v2 := map[string]Map{"a": {"b": "No"}}

	r1 := Rule{Type: Allow, MatchPath: NewRulePath("a.b")}
	r2 := Rule{Type: Warn, MatchPath: NewRulePath("a.b")}
	r3 := Rule{Type: Error, MatchPath: NewRulePath("a.b")}

	events := mustGenEvents(t, v1, v2, []Rule{r1, r2, r3})
	require.Len(t, events, 1)
	assert.Equal(t, r3, events[0].Rule)
}

// Test expects a rule with with higher custom priority to take precedence over
// other rules with the same paths.
func TestEvent_RulePrecedenceByPriority_Custom(t *testing.T) {
	type Map map[string]any

	v1 := map[string]Map{"a": {"b": "Yes"}}
	v2 := map[string]Map{"a": {"b": "No"}}

	r1 := Rule{Type: 10, MatchPath: NewRulePath("a.b")}
	r2 := Rule{Type: 50, MatchPath: NewRulePath("a.b")}
	r3 := Rule{Type: 30, MatchPath: NewRulePath("a.b")}

	events := mustGenEvents(t, v1, v2, []Rule{r1, r2, r3})
	require.Len(t, events, 1)
	assert.Equal(t, r2, events[0].Rule)
}

func TestEvent_RulePathWildcard(t *testing.T) {
	v1 := map[string]map[string]string{"A": {"a": "Yes"}}
	v2 := map[string]map[string]string{"A": {"a": "No"}}

	r1 := Rule{MatchPath: NewRulePath("*")}
	r2 := Rule{MatchPath: NewRulePath("@")}

	events := mustGenEvents(t, v1, v2, []Rule{r1})
	require.Len(t, events, 1)
	require.Len(t, events[0].MatchedChangePaths, 1)
	assert.Equal(t, "A.a", events[0].Change.Path)
	assert.Equal(t, "A.a", events[0].MatchedChangePaths[0])

	events = mustGenEvents(t, v1, v2, []Rule{r2})
	require.Len(t, events, 1)
	require.Len(t, events[0].MatchedChangePaths, 1)
	assert.Equal(t, "A", events[0].Change.Path)
	assert.Equal(t, "A.a", events[0].MatchedChangePaths[0])
}

func TestEvent_RulePathOption(t *testing.T) {
	v1 := map[string]map[string]string{"A": {"a": "Yes"}}
	v2 := map[string]map[string]string{"A": {"b": "No"}, "B": {"c": ""}}

	r1 := Rule{MatchPath: NewRulePath("A.{a, b}")}

	events := mustGenEvents(t, v1, v2, []Rule{r1})
	require.Len(t, events, 2)
}

// Test expects a rule with a single asterisk path to match everything.
func TestEvent_MatchEverythingRule(t *testing.T) {
	type Map map[string]any

	v1 := map[string]Map{"a": {"b": "Yes", "c": "Old"}}
	v2 := map[string]Map{"a": {"b": "No", "d": "New"}}

	r1 := Rule{MatchPath: NewRulePath("x.x")}
	r2 := Rule{MatchPath: NewRulePath("*")}
	r3 := Rule{MatchPath: NewRulePath("x.x")}

	events := mustGenEvents(t, v1, v2, []Rule{r1, r2, r3})
	require.Len(t, events, 3)

	for _, e := range events {
		assert.Equal(t, r2, e.Rule)
	}
}
