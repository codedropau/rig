package buildctx

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Package files into a build context for Docker.
func Package(paths ...string) (*os.File, error) {
	tmp, err := ioutil.TempFile(os.TempDir(), "buildctx.*.tar")
	if err != nil {
		return nil, err
	}

	tarWriter := tar.NewWriter(tmp)
	defer tarWriter.Close()

	var filePaths []string

	for _, path := range paths {
		tree, err := BuildTree(path)
		if err != nil {
			return nil, err
		}

		filePaths = append(filePaths, tree...)
	}

	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		stat, err := file.Stat()
		if err != nil {
			return nil, err
		}

		var target string

		if stat.Mode()&os.ModeSymlink != 0 {
			target, err = os.Readlink(stat.Name())
			if err != nil {
				return nil, err
			}
		}

		header, err := tar.FileInfoHeader(stat, filepath.ToSlash(target))
		if err != nil {
			return nil, err
		}

		// We override this value because we want the full path in our tar.
		//   eg.
		//     Want = /path/to/file
		//     Dont = /file
		header.Name = filePath

		err = tarWriter.WriteHeader(header)
		if err != nil {
			return nil, err
		}

		if stat.IsDir() {
			continue
		}

		_, err = io.Copy(tarWriter, file)
		if err != nil {
			return nil, err
		}

		err = file.Close()
		if err != nil {
			return nil, err
		}
	}

	// Docker requires that the file be closed and reopened.
	// If this is not done then Docker cannot find the files loaded into the context.
	err = tmp.Close()
	if err != nil {
		return nil, err
	}

	return os.Open(tmp.Name())
}

func BuildTree(path string) ([]string, error) {
	var list []string

	stat, err := os.Stat(path)
	if err != nil {
		return list, err
	}

	if !stat.Mode().IsDir() {
		return []string{path}, nil
	}

	dir, err := os.Open(path)
	if err != nil {
		return list, err
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		return list, err
	}

	for _, file := range files {
		p := filepath.Join(dir.Name(), file.Name())

		list = append(list, p)

		if file.IsDir() {
			subList, err := BuildTree(p)
			if err != nil {
				return list, err
			}

			list = append(list, subList...)

			continue
		}
	}

	return list, nil
}