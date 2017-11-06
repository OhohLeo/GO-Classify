package data

import (
	"encoding/json"
	"fmt"
	"github.com/quirkey/magick"
)

type IconsConfig struct {
	Enable bool   `json:"enable"`
	Size   string `json:"size"`
	Path   string `json:"path" kind:"path"`
}

func (c *IconsConfig) UpdateConfig(json.RawMessage) error {
	return nil
}

func (c *IconsConfig) Check() error {

	// Check size has correct format

	// Check path is valid

	return nil
}

type CanIcons interface {
	HasIcon(size string) bool
	GetIcon(size string) (path string, ok bool)
	SetIcon(input Data, config *IconsConfig) (path string, err error)
}

type Icons struct {
	icons  map[string]string
	config *IconsConfig
}

func (i *Icons) InitIcons(config *IconsConfig) {
	i.icons = make(map[string]string)
	i.config = config
}

func (i *Icons) HasIcon(size string) bool {
	_, ok := i.icons[size]
	return ok
}

func (i *Icons) GetIcon(size string) (path string, ok bool) {
	path, ok = i.icons[size]
	return
}

func (i *Icons) SetIcon(input Data, config *IconsConfig) (path string, err error) {

	// Check if the icon doesn't already exist
	var ok bool
	if path, ok = i.GetIcon(config.Size); ok {
		return
	}

	// Check data compatibility
	var absolutePath string
	switch input.(type) {
	case *File:
		absolutePath = input.(*File).FullPath
	default:
		err = fmt.Errorf("Couldn't set icon from specified type")
		return
	}

	icon, err := magick.NewFromFile(absolutePath)
	if err != nil {
		return
	}

	defer icon.Destroy()

	// Resize icon with specified size
	err = icon.Resize(config.Size)
	if err != nil {
		return
	}

	dstPath := config.Path
	if dstPath == "" {
		dstPath = absolutePath
	}

	// Store file
	path = dstPath + "/" + "/" + input.GetName() + "_" + config.Size + ".jpg"
	err = icon.ToFile(path)
	if err != nil {
		return
	}

	i.icons[config.Size] = path
	return
}
