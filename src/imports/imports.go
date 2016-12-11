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
	Path     string
	FullPath string
	FileInfo os.FileInfo
	Infos    map[string]string
}

func (f File) GetType() string {
	return "file"
}

func (f File) String() string {
	return f.Path
}

func (f File) GetUniqKey() string {
	return f.GetType() + ":" + f.FullPath
}
