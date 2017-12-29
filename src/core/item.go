package core

import (
	//	"encoding/json"
	"fmt"
	"github.com/ohohleo/classify/data"
	"os"
)

type Item struct {
	Id       Id        `json:"id,string"`
	Ref      string    `json:"ref"`
	Engine   data.Data `json:"data"`
	Contents []string  `json:"contents"`
	contents map[string]string
}

func (i *Item) SetData(input data.Data) {
	i.Ref = input.GetRef().String()
	i.Engine = input
}

func (i *Item) LinkToData(d data.Data) error {

	// Sync data content list to item
	if dataContents, ok := d.(data.HasContents); ok {

		for name, src := range dataContents.GetContents() {

			if err := i.SetContent(name, src); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *Item) SetContent(name string, src string) error {

	if i.contents == nil {
		i.contents = make(map[string]string)
	}

	// Check if name not already used
	if _, ok := i.contents[name]; ok {
		return fmt.Errorf("Content '%s' already existing in item '%d'",
			name, i.Id)
	}

	// Check if content file does exist
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("Try to add content '%s' "+
			"with not existing file '%s' in item '%d'",
			name, src, i.Id)
	}

	i.contents[name] = src
	i.Contents = append(i.Contents, name)
	return nil
}

func (i *Item) GetContent(name string) string {

	src, ok := i.contents[name]
	if ok {
		return src
	}

	return ""
}
