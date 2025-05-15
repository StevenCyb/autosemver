package generator

import (
	"testing"
	"time"

	"github.com/StevenCyb/autosemver/internal/model"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
)

var conventionalCommitToSemVer = []model.Tuple[string, string]{
	{First: "breaking change", Second: "major"},
	{First: "fix!", Second: "major"},
	{First: "feat!", Second: "major"},
	{First: "feat", Second: "minor"},
	{First: "pref", Second: "patch"},
	{First: "fix", Second: "patch"},
}

func NewSimulatedRepository(t *testing.T) (*git.Repository, billy.Filesystem) {
	t.Helper()

	fs := memfs.New()
	storer := memory.NewStorage()

	repo, err := git.InitWithOptions(storer, fs, git.InitOptions{
		DefaultBranch: plumbing.NewBranchReferenceName("main"),
	})
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	fakeCommit(t, repo, fs, "README.md", "init")

	ref, err := repo.Head()
	assert.NoError(t, err)
	assert.Equal(t, plumbing.NewBranchReferenceName("main"), ref.Name())

	return repo, fs
}

func fakeCommit(t *testing.T, repo *git.Repository, fs billy.Filesystem, fileName, commitMessage string) {
	wt, err := repo.Worktree()
	assert.NoError(t, err)

	f, err := fs.Create(fileName)
	assert.NoError(t, err)
	_, err = f.Write([]byte("content"))
	assert.NoError(t, err)

	_, err = wt.Add(fileName)
	assert.NoError(t, err)

	_, err = wt.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test Bot",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	assert.NoError(t, err)
}
