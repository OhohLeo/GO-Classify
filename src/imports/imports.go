package imports

import (
	"os"
)

type Import interface {
	Launch() (chan Data, error)
	GetType() string
}

type Data interface {
	GetType() string
	String() string
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
