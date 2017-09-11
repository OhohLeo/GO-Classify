package file

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ohohleo/classify/data"
	"github.com/ohohleo/classify/exports"
	"io/ioutil"
	"os"
	"strconv"
)

type File struct {
	Path        string `json:"path"`
	Permissions string `json:"permissions"`
	mode        os.FileMode
}

func ToBuild() exports.Build {
	return exports.Build{
		Create: Create,
	}
}

func Create(input json.RawMessage,
	config json.RawMessage,
	collections []string) (e exports.Export, err error) {

	var file File
	err = json.Unmarshal(input, &file)
	if err != nil {
		return
	}

	// Validate path
	var stat os.FileInfo
	if stat, err = os.Stat(file.Path); err != nil || stat.IsDir() == false {
		err = fmt.Errorf("'%s' should be a valid directory path ", file.Path)
		return
	}

	// Validate permission
	var mode uint64
	if mode, err = strconv.ParseUint(file.Permissions, 8, 32); err != nil {
		return
	}

	file.mode = os.FileMode(mode)

	e = &file
	return
}

func (f *File) GetRef() exports.Ref {
	return exports.FILE
}

func (f *File) OnInput(input data.Data) error {

	switch input.GetRef() {

	case data.ATTACHMENT:
		return f.onAttachment(input)

	default:
		return fmt.Errorf("Unhandled reference '%s'", input.GetRef().String())
	}

	return nil
}

func (f *File) Stop() error {
	return nil
}

func (f *File) Eq(new exports.Export) bool {
	// newFile, _ := new.(*File)
	return true
}

func (f *File) onAttachment(input data.Data) error {

	attachment, ok := input.(*data.Attachment)
	if ok == false {
		return fmt.Errorf("Invalid attachment data")
	}

	reader := bytes.NewReader(attachment.Data)
	data := base64.NewDecoder(base64.StdEncoding, reader)

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(data)

	err := ioutil.WriteFile(f.Path+"/"+attachment.Name, buffer.Bytes(), f.mode)
	if err != nil {
		return err
	}

	fmt.Printf("ATTACHMENT %s\n", attachment.Name)
	return nil
}
