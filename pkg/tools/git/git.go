package git

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/MusicDin/kubitect/pkg/ui"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var (
	ErrInvalidRepositoryURL = errors.New("repository URL must start with https://")
	ErrCloneFailed          = errors.New("failed to clone repository")
	ErrCheckoutFailed       = errors.New("failed to checkout")
	ErrFetchingFailed       = errors.New("failed to fetch")
)

// tagRegex instructs git to clone repository by tag instead of
// branch (if matched).
var tagRegex = regexp.MustCompile("^v(\\d+)(.{1}\\d+){2}$")

type GitRepo struct {
	url        string
	version    string
	commitHash string
}

// NewGitRepo returns new instance of Git repository linked to the
// given URL.
func NewGitRepo(url string) GitRepo {
	return GitRepo{
		url: url,
	}
}

// WithRef sets the repository reference to the given branch or tag.
// Reference is used when cloning a repository.
func (r GitRepo) WithRef(branchOrTag string) GitRepo {
	r.version = branchOrTag
	return r
}

// WithCommitHash sets the repository checkout commit hash which is
// used when repository is cloned.
func (r GitRepo) WithCommitHash(commitHash string) GitRepo {
	r.commitHash = commitHash
	return r
}

func (p GitRepo) Url() string {
	return p.url
}

// Clone clones a git project with the given URL and version into
// a specific directory.
func (g GitRepo) Clone(dstPath string) error {
	if !strings.HasPrefix(g.url, "https://") {
		return ErrInvalidRepositoryURL
	}

	opts := &git.CloneOptions{
		URL:               g.url,
		Tags:              git.NoTags,
		RecurseSubmodules: git.NoRecurseSubmodules,
		SingleBranch:      true,
	}

	if g.version != "" {
		// If version matches version regex, set reference
		// name to tag, otherwise set it to branch.
		opts.ReferenceName = plumbing.NewBranchReferenceName(g.version)
		if tagRegex.MatchString(g.version) {
			opts.ReferenceName = plumbing.NewTagReferenceName(g.version)
		}
	}

	if g.commitHash == "" {
		opts.Depth = 1
	}

	if ui.Debug() {
		opts.Progress = ui.Streams().Out().File()
	}

	// Ensure destination directory exists.
	err := os.MkdirAll(dstPath, 0700)
	if err != nil {
		return err
	}

	// Clone repository.
	repo, err := git.PlainClone(dstPath, false, opts)
	if err != nil {
		return fmt.Errorf("%w (url: %s, version: %s): %v", ErrCloneFailed, g.url, g.version, err)
	}

	if g.commitHash != "" {
		// Fetch repository work tree.
		tree, err := repo.Worktree()
		if err != nil {
			return err
		}

		opts := &git.CheckoutOptions{
			Hash: plumbing.NewHash(g.commitHash),
		}

		// Checkout to specific commit hash.
		err = tree.Checkout(opts)
		if err != nil {
			return fmt.Errorf("%w (commitHash: %s): %v", ErrCheckoutFailed, g.commitHash, err)
		}
	}

	return nil
}
