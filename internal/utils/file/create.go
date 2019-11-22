package file

import (
	"fmt"
	"os"
	"path/filepath"
)

func Write(path, data string) error {
	directory := filepath.Dir(path)

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprint(f, data)
	if err != nil {
		return err
	}

	return nil
}
