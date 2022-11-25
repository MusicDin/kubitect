package git

import (
	"cli/env"
	"cli/file"
	"cli/ui"
	"fmt"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Version regex ("v" any number "dot" any number "dot" any number)
var versionRegex = regexp.MustCompile("^v(\\d+)(.{1}\\d+){2}$")

// Clone clones a git project with the given URL and version into
// a specific directory.
func Clone(path, url, version string) error {
	if len(url) < 1 {
		return fmt.Errorf("git clone: URL not provided")
	}

	if len(version) < 1 {
		return fmt.Errorf("git clone: version not provided")
	}

	// If version matches version regex, set reference name to tag,
	// otherwise set it to branch.
	var refName plumbing.ReferenceName
	if versionRegex.MatchString(version) {
		refName = plumbing.NewTagReferenceName(version)
	} else {
		refName = plumbing.NewBranchReferenceName(version)
	}

	opts := &git.CloneOptions{
		URL:               url,
		ReferenceName:     refName,
		Tags:              git.NoTags,
		RecurseSubmodules: git.NoRecurseSubmodules,
		SingleBranch:      true,
		Depth:             1,
	}

	if env.Debug {
		opts.Progress = ui.GlobalUi().Streams.Out.File
	}

	if err := file.MakeDir(path); err != nil {
		return fmt.Errorf("git clone: %v", err)
	}

	if _, err := git.PlainClone(path, false, opts); err != nil {
		return fmt.Errorf("git clone: failed to clone project (url: %s, version: %s): %v", url, version, err)
	}

	return nil
}
