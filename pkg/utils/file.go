package utils

import (
	"errors"
	"io/fs"
	"path/filepath"
)

// GetConfigFilepathFromFilename searches the config directory for the specified file, then the working directory, and then "/etc/config"
func GetConfigFilepathFromFilename(filename string) string {
	// lets search the config directory first
	path, err := getFilepathFromFilename("./etc/config", filename)
	// if there is NOT an error, simply return the path
	if err == nil {
		return path
	}

	// lets try the working directory
	path, err = getFilepathFromFilename(".", filename)
	if err == nil {
		return path
	}

	path, _ = getFilepathFromFilename("/etc/config", filename)
	return path
}

func getFilepathFromFilename(root, filename string) (string, error) {
	var pathToFile string
	var found bool

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info.Name() == filename {
			pathToFile = path
			found = true
			return nil
		}

		return err
	})

	if err != nil {
		return "", err
	}

	if !found {
		return "", errors.New("file not found")
	}

	return pathToFile, nil
}
