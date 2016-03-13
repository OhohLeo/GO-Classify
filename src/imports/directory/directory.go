package directory

import (
	"github.com/ohohleo/classify/imports"
	"io/ioutil"
	"os"
)

type Directory struct {
	Path        string `json:"path"`
	IsRecursive bool   `json:"is_recursive"`
}

// Return a channel of files in the directory
func (r *Directory) Launch() (chan imports.File, error) {

	// Check we have an existing directory
	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return nil, err
	}

	c := make(chan imports.File)

	go func() {
		readDirectory(c, r.Path, r.IsRecursive)
		close(c)
	}()

	return c, nil
}

func readDirectory(c chan imports.File, path string, isRecursive bool) {

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
		c <- imports.File{
			Path:     path,
			FullPath: fullpath,
			FileInfo: f,
		}
	}
}
