package directory

import (
	"github.com/ohohleo/classify/imports"
	"io/ioutil"
	//"log"
	"errors"
	"os"
)

type Directory struct {
	Path        string `json:"path"`
	IsRecursive bool   `json:"is_recursive"`
	needToStop  bool
}

func (r *Directory) Check(config map[string][]string, collections []string) error {

	// Check we have an existing directory
	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return err
	}

	// Check that the directory is in the global directories
	globalPaths, ok := config["*"]
	if ok {
		for _, path := range globalPaths {
			if r.Path == path {
				return nil
			}
		}
	}

	// Check that the directory is authorised for all specified collections
	for _, name := range collections {

		paths, ok := config[name]
		if ok == false {
			continue
		}

		for _, path := range paths {
			if r.Path == path {
				return nil
			}
		}
	}

	return errors.New("invalid or unauthorised import path '" + r.Path + "'")
}

// Return a channel of files in the directory
func (r *Directory) Start() (chan imports.Data, error) {

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

func (r *Directory) Eq(new imports.Import) bool {
	newDirectory, _ := new.(*Directory)
	return r.Path == newDirectory.Path &&
		r.IsRecursive == newDirectory.IsRecursive
}
