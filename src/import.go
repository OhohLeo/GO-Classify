package main

import (
	"io/ioutil"
	"os"
)

type File struct {
	Path     string
	FullPath string
	FileInfo os.FileInfo
}

// Return a channel of files in the directory
func ReadDirectory(path string, isRecursive bool) (chan File, error) {

	// Check we have an existing directory
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	c := make(chan File)

	go func() {
		readDirectory(c, path, isRecursive)
		close(c)
	}()

	return c, nil
}

func readDirectory(c chan File, path string, isRecursive bool) {

	// Read directory
	files, _ := ioutil.ReadDir(path)

	for _, f := range files {

		fullpath := path + "/" + f.Name()

		if f.IsDir() {

			// Read recursively
			if isRecursive {
				readDirectory(c, fullpath, isRecursive)
			}

			continue
		}

		// Send file info through channel
		c <- File{
			Path:     path,
			FullPath: fullpath,
			FileInfo: f,
		}
	}
}
