package imports

import (
	"os"
)

type Import interface {
	Check(map[string][]string, []string) error
	Start() (chan Data, error)
	Stop()
	GetType() string
	Eq(Import) bool
}

type Data interface {
	GetType() string
	String() string
	GetUniqKey() string
}

type File struct {
	Type     string            `json:"type"`
	Name     string            `json:"name"`
	Path     string            `json:"path"`
	FullPath string            `json:"fullpath"`
	Infos    map[string]string `json:"infos"`
	FileInfo os.FileInfo       `json:"-"`
}

func (f *File) GetType() string {
	return "file"
}

func (f *File) String() string {
	return f.Name
}

func (f *File) GetUniqKey() string {
	return f.GetType() + ":" + f.FullPath
}
