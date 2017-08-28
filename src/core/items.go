package core

import (
	"fmt"
)

type Items struct {
	items map[string]Data
}

func NewItems() *Items {
	return &Items{
		items: make(map[string]Data),
	}
}

func (i *Items) Add(id string, item Data) error {

	if _, ok := i.items[id]; ok {
		return fmt.Errorf("already existing item '" + id + "' in items")
	}

	// Add to the hash list
	i.items[id] = item

	return nil
}

func (i *Items) Remove(id string) error {

	if _, ok := i.items[id]; ok == false {
		return fmt.Errorf("not existing item '" + id + "' in items")
	}

	// Remove the hash list
	delete(i.items, id)

	return nil
}

func (i *Items) GetCurrentList() (items []Data) {

	items = make([]Data, 0)

	for _, item := range i.items {
		items = append(items, item)
	}

	return items
}
