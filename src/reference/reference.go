package reference

import (
	"log"
	"reflect"

	"github.com/ohohleo/classify/params"
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

type GetRefOptions struct {
	Src    string
	Params map[string]params.HasParam
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

func GetRefs(data interface{}) []*Ref {
	return getRefsByValue(reflect.ValueOf(data).Elem(), nil)
}

func GetParamRefs(prefix string, data interface{}) ([]*Ref, map[string]params.HasParam) {

	params := make(map[string]params.HasParam)
	return getRefsByValue(reflect.ValueOf(data).Elem(), &GetRefOptions{prefix, params}), params
}

func getRefsByValue(val reflect.Value, opt *GetRefOptions) []*Ref {

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
	if opt != nil {

		// Handle HasParam interface
		if opt.Params != nil {
			if param, ok := val.Interface().(params.HasParam); ok {
				opt.Params[opt.Src] = param
			}
		}
	}

	for i := 0; i < val.NumField(); i++ {

		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		name := tag.Get("json")
		if name == "" || name == "-" {
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
			Value:    valueField.Interface(),
		}

		switch kind {

		case "struct":
			if opt != nil {
				opt.Src = opt.Src + "-" + name
			}

			ref.Childs = getRefsByValue(valueField, opt)

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
				if opt != nil {
					opt.Src = opt.Src + "-" + mapKey
				}

				ref.Map[mapKey] = getRefsByValue(
					reflect.ValueOf(mapValue), opt)
			}

		case "interface":
		case "slice":
		case "stringlist":
		case "string":
		case "bool":
		case "int":
		case "uint64":
			// nothing to do

		default:
			log.Printf("Unhandled kind '%s'\n", kind)
		}

		refs = append(refs, ref)
	}

	return refs
}
