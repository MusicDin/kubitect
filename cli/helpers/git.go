package helpers

import (
	"cli/env"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GitClone clones a git project with the given URL and version into
// a specific directory.
func GitClone(path string, url string, version string) error {

	if len(path) < 1 {
		return fmt.Errorf("Git clone destination not provided.")
	}

	if len(url) < 1 {
		return fmt.Errorf("Git URL not provided.")
	}

	if len(version) < 1 {
		return fmt.Errorf("Git version not provided.")
	}

	gitCloneOptions := &git.CloneOptions{
		URL:               url,
		ReferenceName:     plumbing.NewBranchReferenceName(version),
		Tags:              git.NoTags,
		RecurseSubmodules: git.NoRecurseSubmodules,
		SingleBranch:      true,
		Depth:             1,
	}

	if env.DebugMode {
		gitCloneOptions.Progress = os.Stdout
	}

	_, err := git.PlainClone(path, false, gitCloneOptions)
	if err != nil {
		return fmt.Errorf("Error cloning project for 'url=%s' and 'version=%s': %w", url, version, err)
	}

	return nil
}
