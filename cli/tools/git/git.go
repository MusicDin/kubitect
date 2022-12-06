package git

import (
	"cli/file"
	"cli/ui"
	"fmt"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Version regex ("v" any number "dot" any number "dot" any number)
var versionRegex = regexp.MustCompile("^v(\\d+)(.{1}\\d+){2}$")

type GitProject struct {
	Url     string
	Version string
	Path    string

	Ui *ui.Ui
}

// Clone clones a git project with the given URL and version into
// a specific directory.
func (g *GitProject) Clone() error {
	if len(g.Url) < 1 {
		return fmt.Errorf("git clone: project URL not set")
	}

	if len(g.Version) < 1 {
		return fmt.Errorf("git clone: project version not set")
	}

	// If version matches version regex, set reference name to tag,
	// otherwise set it to branch.
	var refName plumbing.ReferenceName
	if versionRegex.MatchString(g.Version) {
		refName = plumbing.NewTagReferenceName(g.Version)
	} else {
		refName = plumbing.NewBranchReferenceName(g.Version)
	}

	opts := &git.CloneOptions{
		URL:               g.Url,
		ReferenceName:     refName,
		Tags:              git.NoTags,
		RecurseSubmodules: git.NoRecurseSubmodules,
		SingleBranch:      true,
		Depth:             1,
	}

	if g.Ui.Debug {
		opts.Progress = g.Ui.Streams.Out.File
	}

	if err := file.MakeDir(g.Path); err != nil {
		return fmt.Errorf("git clone: %v", err)
	}

	if _, err := git.PlainClone(g.Path, false, opts); err != nil {
		return fmt.Errorf("git clone: failed to clone project (url: %s, version: %s): %v", g.Url, g.Version, err)
	}

	return nil
}
