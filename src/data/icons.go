package data

import (
	"encoding/json"
	"fmt"
	"github.com/quirkey/magick"
	"log"
	"os"
	"path/filepath"
)

type IconsConfig struct {
	Enable  bool   `json:"enable"`
	Size    string `json:"size"`
	SubPath string `json:"subPath"`
}

func (i IconsConfig) GetParam(name string, data json.RawMessage) (result interface{}, err error) {
	err = fmt.Errorf("import 'icons' invalid param '%s'", name)
	return
}

func (i IconsConfig) Update(rawMsg *json.RawMessage) error {
	return json.Unmarshal(*rawMsg, &i)
}

func (c IconsConfig) Check() error {

	// Check size has correct format

	// Check path is valid

	return nil
}

type Icon struct {
	Name string `json:"name"`
	Path string `json:"-"`
}

func (i *Icon) GetAbsolutePath() string {
	return i.Path + "/" + i.Name
}

type Icons map[string]*Icon

func NewIcons() Icons {
	return make(map[string]*Icon)
}

func (i Icons) NewConfig() Config {
	return new(IconsConfig)
}

func (i Icons) SetIcon(src string, name string, config *IconsConfig) (path string, err error) {

	if config.Enable == false {
		return
	}

	// Check if the icon doesn't already referenced
	if _, ok := i[config.Size]; ok {
		return
	}

	dst := filepath.Dir(src)

	// Handle configuration sub-path
	if config.SubPath != "" {
		dst += "/" + config.SubPath
	}

	// Otherwise create icon
	icon := &Icon{
		Name: name + "_" + config.Size + ".jpg",
		Path: dst,
	}

	path = icon.GetAbsolutePath()

	// Check if the icon file doesn't already exist
	if _, err = os.Stat(icon.GetAbsolutePath()); err == nil {
		i[config.Size] = icon
		return
	}

	// Create destination path if it doesn't exist
	if _, err = os.Stat(dst); os.IsNotExist(err) {
		err = os.MkdirAll(dst, os.ModePerm)
		if err != nil {
			return
		}
	}

	iconMagick, err := magick.NewFromFile(src)
	if err != nil {
		return
	}

	defer iconMagick.Destroy()

	// Resize iconMagick with specified size
	err = iconMagick.Resize(config.Size)
	if err != nil {
		return
	}

	// Store file
	err = iconMagick.ToFile(path)
	if err != nil {
		return
	}

	log.Printf("Create icon at '%s'\n", icon.GetAbsolutePath())

	i[config.Size] = icon
	return
}
