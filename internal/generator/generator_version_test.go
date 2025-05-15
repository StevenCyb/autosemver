package generator

import (
	"testing"

	"github.com/StevenCyb/autosemver/internal/logger"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
)

func TestFindNextVersion_EmptyRepository(t *testing.T) {
	t.Parallel()

	repo, _ := NewSimulatedRepository(t)
	tag, err := findNextVersion(repo, conventionalCommitToSemVer, logger.Silent{}, false)

	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "0.0.0", *tag)
}

func TestFindNextVersion_NoTag_PatchCommit(t *testing.T) {
	t.Parallel()

	repo, fs := NewSimulatedRepository(t)
	fakeCommit(t, repo, fs, "main.go", "fix: fix a bug")
	tag, err := findNextVersion(repo, conventionalCommitToSemVer, logger.Silent{}, false)

	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "0.0.1", *tag)
}

func TestFindNextVersion_NoTag_FeatCommit(t *testing.T) {
	t.Parallel()

	repo, fs := NewSimulatedRepository(t)
	fakeCommit(t, repo, fs, "main.go", "feat: some new feature")
	tag, err := findNextVersion(repo, conventionalCommitToSemVer, logger.Silent{}, false)

	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "0.1.0", *tag)
}

func TestFindNextVersion_NoTag_BreakingChangeCommit(t *testing.T) {
	t.Parallel()

	repo, fs := NewSimulatedRepository(t)
	fakeCommit(t, repo, fs, "main.go", "feat!: some new feature")
	tag, err := findNextVersion(repo, conventionalCommitToSemVer, logger.Silent{}, false)

	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "1.0.0", *tag)
}

func TestFindNextVersion_Tag1_0_0_PatchCommit(t *testing.T) {
	t.Parallel()

	repo, fs := NewSimulatedRepository(t)
	headRef, err := repo.Reference(plumbing.HEAD, true)
	assert.NoError(t, err)
	_, err = repo.CreateTag("1.0.0", headRef.Hash(), nil)
	assert.NoError(t, err)
	fakeCommit(t, repo, fs, "main.go", "fix: fix a bug")
	tag, err := findNextVersion(repo, conventionalCommitToSemVer, logger.Silent{}, false)

	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "1.0.1", *tag)
}

func TestFindNextVersion_Tag1_0_0_FeatCommit(t *testing.T) {
	t.Parallel()

	repo, fs := NewSimulatedRepository(t)
	headRef, err := repo.Reference(plumbing.HEAD, true)
	assert.NoError(t, err)
	_, err = repo.CreateTag("1.0.0", headRef.Hash(), nil)
	assert.NoError(t, err)
	fakeCommit(t, repo, fs, "main.go", "feat: some new feature")
	tag, err := findNextVersion(repo, conventionalCommitToSemVer, logger.Silent{}, false)

	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "1.1.0", *tag)
}

func TestFindNextVersion_Tag1_0_0_BreakingChangeCommit(t *testing.T) {
	t.Parallel()

	repo, fs := NewSimulatedRepository(t)
	headRef, err := repo.Reference(plumbing.HEAD, true)
	assert.NoError(t, err)
	_, err = repo.CreateTag("1.0.0", headRef.Hash(), nil)
	assert.NoError(t, err)
	fakeCommit(t, repo, fs, "main.go", "feat!: some new feature")
	tag, err := findNextVersion(repo, conventionalCommitToSemVer, logger.Silent{}, false)

	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "2.0.0", *tag)
}

func TestFindNextVersion_Tag1_0_0_rc_BreakingChangeCommit(t *testing.T) {
	t.Parallel()

	repo, fs := NewSimulatedRepository(t)
	headRef, err := repo.Reference(plumbing.HEAD, true)
	assert.NoError(t, err)
	_, err = repo.CreateTag("1.0.0-rc.1", headRef.Hash(), nil)
	assert.NoError(t, err)
	fakeCommit(t, repo, fs, "main.go", "feat!: some new feature")
	tag, err := findNextVersion(repo, conventionalCommitToSemVer, logger.Silent{}, false)

	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "1.0.0", *tag)
}

func TestFindNextVersion_InvalidTag_NoCommit(t *testing.T) {
	t.Parallel()

	repo, _ := NewSimulatedRepository(t)
	headRef, err := repo.Reference(plumbing.HEAD, true)
	assert.NoError(t, err)
	_, err = repo.CreateTag("invalid", headRef.Hash(), nil)
	assert.NoError(t, err)
	tag, err := findNextVersion(repo, conventionalCommitToSemVer, logger.Silent{}, false)

	assert.Error(t, err)
	assert.Nil(t, tag)
}

func TestFindNextVersion_InvalidTagButIgnored_BreakingChangeCommit(t *testing.T) {
	t.Parallel()

	repo, fs := NewSimulatedRepository(t)
	headRef, err := repo.Reference(plumbing.HEAD, true)
	assert.NoError(t, err)
	_, err = repo.CreateTag("1.0.0-rc.1", headRef.Hash(), nil)
	assert.NoError(t, err)
	fakeCommit(t, repo, fs, "main.go", "feat!: some new feature")
	tag, err := findNextVersion(repo, conventionalCommitToSemVer, logger.Silent{}, true)

	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "1.0.0", *tag)
}
