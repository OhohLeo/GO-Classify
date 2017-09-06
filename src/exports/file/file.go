package file

import (
	"encoding/json"
	"github.com/ohohleo/classify/data"
	"github.com/ohohleo/classify/exports"
)

type File struct {
}

func ToBuild() exports.BuildExport {
	return exports.BuildExport{
		Create: Create,
	}
}

func Create(input json.RawMessage,
	config map[string][]string,
	collections []string) (e exports.Export, err error) {

	var file File
	err = json.Unmarshal(input, &file)
	if err != nil {
		e = &file
	}
	return
}

func (f *File) GetRef() exports.Ref {
	return exports.FILE
}

func (f *File) OnInput(data data.Data) error {
	return nil
}

func (f *File) Stop() error {
	return nil
}

func (f *File) Eq(new exports.Export) bool {
	// newFile, _ := new.(*File)
	return true
}
