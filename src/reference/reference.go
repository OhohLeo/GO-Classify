package reference

import (
	"github.com/ohohleo/classify/params"
	"log"
	"reflect"
)

type Ref struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Comments string            `json:"comments,omitempty"`
	Childs   []*Ref            `json:"childs,omitempty"`
	Map      map[string][]*Ref `json:"map,omitempty"`
	Key      string            `json:"key,omitempty"`
}

type Reference struct {
	Refs []*Ref      `json:"refs"`
	Data interface{} `json:"data"`
}

func New(refs []*Ref, data interface{}) *Reference {

	return &Reference{
		Refs: refs,
		Data: data,
	}
}

func GetRefs(prefix string, data interface{}) ([]*Ref, map[string]params.HasParam) {

	params := make(map[string]params.HasParam)
	return getRefsByValue(prefix, reflect.ValueOf(data).Elem(), params), params
}

func getRefsByValue(src string, val reflect.Value, p map[string]params.HasParam) []*Ref {

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

	// Store structure handling params
	if param, ok := val.Interface().(params.HasParam); ok {
		p[src] = param
	}

	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		name := tag.Get("json")
		if name == "" {
			continue
		}

		kind := tag.Get("kind")
		if kind == "" {
			kind = typeField.Type.Kind().String()
		}

		ref := &Ref{
			Name:     name,
			Comments: tag.Get("comments"),
			Type:     kind,
		}

		switch kind {
		case "struct":
			ref.Childs = getRefsByValue(src+"-"+name, valueField, p)
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
				ref.Map[mapKey] = getRefsByValue(
					src+"-"+mapKey, reflect.ValueOf(mapValue), p)
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
