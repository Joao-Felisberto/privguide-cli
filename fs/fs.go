// Package to abstract file system accesses,
// namely by handling lookup from both the global and local directories
//
// By default, the local path is `.devprivops/` and the global path is `/etc/devprivops/`.
// Files in the local path override those in the global path
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

var LocalDir = fmt.Sprintf("./.%s", util.AppName)
var GlobalDir = fmt.Sprintf("/etc/%s", util.AppName)

// Returns the full path of a file using the pre-determined paths to the local and global directories
//
// `relativePath`: the path relative to either root
//
// returns: the full path to the provided file
func GetFile(relativePath string) (string, error) {
	return getFile(
		relativePath,
		LocalDir,
		GlobalDir,
	)
}

// Returns the full path of a file using the provided paths to the local and global directories
//
// `localRoot`: the root of the local directory
//
// `globalRoot`: the root of the global directory
//
// `relativePath` the path relative to either root
//
// returns: the full path to the provided file
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

// Returns the relative paths of the system descriptions under `descriptions/` using the default paths to the local and global directories
//
// `relativePath` the path relative to either root
//
// returns: the relative paths of the system descriptions
func GetDescriptions(descriptionRoot string) ([]string, error) {
	return getDescriptions(
		descriptionRoot,
		LocalDir,
		GlobalDir,
	)
}

// Returns the relative paths of the system descriptions under `descriptions/` using the provided paths to the local and global directories
//
// `localRoot`: the root of the local directory
//
// `globalRoot`: the root of the global directory
//
// `relativePath` the path relative to either root
//
// returns: the relative paths of the system descriptions
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
// returns: the directory names of the system regulation directories
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
// returns: the directory names of the system regulation directories
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
