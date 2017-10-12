package config

import (
	"reflect"
)

type Ref struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Comments string `json:"comments,omitempty"`
	Childs   []*Ref `json:"childs,omitempty"`
}

type Config struct {
	Refs []*Ref      `json:"refs"`
	Data interface{} `json:"data"`
}

func GetRef(data interface{}) *Config {

	val := reflect.ValueOf(data).Elem()

	config := &Config{
		Refs: getRefsByValue(val),
		Data: data,
	}

	return config
}

func getRefsByValue(val reflect.Value) []*Ref {

	refs := make([]*Ref, 0)

	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		kind := tag.Get("kind")
		if kind == "" {
			kind = typeField.Type.Kind().String()
		}

		ref := &Ref{
			Name:     tag.Get("json"),
			Comments: tag.Get("comments"),
			Type:     kind,
		}

		if kind == "struct" {
			ref.Childs = getRefsByValue(valueField)
		}

		refs = append(refs, ref)
	}

	return refs
}
