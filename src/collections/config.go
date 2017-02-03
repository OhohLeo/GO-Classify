package collections

import (
	"strings"
)

type Config struct {
	ImportFilters []string `json:"import_filters"`
	Separators    []string `json:"separators"`
	WordBanned    []string `json:"word_banned"`
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
		for _, ban := range c.WordBanned {

			if strings.Contains(toClean, ban) {
				banned = append(banned, ban)
				toClean = strings.Replace(toClean, ban, "", -1)
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
			for _, ban := range c.WordBanned {
				if word == ban {
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

func (c *Collection) SetConfig(string, []string) error {
	return nil
}

func (c *Collection) GetConfig() *Config {
	return c.config
}
