package config

import (
	"log"
	"reflect"
)

type Ref struct {
	Name     string            `json:"name"`
	Key      string            `json:"key"`
	Type     string            `json:"type"`
	Comments string            `json:"comments,omitempty"`
	Childs   []*Ref            `json:"childs,omitempty"`
	Map      map[string][]*Ref `json:"map,omitempty"`
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

	if val.Kind() == reflect.Interface && !val.IsNil() {

		elm := val.Elem()

		if elm.Kind() == reflect.Ptr &&
			!elm.IsNil() &&
			elm.Elem().Kind() == reflect.Ptr {
			val = elm
		}
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

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

		switch kind {
		case "struct":
			ref.Childs = getRefsByValue(valueField)
		case "map":

			if ref.Map == nil {
				ref.Map = make(map[string][]*Ref)
			}

			for _, k := range valueField.MapKeys() {

				mapKey, ok := k.Interface().(string)
				if ok == false {
					log.Printf("Unhandled map key '%+v'\n", k.Interface())
					return nil
				}

				mapValue := valueField.MapIndex(k).Interface()
				ref.Map[mapKey] = getRefsByValue(reflect.ValueOf(mapValue))
			}

		case "bool":
		case "int":
		case "stringlist":
			// nothing to do
		default:
			log.Printf("Unhandled kind '%s'\n", kind)
		}

		refs = append(refs, ref)
	}

	return refs
}
