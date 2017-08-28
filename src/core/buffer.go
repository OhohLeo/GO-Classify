package core

import (
	"fmt"
)

type Buffer struct {
	waitings   []string
	items      map[string]*BufferItem
	maxSize    int
	collection *Collection
}

func NewBuffer(collection *Collection, defaultSize int) *Buffer {
	return &Buffer{
		waitings:   make([]string, 0),
		items:      make(map[string]*BufferItem),
		maxSize:    defaultSize,
		collection: collection,
	}
}

func (b *Buffer) CleanedNames(bannedList []string, separators []string) {

	for _, item := range b.items {

		if item.SetCleanedName(bannedList, separators) {

			b.collection.SendCollectionEvent("buffer", "update", item)

			// Launch research on item
			b.collection.Search("buffer", item)
		}
	}
}

func (b *Buffer) SetSize(size int) {

	if b.maxSize == size {
		return
	}

	isSizeInc := (b.maxSize < size)

	b.maxSize = size

	// Size increase : attempt to send remaings
	if isSizeInc {
		b.SendNext("")
		return
	}

	// No need to resize
	if len(b.waitings) < size {
		return
	}

	// Resize needed
	for _, key := range b.waitings[size:] {

		// Items back to waiting status
		item, ok := b.items[key]
		if ok {
			item.Status = BUFFER_WAITING
		} else {
			fmt.Println("key '" + key + "' not in buffer")
		}
	}

	b.waitings = b.waitings[:size]
}

func (b *Buffer) Add(name string, item *BufferItem) bool {

	if _, ok := b.items[name]; ok {
		fmt.Println("already existing item '" + name + "' in buffer")
		return false
	}

	// Add to the hash list
	b.items[name] = item

	b.SendNext(name)

	return true
}

func (b *Buffer) Remove(name string) bool {

	if _, ok := b.items[name]; ok == false {
		fmt.Println("not existing item '" + name + "' in buffer")
		return false
	}

	// Remove from the waiting list
	for idx, itemName := range b.waitings {
		if itemName == name {
			b.waitings = append(b.waitings[:idx], b.waitings[idx+1:]...)
			break
		}
	}

	// Remove the hash list
	delete(b.items, name)

	b.SendNext("")

	return true
}

func (b *Buffer) Validate(name string) *BufferItem {

	item, ok := b.items[name]
	if ok {
		b.Remove(name)
	}

	return item
}

func (b *Buffer) GetCurrentList() (items []*BufferItem) {

	items = make([]*BufferItem, 0)

	for _, key := range b.waitings {

		item, ok := b.items[key]
		if ok == false {
			fmt.Println("key '" + key + "' not in buffer")
			continue
		}

		items = append(items, item)
	}

	return items
}

func (b *Buffer) SendNext(next string) {

	// No items to send
	if len(b.waitings) >= b.maxSize {
		return
	}

	// Search for next item to send
	itemsNb := b.maxSize - len(b.waitings)

	items := make([]*BufferItem, 0)
	keys := make([]string, 0)

	if len(next) != 0 {

		item, ok := b.items[next]
		if ok {
			item.Status = BUFFER_SENDING
			items = append(items, item)
			keys = append(keys, next)
			itemsNb--

			// Launch research on item
			b.collection.Search("buffer", item)
		} else {
			fmt.Println("key '" + next + "' not in buffer")
		}
	}

	// Search for next items
	for key, item := range b.items {

		// No need to send more items
		if itemsNb <= 0 {
			break
		}

		if item.Status > BUFFER_WAITING {
			continue
		}

		item.Status = BUFFER_SENDING
		items = append(items, item)
		keys = append(keys, key)
		itemsNb--

		// Launch research on item
		b.collection.Search("buffer", item)
	}

	// Store keys
	b.waitings = append(b.waitings, keys...)
}
