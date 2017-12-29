package core

import (
	"fmt"
)

type Items struct {
	items map[Id]*Item
}

func NewItems() *Items {
	return &Items{
		items: make(map[Id]*Item),
	}
}

func (i *Items) Add(id Id, item *Item) error {

	if _, ok := i.items[id]; ok {
		return fmt.Errorf("already existing item '%d' in items", id)
	}

	// Add to the hash list
	i.items[id] = item

	return nil
}

func (i *Items) Remove(id Id) error {

	if _, err := i.Get(id); err != nil {
		return err
	}

	// Remove the hash list
	delete(i.items, id)

	return nil
}

func (i *Items) Get(id Id) (item *Item, err error) {

	var ok bool
	if item, ok = i.items[id]; ok == false {
		err = fmt.Errorf("item '%d' not existing", id)
	}

	return
}

func (i *Items) GetCurrentList() (items []*Item) {

	items = make([]*Item, 0)

	for _, item := range i.items {
		items = append(items, item)
	}

	return items
}
