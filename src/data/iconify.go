package data

import (
	"fmt"
	"github.com/quirkey/magick"
)

type IconifyConfig struct {
	Enable bool   `json:"enable"`
	Size   string `json:"size"`
	Path   string `json:"path" kind:"path"`
}

func (cfg *IconifyConfig) Check() error {

	// Check size has correct format

	// Check path is valid

	return nil
}

type CanIconify interface {
	HasIcon(size string) bool
	GetIcon(size string) (path string, ok bool)
	SetIcon(input Data, config *IconifyConfig) (path string, err error)
}

type Iconify struct {
	icons map[string]string
}

func (i *Iconify) HasIcon(size string) bool {
	_, ok := i.icons[size]
	return ok
}

func (i *Iconify) GetIcon(size string) (path string, ok bool) {
	path, ok = i.icons[size]
	return
}

func (i *Iconify) SetIcon(input Data, config *IconifyConfig) (path string, err error) {

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

	// Store file
	path = config.Path + "/" + "/" + input.GetName() + "_" + config.Size + ".jpg"
	err = icon.ToFile(path)
	if err != nil {
		return
	}

	i.icons[config.Size] = path
	return
}
