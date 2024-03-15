package fs

import (
	"errors"
	"fmt"
	"os"
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
