package core

import (
	"fmt"
)

type Items struct {
	items map[string]*Item
}

func NewItems() *Items {
	return &Items{
		items: make(map[string]*Item),
	}
}

func (i *Items) Add(id string, item *Item) error {

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

func (i *Items) GetCurrentList() (items []*Item) {

	items = make([]*Item, 0)

	for _, item := range i.items {
		items = append(items, item)
	}

	return items
}
