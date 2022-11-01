package git

import (
	"cli/env"
	"fmt"
	"os"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Clone clones a git project with the given URL and version into
// a specific directory.
func Clone(path string, url string, version string) error {

	if len(path) < 1 {
		return fmt.Errorf("Git clone destination not provided.")
	}

	if len(url) < 1 {
		return fmt.Errorf("Git URL not provided.")
	}

	if len(version) < 1 {
		return fmt.Errorf("Git version not provided.")
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

	gitCloneOptions := &git.CloneOptions{
		URL:               url,
		ReferenceName:     refName,
		Tags:              git.NoTags,
		RecurseSubmodules: git.NoRecurseSubmodules,
		SingleBranch:      true,
		Depth:             1,
	}

	if env.DebugMode {
		gitCloneOptions.Progress = os.Stdout
	}

	os.MkdirAll(path, os.ModePerm)

	_, err = git.PlainClone(path, false, gitCloneOptions)

	if err != nil {
		return fmt.Errorf("Error cloning git project for 'url=%s' and 'version=%s': %v", url, version, err)
	}

	return nil
}
