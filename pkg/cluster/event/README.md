# Event Generation Behavior

The `event` package is designed to detect changes and generate events based on the provided rules.

## Event

An `Event` represents a detected change and its associated rule.
It contains:

- The rule that was matched.
- The detected change.
- Paths of changes that matched the rule (triggered an event).

## Rule

A `Rule` defines the conditions that trigger events based on detected changes.
It specifies:

- The rule's type.
- The rule's path.
- The type of change it observes.

## Event Generation

The core of the event generation process is the `GenerateEvents` function.
It evaluates the changes from the comparison tree against the provided rules and returns a list of corresponding events.
A single change can match at most one rule and thus produce at most one event.

The process involves:

- Validating the provided rules.
- Recursively traversing through each node in the comparison tree.
- Matching each leaf change against the provided rules and determining the best match.
- Creating and adding events based on best matched rule.

## Matching Rules

The matchRule function evaluates each change to identify the best matching rule.
A change matches a rule if both the path and change type align.

For each change, the `matchRule` function determines the most appropriate rule (best match).
A change matches a rule if both the path and change type align.
If no rules match with a given change, the change is ignored.
If multiple rules match, they are prioritized based on:

1. Path Length: Longer paths are prioritized.
2. Wildcard Count: Paths with fewer wildcards are considered more specific and are prioritized.
3. Rule Priority: Higher priority rules are prioritized.

# RulePath

To improve rule matching flexibility, the rule's path uses specific operators.
While these resemble regex, they are much simpler and less powerful.

## Operators

- `.`: Separates path into multiple segments.
- `{}`: Option block defines multiple options divided with comma for matching path segments.
- `*`: Acts as a wildcard, matching any path segment.
- `@`: Denotes an anchor for a path segment.
- `!`: Indicates path termination, ensuring an exact match.

## Simplifications

- `{a}` can be simplified to `a` when the option block contains only one option.
- `@{a}` can be simplified to `@a` when using an anchor with a single segment option.
- `@*` can be simplified to `@` when using an anchor as a wildcard.
- Spaces within rules are ignored and can be omitted.

## Reserved Characters

The characters `*`, `@`, `{`, `}`, `.`, `!` and space (` `) are considered reserved.
Rule paths containing these characters outside their intended context will be considered invalid.
Therefore, change paths containing these characters will never be matched.

## Matching Scenarios

Consider changes on the following paths:
- `a`
- `a.b`
- `a.b.c`
- `a.B.c`
- `a.C.c`
- `A.B.c`
- `a.b.c.d`

1. **Normal Path**:
   - **Description**: Normally, rule path matches change paths that begin with the same segments.
   - **Example**: Rule path `a.b.c` matches both `a.b.c` and `a.b.c.d`.

2. **Exact Match**:
   - **Description**: For a precise match without considering extended paths, end the rule path with `!`.
   - **Example**: Rule path `a.b.c!` strictly matches `a.b.c`.

3. **Path Options**:
   - **Description**: Option blocks in rule paths allow matching multiple potential segments.
   - **Example**: Rule path `a.{b, B}.c` can match either `a.b.c` or `a.B.c`.

4. **Wildcard Path**:
   - **Description**: Wildcards allow matching any segment in the change path.
   - **Example**: Rule path `a.*.c` can match any of `a.b.c`, `a.B.c`, `a.C.c`, or `a.b.c.d`.

5. **Anchor Path**:
   - **Description**: Anchors ensure that the change in an event is extracted from the anchored path segment, rather than the segment where the change was actually detected.
   - **Examples**:
     - Rule path `a.@*.c` (or its shorthand `a.@.c`) matches paths like `a.b.c`, `a.B.c`, `a.C.c`, and `a.b.c.d`.
     - Rule path `a.@{b}.c` (or its shorthand `a.@b.c`) specifically matches `a.b.c` and its extension `a.b.c.d`.
