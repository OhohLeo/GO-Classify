package config

import (
	"reflect"
)

type Ref struct {
	Name     string `json:"name"`
	Comments string `json:"comments"`
	Type     string `json:"type"`

	Childs []*Ref `json:"childs"`
}

type Config struct {
	Refs []*Ref      `json:"refs"`
	Data interface{} `json:"data"`
}

func GetRefs(data interface{}) *Config {

	config := &Config{
		Refs: make([]*Ref, 0),
		Data: data,
	}

	val := reflect.ValueOf(data).Elem()

	for i := 0; i < val.NumField(); i++ {

		typeField := val.Type().Field(i)
		tag := typeField.Tag

		kind := tag.Get("kind")
		if kind == "" {
			kind = typeField.Type.Kind().String()
		}

		config.Refs = append(config.Refs, &Ref{
			Name:     tag.Get("json"),
			Comments: tag.Get("comments"),
			Type:     kind,
		})
	}

	return config
}
