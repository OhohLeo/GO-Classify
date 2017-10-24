package data

import (
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	*Iconify
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	FullName    string            `json:"fullname"`
	Extension   string            `json:"extension"`
	Path        string            `json:"path"`
	FullPath    string            `json:"fullpath"`
	ContentType string            `json:"contentType"`
	Infos       map[string]string `json:"infos"`
	FileInfo    os.FileInfo       `json:"-"`
	file        *os.File          `json:"-"`
}

func NewFileFromPath(path string, name string) (*File, error) {

	var f *os.File
	f, err := os.Open(path + "/" + name)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return NewFileFromOsFile(path, f)
}

func NewFileFromOsFile(path string, f *os.File) (file *File, err error) {

	// Get fullpath & extension
	fullname := f.Name()
	fullpath := path + "/" + fullname
	extension := filepath.Ext(fullpath)

	file = &File{
		Name:      strings.TrimRight(fullname, extension),
		FullName:  fullname,
		Extension: extension,
		Path:      path,
		FullPath:  fullpath,
		file:      f,
	}
	return
}

func (f *File) GetRef() Ref {
	return FILE
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) GetContentType() string {
	return f.ContentType
}

func (f *File) GetFileInfo() (fileInfo os.FileInfo, err error) {

	if f.file == nil {
		f.file, err = os.Open(f.FullPath)
		if err != nil {
			return
		}

		defer f.file.Close()
	}

	// Get file infos
	fileInfo, err = f.file.Stat()
	if err == nil {
		f.FileInfo = fileInfo
	}

	return
}
