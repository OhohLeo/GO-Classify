package data

import (
	"encoding/json"
	"fmt"
	"github.com/ohohleo/classify/params"
	"github.com/quirkey/magick"
	"path/filepath"
)

type IconsConfig struct {
	Enable bool   `json:"enable"`
	Size   string `json:"size"`
}

func (i IconsConfig) GetParam(name string, data json.RawMessage) (result interface{}, err error) {

	switch name {
	case "path":
		result, err = params.GetPath(data)
	default:
		err = fmt.Errorf("import 'directory' invalid param '%s'", name)
	}

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

type Icons map[string]string

func (i Icons) NewConfig() Config {
	return new(IconsConfig)
}

func (i Icons) SetIcon(src string, name string, config *IconsConfig) (path string, err error) {

	if config.Enable == false {
		fmt.Printf("MAGICK DISABLE!\n")
		return
	}

	fmt.Printf("MAGICK %+v!\n", config)

	// Check if the icon doesn't already exist
	var ok bool
	if path, ok = i[config.Size]; ok {
		fmt.Printf("ALREADY MAGICK %s EXIST\n", src)
		return
	}

	fmt.Printf("MAGICK %s\n", src)
	icon, err := magick.NewFromFile(src)
	if err != nil {
		return
	}

	defer icon.Destroy()

	// Resize icon with specified size
	err = icon.Resize(config.Size)
	if err != nil {
		return
	}

	dst := filepath.Dir(src)

	fmt.Printf("MAGICK %s\n", dst)

	// Store file
	path = dst + "/" + "/" + name + "_" + config.Size + ".jpg"
	err = icon.ToFile(path)
	if err != nil {
		return
	}

	i[config.Size] = path
	return
}
