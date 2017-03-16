package collections

type Items struct {
	items map[string]*Item
}

func NewItems(defaultSize int) *Items {
	return &Items{
		items: make(map[string]*Item),
	}
}
