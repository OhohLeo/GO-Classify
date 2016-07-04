package collections

import (
	"github.com/ohohleo/classify/imports"
	"github.com/ohohleo/classify/websites"
)

type Item struct {
	status    int
	key       string
	inputs    map[string]imports.Data
	websites  map[string]websites.Data
	proposed  map[string]interface{}
	validated map[string]interface{}
}

func (i *Item) AddImportData(data imports.Data) error {
	return nil
}

func (i *Item) RemoveImportData(data imports.Data) error {
	return nil
}

func (i *Item) AddWebsiteData(data websites.Data) error {
	return nil
}

func (i *Item) RemoveWebsiteData(data websites.Data) error {
	return nil
}

func (i *Item) GetUniqKey() string {

	if i.key == "" {
		i.key = "UNQI"
	}

	return i.key
}
