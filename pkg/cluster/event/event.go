package event

import (
	"fmt"

	"github.com/MusicDin/kubitect/pkg/utils/cmp"
)

// Event represents a detected change and its associated rule.
type Event struct {
	Rule   Rule
	Change cmp.Change

	// Paths of changes that matched the rule.
	MatchedChangePaths []string
}

func (e Event) String() string {
	return fmt.Sprintf("(%s) Change: [Type: %s, Path: %s]", e.Rule.Type, e.Change.Type, e.Change.Path)
}

// Events abstracts Event list in order to provide more advance operations on
// the list.
type Events []Event

// Filter iterates over the events and returns a new slice of events containing
// only those events that satisfy the provided filter function (return true).
func (events Events) Filter(filter func(event Event) bool) Events {
	filtered := []Event{}
	for _, e := range events {
		if filter(e) {
			filtered = append(filtered, e)
		}
	}

	return filtered
}

// FilterByRuleType returns a new slice containing only those events that
// match the provided RuleType.
func (events Events) FilterByRuleType(t RuleType) Events {
	filter := func(e Event) bool {
		return e.Rule.Type == t
	}

	return events.Filter(filter)
}

// FilterByActionType returns a new slice containing only those events that
// match the provided ActionType.
func (events Events) FilterByAction(t ActionType) Events {
	filter := func(e Event) bool {
		return e.Rule.ActionType == t
	}

	return events.Filter(filter)
}

// GenerateEvents evaluates the changes from the comparison tree against the
// provided rules and returns a list of corresponding events. Each event
// encapsulates a matched change and its associated rule. A single change can
// match at most one rule and thus produce at most one event.
// Note that provided rules are validated prior the event generation.
func GenerateEvents(node *cmp.DiffNode, rules []Rule) ([]Event, error) {
	for _, r := range rules {
		err := r.Validate()
		if err != nil {
			return nil, err
		}
	}

	return generateEvents(node, rules, []Event{}), nil
}

func generateEvents(node *cmp.DiffNode, rules []Rule, events []Event) []Event {
	if node == nil {
		return events
	}

	if node.IsLeaf() && node.HasChanged() {
		rule := matchRule(node, rules)
		if rule != nil {
			events = createAndAddEvent(node, rule, events)
		}
	}

	for _, c := range node.Children() {
		events = generateEvents(c, rules, events)
	}

	return events
}

// createAndAddEvent constructs an event based on the provided node and rule.
// If the rule's path contains an anchor, the change is extracted from the
// corresponding parent node. In such case, the function checks whether an
// event with the same rule, change, and rule path already exists in the list.
// If it does, it only appends the node's path to the MatchedChangePaths of an
// existing event. Otherwise, it appends the new event to the list.
func createAndAddEvent(node *cmp.DiffNode, rule *Rule, events []Event) []Event {
	targetNode := node
	rulePath := rule.MatchPath

	// If the rule's path contains an anchor, determine the anchor path
	// and find the corresponding node.
	if rulePath.IsAnchorPath() {
		anchorNodePath := rulePath.FindAnchorPath(node.Path())
		anchorNode := node.ParentByPath(anchorNodePath)
		if anchorNode != nil {
			targetNode = anchorNode
		}
	}

	event := &Event{
		Rule:               *rule,
		Change:             targetNode.ToChange(),
		MatchedChangePaths: []string{node.Path()},
	}

	// If the rule's path contains an anchor, try to group the new event
	// with an exiting event. Events are grouped if their rule, change and
	// rule path match.
	if rulePath.IsAnchorPath() {
		for i, e := range events {
			isSameRulePath := e.Rule.MatchPath.Path() == event.Rule.MatchPath.Path()
			isSameRuleType := e.Rule.Type == event.Rule.Type
			isSameChangeType := e.Change.Type == event.Change.Type

			if isSameRuleType && isSameChangeType && isSameRulePath {
				events[i].MatchedChangePaths = append(events[i].MatchedChangePaths, node.Path())
				return events
			}
		}
	}

	// If no matching event is found or the rule's path doesn't contain
	// an anchor, append the new event to the list.
	return append(events, *event)
}

// matchRule determines the most appropriate rule for a given change. A change
// matches a rule if both the path and change type align. Each change can
// correspond to either no rule or one rule. If multiple rules match, they are
// prioritized by path length, wildcard count, and rule priority. If no rule
// matches, it returns nil.
func matchRule(node *cmp.DiffNode, rules []Rule) *Rule {
	var bestMatch *Rule
	for i, rule := range rules {
		if rule.MatchChangeType != node.ChangeType() && rule.MatchChangeType != cmp.Any {
			continue
		}

		if rule.MatchPath.Matches(node.Path()) {
			if bestMatch == nil || isBetterMatch(rule, *bestMatch) {
				bestMatch = &rules[i]
			}
		}

	}

	return bestMatch
}

// isBetterMatch determines whether rule r1 is a better match then rule r2
// based on the following rules:
//   1. Path Length: Longer paths are prioritized.
//   2. Wildcard Count: Paths with fewer wildcards are deemed more specific
//	and are prioritized.
//   3. Rule Priority: Higher priority rules are prioritized.
func isBetterMatch(r1, r2 Rule) bool {
	r1Len := r1.MatchPath.Len()
	r2Len := r2.MatchPath.Len()

	// Longer path has precedence.
	if r1Len != r2Len {
		return r1Len > r2Len
	}

	r1Wcs := r1.MatchPath.WildcardCount()
	r2Wcs := r2.MatchPath.WildcardCount()

	// Path with fewer wildcards has precedence.
	if r1Wcs != r2Wcs {
		return r1Wcs < r2Wcs
	}

	// Higher priority has precedence.
	return r1.Type > r2.Type
}
