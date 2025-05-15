package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/StevenCyb/autosemver/internal/generator"
	"github.com/StevenCyb/autosemver/internal/logger"
	"github.com/StevenCyb/autosemver/internal/model"
)

const version = "1.0.0"

var errorExitCode = 1
var ignoreInvalidTags = false
var asRC = false
var log logger.Logger = logger.Silent{}
var conventionalCommitToSemVer = []model.Tuple[string, string]{
	{First: "breaking change", Second: "major"},
	{First: "fix!", Second: "major"},
	{First: "feat!", Second: "major"},
	{First: "feat", Second: "minor"},
	{First: "pref", Second: "patch"},
	{First: "fix", Second: "patch"},
}

func main() {
	repoPath := "."
	if len(os.Args) > 1 {
		args := os.Args[1:]

		if args[0] == "version" {
			fmt.Println(version)
			os.Exit(0)
		} else if args[0] == "help" {
			printHelp()
			os.Exit(0)
		}

		if !strings.HasPrefix(args[0], "-") {
			repoPath = args[0]
			args = args[1:]
			if _, err := os.Stat(repoPath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error: path '%s' does not exist\n", repoPath)
				os.Exit(errorExitCode)
			}
		}

		for _, arg := range args {
			if arg == "--disable-exit-1" {
				errorExitCode = 0
			} else if arg == "--verbose" || arg == "-v" {
				log = logger.Verbose{}
			} else if arg == "--release-candidate" || arg == "-r" {
				asRC = true
			} else if arg == "--ignore-invalid-tag" || arg == "-i" {
				ignoreInvalidTags = true
			} else if arg == "--help" || arg == "-h" {
				printHelp()
				os.Exit(0)
			} else if strings.HasPrefix(arg, "--mapping=") || strings.HasPrefix(arg, "-m=") {
				mapping := strings.TrimPrefix(arg, "--mapping=")
				mapping = strings.TrimPrefix(mapping, "-m=")
				splitMapping := strings.Split(mapping, ":")
				if len(splitMapping) != 2 || len(splitMapping[0]) == 0 || len(splitMapping[1]) == 0 ||
					(splitMapping[1] != "major" && splitMapping[1] != "minor" && splitMapping[1] != "patch") {
					fmt.Fprintf(os.Stderr, "Error: invalid mapping format '%s'\n", mapping)
					printHelp()
					os.Exit(errorExitCode)
				}
				conventionalCommitToSemVer = append(conventionalCommitToSemVer, model.Tuple[string, string]{First: splitMapping[0], Second: splitMapping[1]})
			} else {
				fmt.Fprintf(os.Stderr, "Error: unknown option '%s'\n", arg)
				printHelp()
				os.Exit(errorExitCode)
			}
		}
	}

	if !asRC {
		version, err := generator.FindNextVersion(repoPath, conventionalCommitToSemVer, log, ignoreInvalidTags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(errorExitCode)
		}
		fmt.Println(*version)
	} else {
		version, err := generator.FindNextRC(repoPath, conventionalCommitToSemVer, log, ignoreInvalidTags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(errorExitCode)
		}
		fmt.Println(*version)
	}
}

func printHelp() {
	fmt.Println("Usage: autosemver {version|help|[repository_path]} [options]")
	fmt.Println("\nCommands:")
	fmt.Println("\t[repository_path]: path to the git repository (default: current directory)")
	fmt.Println("\tversion: show the version of autosemver")
	fmt.Println("\thelp: show this help message")
	fmt.Println("\nOptions:")
	fmt.Println("\t--help, -h: show this help message")
	fmt.Println("\t--verbose, -v: enable verbose output")
	fmt.Println("\t--release-candidate, -r: mark the version as a release candidate (append '-rc.N' to the version)")
	fmt.Println("\t--ignore-invalid-tag, -i: ignore invalid tags (not a valid semantic version)")
	fmt.Println("\t--disable-exit-1: do not exit with a non-zero code on error")
	fmt.Println("\t--mapping=feat:minor, -m=fix:patch: add mapping for commit types (prefix) to version increments {major, minor, patch}")
	fmt.Println("\nDefault Mapping (ignores not matching commits):")
	for _, i := range conventionalCommitToSemVer {
		fmt.Printf("\t\"%s\": %s\n", i.First, i.Second)
	}
}
