package data

import (
	"os"
)

type File struct {
	Type      string            `json:"type"`
	Name      string            `json:"name"`
	FullName  string            `json:"fullname"`
	Extension string            `json:"extension"`
	Path      string            `json:"path"`
	FullPath  string            `json:"fullpath"`
	Infos     map[string]string `json:"infos"`
	FileInfo  os.FileInfo       `json:"-"`
}

func (f *File) GetRef() Ref {
	return FILE
}

func (f *File) GetName() string {
	return f.Name
}
