package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileConfig struct {
	Icons IconsConfig `json:"icons"`
}

func (c *FileConfig) Update(rawMsg *json.RawMessage) error {
	return json.Unmarshal(*rawMsg, &c)
}

type File struct {
	Type         string            `json:"type"`
	Name         string            `json:"name"`
	BaseName     string            `json:"basename"`
	Extension    string            `json:"extension"`
	Path         string            `json:"path"`
	AbsolutePath string            `json:"absolutePath"`
	ContentType  string            `json:"contentType"`
	Icons        Icons             `json:"icons,omitempty"`
	Infos        map[string]string `json:"infos"`
	FileInfo     os.FileInfo       `json:"-"`

	file *os.File
}

func NewFileFromPath(path string, name string) (*File, error) {

	var f *os.File
	f, err := os.Open(path + "/" + name)
	if err != nil {
		return nil, err
	}

	fmt.Printf("NEW FILE FROM PATH %s %s %s\n", path, name, f.Name())

	defer f.Close()

	return NewFileFromOsFile(f)
}

func NewFileFromOsFile(f *os.File) (file *File, err error) {

	// Get absolutePath & extension
	absolutePath := f.Name()
	basename := filepath.Base(absolutePath)
	extension := filepath.Ext(absolutePath)

	file = &File{
		file:         f,
		Name:         strings.TrimRight(basename, extension),
		BaseName:     basename,
		Extension:    extension,
		Path:         filepath.Dir(absolutePath),
		AbsolutePath: absolutePath,
		Icons:        make(map[string]string),
	}

	return
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) GetRef() Ref {
	return FILE
}

func (f *File) NewConfig() Config {
	return new(FileConfig)
}

func (f *File) OnCollection(config Config) (err error) {

	cfg, ok := config.(*FileConfig)
	if ok == false {
		err = fmt.Errorf("expected file configuration %+v", config)
		return
	}

	_, err = f.Icons.SetIcon(f.AbsolutePath, f.Name, &cfg.Icons)
	fmt.Printf("File OnCollection %s %+v\n", f.AbsolutePath, err)
	return
}

func (f *File) GetContentType() string {
	return f.ContentType
}

func (f *File) GetFileInfo() (fileInfo os.FileInfo, err error) {

	if f.FileInfo != nil {
		fileInfo = f.FileInfo
		return
	}

	if f.file == nil {
		f.file, err = os.Open(f.AbsolutePath)
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
