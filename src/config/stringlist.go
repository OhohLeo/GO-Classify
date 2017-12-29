package config

import (
	"fmt"
	"sort"
	"strings"
)

type StringList []string

func (c *StringList) Add(toAdd []string) error {

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
		sort.Sort(StringList(list))
		*c = list
		return nil
	}

	return fmt.Errorf("invalid config to add '%s': nothing changed",
		strings.Join(toAdd, ","))
}

func (c *StringList) Remove(toRemove []string) error {

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
func (c StringList) Len() int {
	return len(c)
}

func (c StringList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c StringList) Less(i, j int) bool {
	return len(c[j]) < len(c[i])
}
