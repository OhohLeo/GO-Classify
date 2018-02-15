package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type FileConfig struct {
	Icons IconsConfig `json:"icons"`
}

func (c *FileConfig) Update(rawMsg *json.RawMessage) error {
	return json.Unmarshal(*rawMsg, &c)
}

type File struct {
	Name        string            `json:"name"`
	Extension   string            `json:"extension"`
	ContentType string            `json:"contentType"`
	Icons       Icons             `json:"icons"`
	Infos       map[string]string `json:"infos"`

	Path         string      `json:"-"`
	AbsolutePath string      `json:"-"`
	FileInfo     os.FileInfo `json:"-"`

	file *os.File
}

func NewFileFromPath(path string, name string) (*File, error) {

	var f *os.File
	f, err := os.Open(path + "/" + name)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("NEW FILE at %s/%s\n", path, name)

	defer f.Close()

	return NewFileFromOsFile(f)
}

func NewFileFromOsFile(f *os.File) (file *File, err error) {

	// Get absolutePath & extension
	absolutePath := f.Name()
	extension := filepath.Ext(absolutePath)

	file = &File{
		file:         f,
		Name:         filepath.Base(absolutePath),
		Extension:    extension,
		Path:         filepath.Dir(absolutePath),
		AbsolutePath: absolutePath,
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

	if f.Icons == nil {
		f.Icons = NewIcons()
	}

	_, err = f.Icons.SetIcon(f.AbsolutePath, f.Name, &cfg.Icons)

	return
}

func (f *File) GetContentType() string {
	return f.ContentType
}

func (f *File) GetContents() (contents map[string]string) {

	contents = make(map[string]string)

	contents[f.Name] = f.AbsolutePath

	// If no icons to handle
	if f.Icons == nil {
		return
	}

	// Otherwise add icons
	for size, icon := range f.Icons {
		contents[size] = icon.GetAbsolutePath()
	}

	return
}

func (f *File) GetBlackList() (blacklist []string) {

	// If no icons to handle
	if f.Icons == nil {
		return
	}

	// Otherwise add icons
	for _, icon := range f.Icons {
		blacklist = append(blacklist, icon.Name)
	}

	return
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
