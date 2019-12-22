package core

import (
	"github.com/ohohleo/classify/reference"
)

type DatasReference map[string]DataReference

type DataReference struct {
	Name       string                         `json:"-"`
	Attributes map[string]reference.Attribute `json:"attributes,omitempty"`
}

func GetDatasReference(datas map[string]interface{}) DatasReference {
	result := make(map[string]DataReference)
	for dataName, data := range datas {
		result[dataName] = DataReference{
			Name:       dataName,
			Attributes: reference.GetAttributes(data),
		}
	}
	return result
}

type References struct {
	Datas   DatasReference            `json:"datas"`
	Imports map[string]DatasReference `json:"imports,omitempty"`
	Exports map[string]DatasReference `json:"exports,omitempty"`
}
