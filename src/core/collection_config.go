package core

import (
	"fmt"
	"github.com/ohohleo/classify/config"
	"strconv"
	"strings"
)

type CollectionConfig struct {
	EnableStore  bool              `json:"enableStore"`
	EnableBuffer bool              `json:"enableBuffer"`
	BufferSize   int               `json:"bufferSize"`
	Filters      config.StringList `json:"filters"`
	Separators   config.StringList `json:"separators"`
	Banned       config.StringList `json:"banned"`
}

func NewCollectionConfig(bufferSize int) *CollectionConfig {
	return &CollectionConfig{
		BufferSize: bufferSize,
	}
}

// Returns the list of string with banned chain removed
func (c *CollectionConfig) clean(toClean string) (result []string, banned []string) {

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

func (c *CollectionConfig) GetBannedList() []string {

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

	var currentList *config.StringList

	params := map[string]*config.StringList{
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

	switch name {

	case "enableStore":
		previousStoreState := c.config.EnableStore

		c.config.EnableStore = (value == "true")
		if previousStoreState == c.config.EnableStore {
			return nil
		}

		if c.config.EnableStore {
			c.ActivateStore()
		} else {
			return c.DisableStore()
		}

	case "enableBuffer":
		previousBufferState := c.config.EnableBuffer

		c.config.EnableBuffer = (value == "true")
		if previousBufferState == c.config.EnableBuffer {
			return nil
		}

		if c.config.EnableBuffer {
			c.ActivateBuffer()
		} else {
			return c.DisableBuffer()
		}

	case "bufferSize":

		if c.config.EnableBuffer == false {
			return fmt.Errorf("Should enable buffer first")
		}

		// Integer is expected
		bufferSize, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("Parameters '%s' need an integer", name)
		}

		// Set new buffer size
		c.config.BufferSize = bufferSize

		// Update buffer size
		c.buffer.SetSize(bufferSize)

	default:
		return fmt.Errorf("Invalid config parameters '%s'", name)
	}

	return nil
}

func (c *Collection) GetConfig() *CollectionConfig {
	return c.config
}
