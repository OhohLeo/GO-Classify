package collections

import (
	"fmt"
	"strings"
)

type Config struct {
	BufferSize int           `json:"buffer_size"`
	Filters    CfgStringList `json:"filters"`
	Separators CfgStringList `json:"separators"`
	Banned     CfgStringList `json:"banned"`
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

func (c *Collection) ModifyConfig(name string, action string, list []string) error {

	var currentList CfgStringList

	params := map[string]CfgStringList{
		"filters":    c.config.Filters,
		"separators": c.config.Separators,
		"banned":     c.config.Banned,
	}

	if param, ok := params[name]; ok {
		currentList = param
	} else {
		return fmt.Errorf("Invalid config parameters '%s'", name)
	}

	if action == "add" {
		currentList.Add(list)
	} else if action == "remove" {
		return currentList.Remove(list)
	} else {
		return fmt.Errorf("Invalid config action '%s'", action)
	}

	return nil
}

func (c *Collection) ModifyConfigValue(name string, value int) error {

	if name == "buffer_size" {
		c.config.BufferSize = value
	} else {
		return fmt.Errorf("Invalid config parameters '%s'", name)
	}

	return nil
}

func (c *Collection) GetConfig() *Config {
	return c.config
}

type CfgStringList []string

func (c CfgStringList) Add(toAdd []string) {
	c = append(c, toAdd...)
}

func (c CfgStringList) Remove(toRemove []string) error {

	mapToRemove := map[string]struct{}{}
	for _, remove := range toRemove {
		mapToRemove[remove] = struct{}{}
	}

	for i := len(c) - 1; i >= 0; i-- {

		current := c[i]

		if _, ok := mapToRemove[current]; ok {
			delete(mapToRemove, current)
			continue
		}

		c = append(c[:i], c[:i+1]...)
	}

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
