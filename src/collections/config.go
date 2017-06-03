package collections

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Config struct {
	BufferSize int           `json:"bufferSize"`
	Filters    CfgStringList `json:"filters"`
	Separators CfgStringList `json:"separators"`
	Banned     CfgStringList `json:"banned"`
}

func NewConfig(bufferSize int) *Config {
	return &Config{
		BufferSize: bufferSize,
	}
}

// Returns the list of string with banned chain removed
func (c *Config) clean(toClean string) (result []string, banned []string) {

	// Remove white space
	toClean = strings.TrimSpace(toClean)

	result = make([]string, 0)
	banned = make([]string, 0)

	// If no separators
	if len(c.Separators) == 0 {

		// Remove banned word
		for _, toBan := range c.Banned {

			if strings.Contains(toClean, toBan) {
				banned = append(banned, toBan)
				toClean = strings.Replace(toClean, toBan, "", -1)
			}
		}

		result = append(result, toClean)
		return
	}

	// Split words
	for _, sep := range c.Separators {

		words := strings.Split(toClean, sep)

		for _, word := range words {

			canAppend := true

			// Remove banned word
			for _, toBan := range c.Banned {
				if word == toBan {
					canAppend = false
					break
				}
			}

			if canAppend {
				result = append(result, word)
			} else {
				banned = append(banned, word)
			}
		}
	}

	return
}

func (c *Config) GetBannedList() []string {

	if len(c.Separators) == 0 {
		return c.Banned
	}

	bannedList := make([]string, 0)

	// Removed banned strings
	for _, banned := range c.Banned {

		// Separate elements with separators
		for _, separator := range c.Separators {
			banned = strings.Replace(banned, separator, "`", -1)
		}

		bannedList = append(bannedList, strings.Split(banned, "`")...)
	}

	return bannedList
}

func (c *Collection) ModifyConfig(name string, action string, list []string) (err error) {

	var currentList *CfgStringList

	params := map[string]*CfgStringList{
		"filters":    &c.config.Filters,
		"separators": &c.config.Separators,
		"banned":     &c.config.Banned,
	}

	if param, ok := params[name]; ok {
		currentList = param
	} else {
		return fmt.Errorf("Invalid config parameters '%s'", name)
	}

	if action == "add" {
		err = currentList.Add(list)
	} else if action == "remove" {
		err = currentList.Remove(list)
	} else {
		err = fmt.Errorf("Invalid config action '%s'", action)
	}

	// Get cleaned buffer items on change
	if err == nil && (name == "banned" || name == "separators") {
		c.buffer.CleanedNames(c.config.GetBannedList(), c.config.Separators)
	}

	fmt.Printf("%s %+v\n", name, currentList)

	return err
}

func (c *Collection) ModifyConfigValue(name string, value string) error {

	if name == "bufferSize" {

		// Integer is expected
		bufferSize, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("Parameters '%s' need an integer", name)
		}

		// Set new buffer size
		c.config.BufferSize = bufferSize

		// Update buffer size
		c.buffer.SetSize(bufferSize)
	} else {
		return fmt.Errorf("Invalid config parameters '%s'", name)
	}

	return nil
}

func (c *Collection) GetConfig() *Config {
	return c.config
}

type CfgStringList []string

func (c *CfgStringList) Add(toAdd []string) error {

	mapCurrent := map[string]struct{}{}
	for _, current := range *c {
		mapCurrent[current] = struct{}{}
	}

	hasChanged := false
	list := *c

	// Check if item doesn't already exist
	for _, add := range toAdd {

		if _, ok := mapCurrent[add]; ok {
			continue
		}

		list = append(list, add)
		hasChanged = true
	}

	if hasChanged {

		// Sort longest item first, short one at the end
		sort.Sort(CfgStringList(list))
		*c = list
		return nil
	}

	return fmt.Errorf("invalid config to add '%s': nothing changed",
		strings.Join(toAdd, ","))
}

func (c *CfgStringList) Remove(toRemove []string) error {

	list := *c

	mapToRemove := map[string]struct{}{}
	for _, remove := range toRemove {
		mapToRemove[remove] = struct{}{}
	}

	for i := len(list) - 1; i >= 0; i-- {

		current := list[i]

		if _, ok := mapToRemove[current]; ok {

			if i > 0 {
				list = append(list[:i], list[:i+1]...)
			} else if len(list) > 0 {
				list = list[1:]
			} else {
				list = []string{}
			}

			delete(mapToRemove, current)
		}
	}

	*c = list

	// No errors found
	if len(mapToRemove) == 0 {
		return nil
	}

	keys := make([]string, len(mapToRemove))

	i := 0
	for key, _ := range mapToRemove {
		keys[i] = key
		i++
	}

	return fmt.Errorf("invalid config to remove '%s'", strings.Join(keys, ","))
}

// Methods used to implement sort.Interface
func (c CfgStringList) Len() int {
	return len(c)
}

func (c CfgStringList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c CfgStringList) Less(i, j int) bool {
	return len(c[j]) < len(c[i])
}
