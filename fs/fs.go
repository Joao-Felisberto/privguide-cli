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

const appName = "devprivops"

func GetFile(relativePath string) (string, error) {
	localPath := fmt.Sprintf("./.%s/%s", appName, relativePath)
	if _, err := os.Stat(localPath); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		defaultPath := fmt.Sprintf("/etc/%s/%s", appName, relativePath)
		if _, err := os.Stat(defaultPath); errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		return defaultPath, nil
	}
	return localPath, nil
}

func GetDescriptions() ([]string, error) {
	localPath := fmt.Sprintf("./.%s/descriptions/", appName)
	defaultPath := fmt.Sprintf("/etc/%s/descriptions/", appName)

	files := []string{}

	entries, err := os.ReadDir(localPath)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		// fpath := fmt.Sprintf("%s%s", localPath, e.Name())
		files = append(files, fmt.Sprintf("descriptions/%s", e.Name()))
	}

	entries, err = os.ReadDir(defaultPath)
	if err != nil {
		return files, nil
	}

	for _, e := range entries {
		// fpath := fmt.Sprintf("%s%s", defaultPath, e.Name())
		files = append(files, fmt.Sprintf("descriptions/%s", e.Name()))
	}

	return files, nil
}

func GetRegulations() ([]string, error) {
	localPath := fmt.Sprintf("./.%s/regulations/", appName)
	defaultPath := fmt.Sprintf("/etc/%s/regulations/", appName)

	files := []string{}

	localRegulations, err := getDirsInDir(localPath)
	if err != nil {
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
