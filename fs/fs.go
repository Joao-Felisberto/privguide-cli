// Package to abstract file system accesses,
// namely by handling lookup from both the global and local directories
//
// By default, the local path is `.devprivops/` and the global path is `/etc/devprivops/`.
// Files in the local path override those in the global directory.
//
// The unexported functions are independent of the local and global directories and are made to increase
// testability. These are the ones that should be targeted in unit tests and thus are exported in `export_test.go`.
//
// This package only supports UNIX paths
package fs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/Joao-Felisberto/devprivops/util"
)

/*
	lookup order:
	1. /etc/appName
	2. .appName/
*/

var (
	LocalDir  = fmt.Sprintf("./.%s", util.AppName)   // The local directory
	GlobalDir = fmt.Sprintf("/etc/%s", util.AppName) // The global directory
)

// Returns the path of a file relative to the local or global root using the pre-determined paths to the local and global directories
//
// `relativePath`: the path relative to either root
//
// returns: the path to the provided file relative to the root it is in, or an error if reading any of the directories fails.
func GetFile(relativePath string) (string, error) {
	return getFile(
		relativePath,
		LocalDir,
		GlobalDir,
	)
}

// Returns the path of a file relative to the local or global root using the provided paths to the local and global directories
//
// `localRoot`: the root of the local directory
//
// `globalRoot`: the root of the global directory
//
// `relativePath` the path relative to either root
//
// returns: the path to the provided file relative to the root it is in, or an error if reading any of the directories fails.
func getFile(relativePath string, localRoot string, globalRoot string) (string, error) {
	localPath := fmt.Sprintf("%s/%s", localRoot, relativePath)
	if _, err := os.Stat(localPath); errors.Is(err, os.ErrNotExist) {
		defaultPath := fmt.Sprintf("%s/%s", globalRoot, relativePath)
		if _, err := os.Stat(defaultPath); errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		return defaultPath, nil
	}
	return localPath, nil
}

// Returns the paths of the system descriptions relative to their respective root using the default paths to the local and global directories
//
// `relativePath` the path relative to either root
//
// returns: the relative paths of the system descriptions, or an error if reading any of the directories fails.
func GetDescriptions(descriptionRoot string) ([]string, error) {
	return getDescriptions(
		descriptionRoot,
		LocalDir,
		GlobalDir,
	)
}

// Returns the paths of the system descriptions relative to their respective root using the provided paths to the local and global directories
//
// `localRoot`: the root of the local directory
//
// `globalRoot`: the root of the global directory
//
// `relativePath` the path relative to either root
//
// returns: the relative paths of the system descriptions, or an error if reading any of the directories fails.
func getDescriptions(descriptionRoot string, localRoot string, globalRoot string) ([]string, error) {
	localPath := fmt.Sprintf("%s/%s/", localRoot, descriptionRoot)
	globalPath := fmt.Sprintf("%s/%s/", globalRoot, descriptionRoot)

	files := []string{}

	entries, err := os.ReadDir(localPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("error reading local directory: %s", err)
	}

	for _, e := range entries {
		files = append(files, fmt.Sprintf("%s/%s", descriptionRoot, e.Name()))
	}

	entries, err = os.ReadDir(globalPath)
	if err != nil {
		return files, nil
	}

	for _, e := range entries {
		files = append(files, fmt.Sprintf("%s/%s", descriptionRoot, e.Name()))
	}

	return files, nil
}

// Returns the directory names of the system regulation directories under `regulations/` using the default paths to the local and global directories
//
// returns: the directory names of the system regulation directories, or an error if reading any of the directories fails.
func GetRegulations() ([]string, error) {
	return getRegulations(
		LocalDir,
		GlobalDir,
	)
}

// Returns the directory names of the system regulation directories under `regulations/` using the default paths to the local and global directories
//
// `localRoot`: the root of the local directory
//
// `globalRoot`: the root of the global directory
//
// returns: the directory names of the system regulation directories, or an error if reading any of the directories fails.
func getRegulations(localRoot string, globalRoot string) ([]string, error) {
	localPath := fmt.Sprintf("%s/regulations/", localRoot)
	defaultPath := fmt.Sprintf("%s/regulations/", globalRoot)

	files := []string{}

	localRegulations, err := getDirsInDir(localPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	defaultRegulations, err := getDirsInDir(defaultPath)
	if err != nil {
		files = append(files, localRegulations...)

		return files, nil
	}

	files = append(files, localRegulations...)
	files = append(files, defaultRegulations...)

	return files, nil
}

// Find all top level directories inside a directory
//
// `path`: The parent directory of which we want to know the subdirectories
//
// returns: The list of subdirectories, or an error if reading any of the directories fails.
func getDirsInDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return util.Map(
		util.Filter(entries, func(de fs.DirEntry) bool { return de.IsDir() }),
		func(de fs.DirEntry) string { return de.Name() },
	), nil
}

// Find the relative directories of each configuration file.
// The returned directories contain the root
//
// returns: The list of configuration files, or an error if reading any of the directories fails.
func GetConfigs() ([]string, error) {
	return getConfigs(
		LocalDir,
		GlobalDir,
	)
}

// Find the relative directories of each configuration file.
// The returned directories contain the root
//
// `localRoot`: the root of the local directory
//
// `globalRoot`: the root of the global directory
//
// returns: The list of configuration files, or an error if reading any of the directories fails.
func getConfigs(localRoot string, globalRoot string) ([]string, error) {
	localPath := fmt.Sprintf("%s/config/", localRoot)
	globalPath := fmt.Sprintf("%s/config/", globalRoot)

	files := []string{}

	entries, err := os.ReadDir(localPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("error reading local directory: %s", err)
	}

	for _, e := range entries {
		files = append(files, fmt.Sprintf("config/%s", e.Name()))
	}

	entries, err = os.ReadDir(globalPath)
	if err != nil {
		return files, nil
	}

	for _, e := range entries {
		files = append(files, fmt.Sprintf("config/%s", e.Name()))
	}

	return files, nil
}
