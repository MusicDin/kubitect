package git

import (
	"cli/env"
	"cli/ui"
	"fmt"
	"os"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Clone clones a git project with the given URL and version into
// a specific directory.
func Clone(path, url, version string) error {
	if len(path) < 1 {
		return fmt.Errorf("git clone: destination not provided.")
	}

	if len(url) < 1 {
		return fmt.Errorf("git clone: URL not provided.")
	}

	if len(version) < 1 {
		return fmt.Errorf("git clone: version not provided.")
	}

	// Version regex ("v" any number "dot" any number "dot" any number)
	versionRegex, err := regexp.Compile("^v(\\d+)(.{1}\\d+){2}$")
	if err != nil {
		return err
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

	os.MkdirAll(path, os.ModePerm)

	_, err = git.PlainClone(path, false, opts)

	if err != nil {
		return fmt.Errorf("git clone: failed to clone project (url: %s, version: %s): %v", url, version, err)
	}

	return nil
}
