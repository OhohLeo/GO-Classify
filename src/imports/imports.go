package imports

import (
	"os"
)

type Import interface {
	Start() (chan Data, error)
	Stop()
	GetType() string
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
