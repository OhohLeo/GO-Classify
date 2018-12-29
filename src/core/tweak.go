package core

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ohohleo/classify/reference"
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

func (f Fields) Check(data interface{}, dst bool) (err error) {

	fieldCompatible := make(map[string]struct{})

	for _, ref := range reference.GetRefs(data) {

		// Handle string only
		if ref.Type != "string" {
			continue
		}

		fieldCompatible[ref.Name] = struct{}{}
	}

	for key, value := range f {

		// Reject destination value with regexp
		if dst && value.Regexp != nil {
			err = fmt.Errorf("destination '%s' value can't handle regexp", key)
			return
		}

		// Check field presence
		if _, ok := fieldCompatible[key]; ok == false {
			err = fmt.Errorf("field not found '%s'", key)
			return
		}
	}

	return
}

func (f Fields) GetRawDatas(data interface{}) (results map[string][]string, err error) {

	results = make(map[string][]string)

	for _, ref := range reference.GetRefs(data) {

		// TODO: Handle time.Time

		// Handle string only
		if ref.Type != "string" {
			continue
		}

		// Search matching field name
		v, ok := f[ref.Name]
		if ok == false {
			continue
		}

		// Convert into string
		raw, ok := ref.Value.(string)
		if ok == false {
			err = fmt.Errorf("source expected 'string' on data field '%s'", ref.Name)
			return
		}

		// Check if regexp is handled
		if v.Regexp != nil {

			// Invalid regexp : nothing stored
			if v.Regexp.MatchString(raw) == false {
				log.Printf("source regexp '%s' not matching raw '%s'", v.Regexp.String(), raw)
				results[ref.Name] = make([]string, v.Regexp.NumSubexp())
				continue
			}

			if nb := v.Regexp.NumSubexp(); nb > 0 {
				results[ref.Name] = v.Regexp.FindStringSubmatch(raw)[1:]
			} else {
				results[ref.Name] = []string{raw}
			}

		} else {
			// Otherwise store raw data
			results[ref.Name] = []string{raw}
		}
	}

	return
}

type Tweak struct {

	// Ref => Field => Regex/Value
	Source map[string]Fields `json:"source"`

	// [ ComputedData/Item/Export ] => Field => Value
	Destination map[string]Fields `json:"destination"`
}

func NewTweak(src []byte) (*Tweak, error) {

	var t Tweak
	err := json.Unmarshal(src, &t)
	return &t, err
}

func (t *Tweak) check(raw map[string]interface{}, expected map[string]Fields, dst bool) (err error) {

	// For each data
	for key, data := range raw {

		// Check if data is handled by tweak
		if fields, ok := expected[key]; ok {

			// Check fields compatibility
			if err = fields.Check(data, dst); err != nil {
				err = fmt.Errorf("invalid data %s: %s", key, err.Error())
				return
			}
		}
	}

	return
}

// Check source compatibility
func (t *Tweak) Check(sourceRaw map[string]interface{}, destinationRaw map[string]interface{}) (err error) {

	if err = t.check(sourceRaw, t.Source, false); err != nil {
		return
	}

	if err = t.check(destinationRaw, t.Destination, true); err != nil {
		return
	}

	return
}

var dataIdReg = regexp.MustCompile(`:([a-z0-9]+)-([a-z0-9]+)(-(\d+))?`)

func setResult(format string, raw map[string]map[string][]string) (result string, err error) {

	// Get all ids from format
	submatches := dataIdReg.FindAllStringSubmatch(format, -1)

	// No submatch : nothing to do
	if len(submatches) == 0 {
		return
	}

	toReplace := make(map[string]string)

	for _, submatch := range submatches {

		if len(submatch) != 5 {
			err = fmt.Errorf("invalid submatch size ' from '%s', expected:5 get:%q", format, submatch)
			return
		}

		// Get all importants data
		key := submatch[0]
		dataKey := submatch[1]
		fieldKey := submatch[2]

		var index int
		if submatch[4] != "" {
			index, err = strconv.Atoi(submatch[4])
			if err != nil {
				err = fmt.Errorf("invalid atoi with '%s' from '%s'", submatch[4], format)
				return
			}
		}

		// Search into raw data
		fields, ok := raw[dataKey]
		if ok == false {
			err = fmt.Errorf("data key '%s' not found into raw %+v", dataKey, raw)
			return
		}

		// Get all values
		values, ok := fields[fieldKey]
		if ok == false {
			err = fmt.Errorf("field key '%s' from data '%s' not found into raw %+v", fieldKey, dataKey, raw)
			return
		}

		// Check index validity
		if index >= len(values) {
			err = fmt.Errorf("invalid index %d into max %d field key '%s' from data '%s' in raw %+v",
				index, len(values), fieldKey, dataKey, raw)
			return
		}

		toReplace[key] = values[index]
	}

	// Replace everything
	for old, new := range toReplace {
		format = strings.Replace(format, old, new, -1)
	}

	result = format
	return
}

// Tweak does the transformations between received data and data/field results
func (t *Tweak) Tweak(src map[string]interface{}) (results map[string]map[string]string, err error) {

	raw := make(map[string]map[string][]string)

	// Get selected from data sources
	for key, getters := range t.Source {

		// Check if we specified data are handled
		source, ok := src[key]
		if ok == false {
			continue
		}

		// Check field names & get raw datas & apply regexp
		raw[key], err = getters.GetRawDatas(source)
		if err != nil {
			return
		}
	}

	// No data : no tweak to do!
	if len(raw) == 0 {
		return
	}

	// Based on ordered values : determined results for each result field
	results = make(map[string]map[string]string)

	for name, fields := range t.Destination {

		if results[name] == nil {
			results[name] = make(map[string]string)
		}

		for key, v := range fields {

			results[name][key], err = setResult(v.Value, raw)
			if err != nil {
				err = fmt.Errorf("invalid destination name '%s' key '%s' %s",
					name, key, err.Error())
				return
			}
		}
	}

	return
}

type HasTweak interface {
	GetTweak(*Collection) *Tweak
	SetTweak(*Collection, *Tweak) error

	GetDatas() map[string]interface{}
}

func (c *Classify) GetTweak(t HasTweak, collection *Collection) *Tweak {

	res := t.GetTweak(collection)
	if res == nil {
		return new(Tweak)
	}

	return res
}

func (c *Classify) SetInputTweak(in HasTweak, collection *Collection, new *Tweak) (err error) {

	// Check tweak compatibility
	if err = new.Check(in.GetDatas(), collection.GetDatas()); err != nil {
		return
	}

	// Store tweak
	return in.SetTweak(collection, new)
}

func (c *Classify) SetExportTweak(out HasTweak, collection *Collection, new *Tweak) (err error) {

	// Check tweak compatibility
	if err = new.Check(collection.GetDatas(), out.GetDatas()); err != nil {
		return
	}

	// Store tweak
	return out.SetTweak(collection, new)
}
