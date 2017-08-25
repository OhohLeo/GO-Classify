package imports

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

func (f *File) GetRef() string {
	return "file"
}

func (f *File) String() string {
	return f.Name
}

func (f *File) GetUniqKey() string {
	return f.GetRef() + ":" + f.FullPath
}
