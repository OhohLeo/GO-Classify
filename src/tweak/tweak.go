package tweak

import (
	"encoding/json"
	"fmt"
	"regexp"
	//"sort"

	"github.com/ohohleo/classify/data"
)

type Value struct {
	Value  string         `json:"value,omitempty"`
	Regexp *regexp.Regexp `json:"-"`
	Format string         `json:"format,omitempty"`
}

func (v *Value) MarshalJSON() (src []byte, err error) {

	data := struct {
		Value  string `json:"value,omitempty"`
		Regexp string `json:"regexp,omitempty"`
		Format string `json:"format,omitempty"`
	}{
		Value:  v.Value,
		Format: v.Format,
	}

	if v.Regexp != nil {
		data.Regexp = v.Regexp.String()
	}

	src, err = json.Marshal(data)
	return
}

func (v *Value) UnmarshalJSON(src []byte) (err error) {

	// Decode value
	type value *Value
	if err = json.Unmarshal(src, value(v)); err != nil {

		// Check is it text only?
		if src[0] == '"' && src[len(src)-1] == '"' {
			v.Value = string(src)
			err = nil
			return
		}

		return
	}

	// Try to decode regex
	data := struct {
		Regexp string `json:"regexp"`
	}{}

	if err = json.Unmarshal(src, &data); err != nil {
		return
	}

	if data.Regexp == "" {
		return
	}

	// Set regex
	regex, err := regexp.Compile(data.Regexp)
	if err != nil {
		return
	}

	v.Regexp = regex
	return
}

// Fields is relation between field and value/regex
type Fields map[string]*Value

func (f Fields) Check(dst bool) (err error) {

	for _, value := range f {

		// Reject destination value with regexp
		if dst && value.Regexp != nil {
			err = fmt.Errorf("destination value can't handle regexp")
			return
		}

		// Check field presence

		// Check value

		// Check format

	}

	return
}

// func (f Fields) GetValue(key string, ref data.Ref) (value string, err error) {

// 	var ok bool
// 	value, ok = f[key]
// 	if ok == false {
// 		err = fmt.Errorf("invalid field '%s' of data ref '%s'", key, ref.String())
// 		return
// 	}

// 	return
// }

type Tweak struct {

	// Ref => Field => Regex/Value
	Source map[string]Fields `json:"source"`

	// [ ComputedData/Item/Export ] => Field => Value
	Destination map[string]Fields `json:"destination"`
}

func New(src []byte) (*Tweak, error) {

	var t Tweak
	err := json.Unmarshal(src, &t)
	return &t, err
}

// Check keys compatibility
func (t *Tweak) Check(d data.Data) error {

	ref := d.GetRef()

	var listFields []Fields

	// Check data presence
	if f, ok := t.Source[ref.String()]; ok {
		listFields = append(listFields, f)
	}

	if f, ok := t.Destination[ref.String()]; ok {
		listFields = append(listFields, f)
	}

	// No field matching data : nothing to do
	if len(listFields) == 0 {
		return nil
	}

	// Get data reference content

	for idx, fields := range listFields {

		// Check each fields
		if err := fields.Check(idx > 0); err != nil {
			return err
		}
	}

	return nil
}

type Results map[string]string

func (t *Tweak) Tweak(src map[string]data.Data) (results map[string]Results, err error) {

	// Get values based on regexp
	// var values [][]string

	// // Sort keys by name
	// keys := make([]string, len(t.Source))

	// i := 0
	// for key, _ := range t.Source {
	// 	keys[i] = key
	// 	i++
	// }

	// sort.Strings(keys)

	// // Based on values : determined results for each field
	// for name, fields := range t.Destination {

	// }

	return
}
