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
	return getFile(
		relativePath,
		fmt.Sprintf("./.%s", appName),
		fmt.Sprintf("/etc/%s", appName),
	)
}

func getFile(relativePath string, localRoot string, globalRoot string) (string, error) {
	localPath := fmt.Sprintf("%s/%s", localRoot, relativePath)
	if _, err := os.Stat(localPath); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		defaultPath := fmt.Sprintf("%s/%s", globalRoot, relativePath)
		if _, err := os.Stat(defaultPath); errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		return defaultPath, nil
	}
	return localPath, nil
}

func GetDescriptions(descriptionRoot string) ([]string, error) {
	return getDescriptions(
		descriptionRoot,
		fmt.Sprintf("./.%s", appName),
		fmt.Sprintf("/etc/%s", appName),
	)
}

func getDescriptions(descriptionRoot string, localRoot string, globalRoot string) ([]string, error) {
	localPath := fmt.Sprintf("%s/%s/", localRoot, descriptionRoot)
	globalPath := fmt.Sprintf("%s/%s/", globalRoot, descriptionRoot)

	files := []string{}

	entries, err := os.ReadDir(localPath)
	if err != nil {
		return nil, err
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

func GetRegulations() ([]string, error) {
	return getRegulations(
		fmt.Sprintf("./.%s", appName),
		fmt.Sprintf("/etc/%s", appName),
	)
}

func getRegulations(localRoot string, globalRoot string) ([]string, error) {
	localPath := fmt.Sprintf("%s/regulations/", localRoot)
	defaultPath := fmt.Sprintf("%s/regulations/", globalRoot)

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
