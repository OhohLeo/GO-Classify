package params

import (
	"encoding/json"
	"os"
)

type Path struct {
	Directory string `json:"directory"`
}

type PathResult struct {
	Current     string   `json:"current"`
	Directories []string `json:"directories"`
	Files       []string `json:"files"`
}

func GetPath(params json.RawMessage) (result interface{}, err error) {

	// Get path parameter
	var paramPath Path
	err = json.Unmarshal(params, &paramPath)
	if err != nil {
		return
	}

	if paramPath.Directory == "" {
		paramPath.Directory = "/"
	}

	// Check path validity
	var file *os.File
	file, err = os.Open(paramPath.Directory)
	if err != nil {
		return
	}

	// Search for files inside
	var fileInfos []os.FileInfo
	fileInfos, err = file.Readdir(0)
	if err != nil {
		return
	}

	pathResult := &PathResult{
		Current: paramPath.Directory,
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			pathResult.Directories =
				append(pathResult.Directories, fileInfo.Name())
		} else if fileInfo.Mode().IsRegular() {
			pathResult.Files =
				append(pathResult.Files, fileInfo.Name())
		}
	}

	result = pathResult
	return
}
