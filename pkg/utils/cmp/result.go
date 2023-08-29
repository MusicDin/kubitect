package cmp

// Result wraps the DiffNode and exposes only certain set of operations.
type Result struct {
	diff *DiffNode
}

// Tree returns the root node of the comparison result.
func (r Result) Tree() *DiffNode {
	return r.diff
}

// HasChanges returns true if compared values contain any difference.
func (r Result) HasChanges() bool {
	return r.diff != nil && r.diff.HasChanged()
}

// Changes extracts changes from the leaf nodes of the tree, ignoring changes
// in intermediary nodes.
func (r Result) Changes() []Change {
	if r.diff == nil {
		return nil
	}

	return r.diff.leafChanges([]Change{})
}

// DistinctChanges extracts changes from the result tree, filtering out
// propagated changes. The extraction behavior is determined by the change
// type:
//   - Create/Delete: The change is returned without further traversal into
//     child nodes, as all descendants have the same change type.
//   - Modify: If the node is a leaf, the change is returned; otherwise, the
//     function descends further into child nodes.
func (r Result) DistinctChanges() []Change {
	if r.diff == nil {
		return nil
	}

	return r.diff.distinctChanges([]Change{})
}

// ToYaml returns comparison result in YAML like format.
func (r Result) ToYaml(options ...FormatOptions) string {
	opts := FormatOptions{}
	if len(options) > 0 {
		opts = options[0]
	}

	return opts.toYaml(r.diff, 0, false)
}
