package directory

import (
	"github.com/ohohleo/classify/imports"
	"io/ioutil"
	//"log"
	"os"
)

type Directory struct {
	Path        string `json:"path"`
	IsRecursive bool   `json:"is_recursive"`
	needToStop  bool
}

// Return a channel of files in the directory
func (r *Directory) Start() (chan imports.Data, error) {

	// Check we have an existing directory
	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return nil, err
	}

	c := make(chan imports.Data)

	go func() {
		r.readDirectory(c, r.Path, r.IsRecursive)
		close(c)
	}()

	return c, nil
}

func (r *Directory) Stop() {
	r.needToStop = true
}

func (r *Directory) readDirectory(c chan imports.Data, path string, isRecursive bool) {

	// Read directory
	files, _ := ioutil.ReadDir(path)

	for _, f := range files {

		if r.needToStop {
			break
		}

		fullpath := path + "/" + f.Name()

		if f.IsDir() {

			// Read recursively
			if isRecursive {
				r.readDirectory(c, fullpath, isRecursive)
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

func (r *Directory) GetType() string {
	return "directory"
}

func (r *Directory) GetUniqKey() string {
	return ""
}
