package compose

import (
	"os"
	"path/filepath"
)

func Project(override string) (string, error) {
	if override != "" {
		return override, nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Base(dir), nil
}