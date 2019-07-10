package reference

import (
	"log"
	"reflect"
	"strings"
)

type Ref struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Comments string            `json:"comments,omitempty"`
	Childs   []*Ref            `json:"childs,omitempty"`
	Map      map[string][]*Ref `json:"map,omitempty"`
	Key      string            `json:"key,omitempty"`
	Value    interface{}       `json:"-"`
}

type Reference struct {
	Refs []*Ref      `json:"references"`
	Data interface{} `json:"data"`
}

func New(refs []*Ref, data interface{}) *Reference {

	return &Reference{
		Refs: refs,
		Data: data,
	}
}

func GetRefs(data interface{}) []*Ref {
	return getRefsByValue(reflect.ValueOf(data).Elem())
}

type Attribute struct {
	Name string `json:"-"`
	Type string `json:"type"`
}

func GetAttributes(data interface{}) map[string]Attribute {
	result := make(map[string]Attribute)

	for _, ref := range getRefsByValue(reflect.ValueOf(data).Elem()) {

		// Reject unused fields
		switch ref.Name {
		case "id":
			fallthrough
		case "ref":
			continue
		}

		// Reject unused types
		switch ref.Type {
		case "interface":
			fallthrough
		case "slice":
			fallthrough
		case "map":
			fallthrough
		case "struct":
			continue
		}

		result[ref.Name] = Attribute{
			Name: ref.Name,
			Type: ref.Type,
		}
	}

	return result
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

		name := tag.Get("json")

		// Get name until ','
		comaIdx := strings.Index(name, ",")
		if comaIdx >= 0 {
			name = name[0:comaIdx]
		}

		if name == "" || name == "-" {
			continue
		}

		kind := tag.Get("kind")
		if kind == "" {
			kind = typeField.Type.Kind().String()
		}

		switch kind {
		case "struct":
			if strings.HasPrefix(name, "date") {
				kind = "datetime"
			} else if name == "country" {
				kind = name
			}
		}

		ref := &Ref{
			Name:     name,
			Comments: tag.Get("comments"),
			Type:     kind,
			Value:    valueField.Interface(),
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

		case "interface":
		case "slice":
		case "stringlist":
		case "string":
		case "bool":
		case "int":
		case "uint64":
			// nothing to do

		case "datetime":
		case "country":
		case "string[]":
			// nothing to do

		default:
			log.Printf("Unhandled kind '%s'\n", kind)
		}

		refs = append(refs, ref)
	}

	return refs
}
