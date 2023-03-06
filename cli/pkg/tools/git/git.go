package git

import (
	"fmt"
	"github.com/MusicDin/kubitect/cli/pkg/ui"
	"os"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Version regex ("v" any number "dot" any number "dot" any number)
var versionRegex = regexp.MustCompile("^v(\\d+)(.{1}\\d+){2}$")

type (
	GitProject interface {
		Clone(path string) error
		Url() string
		Version() string
	}

	gitProject struct {
		url     string
		version string
	}
)

func NewGitProject(url, version string) GitProject {
	return &gitProject{
		url:     url,
		version: version,
	}
}

func (p gitProject) Url() string {
	return p.url
}

func (p gitProject) Version() string {
	return p.version
}

// Clone clones a git project with the given URL and version into
// a specific directory.
func (g *gitProject) Clone(dstPath string) error {
	if len(g.url) < 1 {
		return fmt.Errorf("git clone: project URL not set")
	}

	if len(g.version) < 1 {
		return fmt.Errorf("git clone: project version not set")
	}

	// If version matches version regex, set reference name to tag,
	// otherwise set it to branch.
	var refName plumbing.ReferenceName
	if versionRegex.MatchString(g.version) {
		refName = plumbing.NewTagReferenceName(g.version)
	} else {
		refName = plumbing.NewBranchReferenceName(g.version)
	}

	opts := &git.CloneOptions{
		URL:               g.url,
		ReferenceName:     refName,
		Tags:              git.NoTags,
		RecurseSubmodules: git.NoRecurseSubmodules,
		SingleBranch:      true,
		Depth:             1,
	}

	if ui.Debug() {
		opts.Progress = ui.Streams().Out().File()
	}

	if err := os.MkdirAll(dstPath, 0700); err != nil {
		return fmt.Errorf("git clone: %v", err)
	}

	if _, err := git.PlainClone(dstPath, false, opts); err != nil {
		return fmt.Errorf("git clone: failed to clone project (url: %s, version: %s): %v", g.url, g.version, err)
	}

	return nil
}
