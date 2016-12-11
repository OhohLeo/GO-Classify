package directory

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ohohleo/classify/imports"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Directory struct {
	Path        string `json:"path"`
	IsRecursive bool   `json:"is_recursive"`
	exiftoolCmd string
	needToStop  bool
}

func (r *Directory) Check(config map[string][]string, collections []string) error {

	// Check we have an existing directory
	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return err
	}

	// Check if exiftool exists
	cmd, err := exec.LookPath("exiftool")
	if err != nil {
		return err
	}

	// Store exiftool command
	r.exiftoolCmd = cmd

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

		file := imports.File{
			Path:     path,
			FullPath: fullpath,
			FileInfo: f,
		}

		// Search for file header infos
		if r.exiftoolCmd != "" {
			r.Analyse(r.exiftoolCmd, file)
		}

		// Send file info through channel
		c <- file
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

func (r *Directory) Analyse(cmdStr string, file imports.File) {

	fullpath := file.FullPath

	// Prepare command
	cmd := exec.Command(cmdStr, fullpath)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating StdoutPipe for '%s' [%s]: %s\n",
			cmdStr, fullpath, err.Error())
		return
	}

	// Analyse response
	scanner := bufio.NewScanner(cmdReader)
	go func() {

		if file.Infos == nil {
			file.Infos = make(map[string]string)
		}

		for scanner.Scan() {
			// Get result line by line
			res := strings.SplitN(scanner.Text(), ":", 2)
			key := strings.TrimSpace(res[0])
			value := strings.TrimSpace(res[1])

			// Store infos
			file.Infos[key] = value
		}
	}()

	// Execute the command
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error starting '%s' [%s]: %s\n", cmdStr, fullpath, err.Error())
		return
	}

	// Wait for the answer
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error waiting '%s' [%s]: %s\n", cmdStr, fullpath, err.Error())
		return
	}

	// TODO : Check when file not recognized
}
