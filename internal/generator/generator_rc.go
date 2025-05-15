package generator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/StevenCyb/autosemver/internal/logger"
	"github.com/StevenCyb/autosemver/internal/model"
	"github.com/StevenCyb/autosemver/internal/utils"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func FindNextRC(repositoryPath string, incMapping []model.Tuple[string, string], log logger.Logger, ignoreInvalidTags bool) (*string, error) {
	log.Printf("Finding next version in %s\n", repositoryPath)
	repo, err := git.PlainOpen(repositoryPath)
	if err != nil {
		return nil, err
	}
	return findNextRC(repo, incMapping, log, ignoreInvalidTags)
}

func findNextRC(repo *git.Repository, incMapping []model.Tuple[string, string], log logger.Logger, ignoreInvalidTags bool) (*string, error) {
	log.Printf("Finding latest version tag")
	tagRefs, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	var latestVersionTag *model.Tuple[model.SemVer, string]
	err = tagRefs.ForEach(func(ref *plumbing.Reference) error {
		version := ref.Name().Short()
		commitId, err := repo.ResolveRevision(plumbing.Revision(ref.Name()))
		if err != nil {
			return nil
		}

		splitVersion := regexp.MustCompile(`^(?<major>[0-9]+)\.(?<minor>[0-9]+)\.(?<patch>[0-9]+)(?<rc>-rc\.[0-9]+)?$`).FindStringSubmatch(version)
		if len(splitVersion) != 5 {
			if ignoreInvalidTags {
				log.Printf("Tag %s is not a valid semantic version, ignoring\n", version)
				return nil
			} else {
				return errors.New(fmt.Sprintf("Tag %s is not a valid semantic version", version))
			}
		}
		major := utils.MustParseUint(splitVersion[1])
		minor := utils.MustParseUint(splitVersion[2])
		patch := utils.MustParseUint(splitVersion[3])
		hasRC := splitVersion[4] != ""
		var rc *uint = nil
		if hasRC {
			rc = new(uint)
			*rc = utils.MustParseUint(strings.TrimPrefix(splitVersion[4], "-rc."))
		}

		if latestVersionTag == nil ||
			latestVersionTag.First.Major < major ||
			(latestVersionTag.First.Major == major && latestVersionTag.First.Minor < minor) ||
			(latestVersionTag.First.Major == major && latestVersionTag.First.Minor == minor && latestVersionTag.First.Patch < patch) ||
			(hasRC && latestVersionTag.First.Major == major && latestVersionTag.First.Minor == minor && latestVersionTag.First.Patch == patch && *latestVersionTag.First.RC < *rc) {
			latestVersionTag = &model.Tuple[model.SemVer, string]{
				First:  model.SemVer{Major: major, Minor: minor, Patch: patch, RC: rc},
				Second: commitId.String(),
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if latestVersionTag != nil {
		if latestVersionTag.First.RC != nil {
			log.Printf("Latest version tag: %d.%d.%d-rc.%d\n", latestVersionTag.First.Major, latestVersionTag.First.Minor, latestVersionTag.First.Patch, *latestVersionTag.First.RC)
			*latestVersionTag.First.RC++
			newRcVersion := fmt.Sprintf("%d.%d.%d-rc.%d", latestVersionTag.First.Major, latestVersionTag.First.Minor, latestVersionTag.First.Patch, *latestVersionTag.First.RC)
			return &newRcVersion, nil
		} else {
			log.Printf("Latest version tag: %d.%d.%d\n", latestVersionTag.First.Major, latestVersionTag.First.Minor, latestVersionTag.First.Patch)
		}
	} else {
		log.Println("No version tag found")
	}

	incMajor := false
	incMinor := false
	incPatch := false

	log.Println("Finding commits since latest version tag")
	headRef, err := repo.Head()
	if err != nil {
		return nil, err
	}
	commitIter, err := repo.Log(&git.LogOptions{From: headRef.Hash()})
	if err != nil {
		return nil, err
	}
	commitIter.ForEach(func(c *object.Commit) error {
		if latestVersionTag != nil && latestVersionTag.Second == c.Hash.String() {
			return errors.New("END")
		}

		msg := strings.ReplaceAll(strings.ToLower(c.Message), "\n", " ")
		for _, mapping := range incMapping {
			log.Printf("Commit: [%s] %s\n", c.Hash.String(), msg)
			if strings.HasPrefix(msg, mapping.First) {
				switch mapping.Second {
				case "major":
					log.Printf("Found major version bump commit %s\n", c.Hash.String())
					incMajor = true
					break
				case "minor":
					log.Printf("Found minor version bump commit %s\n", c.Hash.String())
					incMinor = true
					break
				case "patch":
					log.Printf("Found patch version bump commit %s\n", c.Hash.String())
					incPatch = true
					break
				}
			}
		}

		return nil
	})

	if latestVersionTag == nil {
		latestVersionTag = &model.Tuple[model.SemVer, string]{
			First:  model.SemVer{Major: 0, Minor: 0, Patch: 0},
			Second: "",
		}
	}
	switch {
	case incMajor:
		latestVersionTag.First.Major++
		latestVersionTag.First.Minor = 0
		latestVersionTag.First.Patch = 0
	case incMinor:
		latestVersionTag.First.Minor++
		latestVersionTag.First.Patch = 0
	case incPatch:
		latestVersionTag.First.Patch++
	}

	newRcVersion := fmt.Sprintf("%d.%d.%d-rc.1", latestVersionTag.First.Major, latestVersionTag.First.Minor, latestVersionTag.First.Patch)

	return &newRcVersion, nil
}
