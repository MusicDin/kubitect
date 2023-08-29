package event

import (
	"fmt"
	"strings"

	"github.com/MusicDin/kubitect/pkg/utils/cmp"
)

// RulePath operators are reserved characters with distinct meanings in rule
// paths. Currently, these characters are not allowed outside their intended
// context.
//
// Note: Spaces are automatically stripped from rule paths upon creation.
const (
	// PathSeparator divides a rule's path into distinct segments.
	PathSeparator = "."

	// PathTerminator is used to prevent rule path from being matched with
	// longer change paths.
	PathTerminator = "!"

	// OptionPrefix marks the beginning of an option block.
	OptionPrefix = "{"

	// OptionSuffix marks the end of an option block.
	OptionSuffix = "}"

	// OptionSeparator divides individual options within an option block.
	OptionSeparator = ","

	// Wildcard represents a flexible match, allowing any segment in the
	// change path to be matched.
	Wildcard = "*"

	// Anchor signifies that the resulting change should be derived from
	// this specific path segment, rather than the segment where the change
	// was actually detected.
	Anchor = "@"
)

// wildcards is a list of all possible wildcard combinations within a rule
// path.
var wildcards = []string{
	Wildcard,
	Anchor,
	fmt.Sprintf("%s%s", Anchor, Wildcard),
}

// RuleType represents the priority or significance of a rule. A rule type with
// a higher value indicates greater importance. However, the 'Ignore' rule type
// is a special case and is always treated with the highest priority.
type RuleType uint8

const (
	// Allow is a normal rule type used for safe operations that do not
	// require any user confirmation.
	Allow RuleType = 0

	// Warn is a rule type used for potentially dangerous operations that
	// should request user permission in order to proceed.
	Warn RuleType = 100

	// Error is a rule type that should prevent any further actions.
	Error RuleType = 200

	// Ignore is a rule type used for ignoring specific changes.
	Ignore RuleType = 255
)

// Normalize normalizes a RuleType by mapping it to one of the known types
// (Allow, Warn, Error, Ignore). It returns the closest RuleType with smaller
// priority.
func (t RuleType) Normalize() RuleType {
	switch {
	case t == Ignore:
		return Ignore
	case t >= Error:
		return Error
	case t >= Warn:
		return Warn
	default:
		return Allow
	}
}

func (t RuleType) String() string {
	normType := t.Normalize()
	switch normType {
	case Ignore:
		return "Ignore"
	case Error:
		return "Error"
	case Warn:
		return "Warn"
	case Allow:
		return "Allow"
	}
	return "Unknown"
}

type ActionType string

const (
	Action_ScaleUp   ActionType = "scale_up"
	Action_ScaleDown ActionType = "scale_down"
)

// Rule defines the conditions that trigger events based on the detected
// changes. It specifies the rule's type, the type and path of a change it
// observes. Other fields are optional and can be used for further event
// processing.
type Rule struct {
	Type RuleType

	MatchPath       RulePath
	MatchChangeType cmp.ChangeType

	// Optional fields.
	ActionType ActionType
	Message    string
}

// Validate ensures the rule's match change type and path are valid. An error
// is returned if validation is not successful.
func (r Rule) Validate() error {
	validChangeTypes := []cmp.ChangeType{cmp.Any, cmp.None, cmp.Create, cmp.Modify, cmp.Delete}
	if !SliceContains(validChangeTypes, r.MatchChangeType) {
		return NewValidationError(r, "Invalid change type %q. Valid change types are: %v", r.MatchChangeType, validChangeTypes)
	}

	return r.MatchPath.Validate()
}

// IsOfType checks if the rule's type matches the given RuleTypes after both
// types are normalized. This is useful for comparing RuleTypes that might not
// be one of the predefined constants but should be treated as if they were.
func (r Rule) IsOfType(t RuleType) bool {
	return r.Type.Normalize() == t.Normalize()
}

// RulePath represents a rule's path, broken down into individual segments.
// It provides utilities for matching and processing change paths.
type RulePath struct {
	path     string
	segments []RulePathSegment

	// Path's metadata.
	len           int
	wildcardCount int
	anchorCount   int
	isExactPath   bool
}

// NewRulePath constructs a new rule path from a given string path. It parses
// the individual segments from the path and sets path's metadata.
func NewRulePath(path string) RulePath {
	// Create new rule path with all spaces removed.
	rp := RulePath{path: strings.ReplaceAll(path, " ", "")}

	// Check if path terminator is present at the end of the path and strip
	// it away.
	if strings.HasSuffix(rp.path, PathTerminator) {
		rp.path = strings.TrimSuffix(rp.path, PathTerminator)
		rp.isExactPath = true
	}

	segments := []RulePathSegment{}
	for _, p := range strings.Split(rp.path, PathSeparator) {
		seg := NewRulePathSegment(p)

		if seg.IsAnchor() {
			rp.anchorCount++
		}

		if seg.IsWildcard() {
			rp.wildcardCount++
		}

		segments = append(segments, seg)
	}

	rp.segments = segments
	rp.len = len(segments)
	return rp
}

// Validate ensures the rule path follows the expected format and constraints.
func (rp RulePath) Validate() error {
	if rp.path == "" {
		return NewValidationError(rp, "Path must not be empty")
	}

	if rp.anchorCount > 1 {
		return NewValidationError(rp, "Only one anchor %q is allowed in a rule path", Anchor)
	}

	for _, p := range rp.segments {
		err := p.Validate()
		if err != nil {
			return NewValidationError(rp, "%v", err)
		}
	}

	return nil
}

// Matches determines if the rule path matches with a given change path by
// comparing each segment of the rule path with the corresponding segment of
// the change path.
func (rp RulePath) Matches(changePath string) bool {
	cpParts := strings.Split(changePath, PathSeparator)

	// If change's path is shorter then rule path, then the rule definitely
	// does not apply.
	if len(cpParts) < rp.Len() {
		return false
	}

	// An exact rule path can match the change path only if they have same
	// length.
	if len(cpParts) != rp.Len() && rp.IsExactPath() {
		return false
	}

	// Ensure all corresponding segments match.
	for i, s := range rp.segments {
		if !s.Matches(cpParts[i]) {
			return false
		}
	}

	return true
}

// FindAnchorPath extracts the portion of the change path that corresponds
// to the anchor in the rule path. If no anchor exists in the rule path, it
// returns the entire change path.
func (rp RulePath) FindAnchorPath(changePath string) string {
	cpParts := strings.Split(changePath, PathSeparator)

	if len(cpParts) < rp.Len() {
		return changePath
	}

	for i, p := range rp.segments {
		if p.IsAnchor() && (p.IsWildcard() || p.ContainsOption(cpParts[i])) {
			return strings.Join(cpParts[:i+1], PathSeparator)
		}
	}

	return changePath
}

// Path returns the rule's path (without any spaces).
func (rp RulePath) Path() string {
	return rp.path
}

// Len returns the number of segments in the rule path.
func (rp RulePath) Len() int {
	return rp.len
}

// WildcardCount returns the number of wildcard segments in the rule path.
func (rp RulePath) WildcardCount() int {
	return rp.wildcardCount
}

// IsAnchorPath indicates whether the rule path contains an anchor segment.
func (rp RulePath) IsAnchorPath() bool {
	return rp.anchorCount > 0
}

// IsExactPath indicates if the rule path is terminated, ensuring an exact
// match.
func (rp RulePath) IsExactPath() bool {
	return rp.isExactPath
}

// RulePathSegment represents a segment within a rule path. It can be a simple
// string, a wildcard, an anchor, or contain multiple options.
type RulePathSegment struct {
	path string

	// Segment metadata.
	options    []string
	isAnchor   bool
	isWildcard bool
}

// NewRulePathSegment constructs a new rule path segment from a given path.
// It parses segment options if they exist, and evaluates whether the
// segment is a wildcard or/and an anchor.
func NewRulePathSegment(path string) RulePathSegment {
	rps := RulePathSegment{path: path}

	if strings.HasPrefix(path, Anchor) {
		rps.isAnchor = true
	}

	// Check if the segment is a wildcard. If it is, no further processing
	// is required.
	if SliceContains(wildcards, path) {
		rps.isWildcard = true
		return rps
	}

	// Remove anchor prefix if present.
	p := strings.TrimPrefix(path, Anchor)

	var options []string

	// Check if the segment contains option blocks.
	if strings.HasPrefix(p, OptionPrefix) && strings.HasSuffix(p, OptionSuffix) {
		p = strings.TrimPrefix(p, OptionPrefix)
		p = strings.TrimSuffix(p, OptionSuffix)
		p = strings.TrimSpace(p)

		// Extract options.
		for _, o := range strings.Split(p, OptionSeparator) {
			options = append(options, o)
		}
	}

	// If segment is a non-wildcard anchor and has no options,
	// treat the rest as a single option.
	if rps.isAnchor && len(options) == 0 {
		options = []string{strings.TrimPrefix(path, Anchor)}
	}

	rps.options = options
	return rps
}

// Validate checks the validity of a rule path segment and returns an error
// if invalid.
func (rps RulePathSegment) Validate() error {
	p := rps.path

	// Ensure the segment is not empty.
	if p == "" {
		return NewValidationError(rps, "Segment must not be empty")
	}

	// Path terminator is stripped away when rule path is created,
	// therefore no segment can contain it.
	if strings.Contains(p, PathTerminator) {
		return NewValidationError(rps, "Path terminator %q is only allowed at the end of last segment", PathTerminator)
	}

	// Validate anchor.
	if strings.Contains(p, Anchor) && !rps.IsAnchor() || strings.Count(p, Anchor) > 1 {
		return NewValidationError(rps, "Only a single anchor %q is allowed, and it must be at start", Anchor)
	}

	// Validate wildcard.
	if strings.Contains(p, Wildcard) && !rps.IsWildcard() {
		return NewValidationError(rps, "Wildcard %q can only be prefixed with an anchor %q", Wildcard, Anchor)
	}

	// Validate option block brackets.
	optsPrefixIndex := strings.Index(rps.path, OptionPrefix)
	optsSuffixIndex := strings.LastIndex(rps.path, OptionSuffix)

	if strings.Count(p, OptionPrefix) > 1 {
		return NewValidationError(rps, "Multiple option prefixes %q are not allowed", OptionPrefix)
	}

	if strings.Count(p, OptionSuffix) > 1 {
		return NewValidationError(rps, "Multiple option suffixes %q are not allowed", OptionSuffix)
	}

	if optsPrefixIndex < 0 && optsSuffixIndex >= 0 {
		return NewValidationError(rps, "Option prefix %q is missing", OptionPrefix)
	}

	if optsPrefixIndex >= 0 && optsSuffixIndex < 0 {
		return NewValidationError(rps, "Option suffix %q is missing", OptionSuffix)
	}

	if optsPrefixIndex > optsSuffixIndex {
		return NewValidationError(rps, "Option prefix %q must precede its suffix %q", OptionPrefix, OptionSuffix)
	}

	// Validate option blocks.
	if optsPrefixIndex >= 0 && optsSuffixIndex >= 0 {
		if !strings.HasSuffix(p, OptionSuffix) {
			return NewValidationError(rps, "Option suffix %q must terminate the segment", OptionSuffix)
		}

		for _, o := range rps.options {
			if o == "" {
				return NewValidationError(rps, "Options must not be empty")
			}
		}
	} else if strings.Contains(p, OptionSeparator) {
		return NewValidationError(rps, "Separator %q is only allowed inside an option block", OptionSeparator)
	}

	return nil
}

// Matches checks if the rule path segment matches a given change path segment.
// A match occurs if the segments are identical, the rule path segment is a
// wildcard, or the change path segment matches any options in the rule segment.
func (rps RulePathSegment) Matches(changePathSeg string) bool {
	return rps.path == changePathSeg || rps.IsWildcard() || rps.ContainsOption(changePathSeg)
}

// IsWildcard checks if the rule path segment is a wildcard.
func (rps RulePathSegment) IsWildcard() bool {
	return rps.isWildcard
}

// IsAnchor checks if the rule path segment is an anchor.
func (rps RulePathSegment) IsAnchor() bool {
	return rps.isAnchor
}

// ContainsOption checks if the rule path segment contains a given option.
func (rps RulePathSegment) ContainsOption(o string) bool {
	for _, s := range rps.options {
		if s == o {
			return true
		}
	}

	return false
}

// SliceContains checks whether a slice of strings contains a given key.
func SliceContains[T ~string](slice []T, key T) bool {
	for _, s := range slice {
		if s == key {
			return true
		}
	}

	return false
}
