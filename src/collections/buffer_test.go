package collections

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuffer(t *testing.T) {

	assert := assert.New(t)

	buffer := NewBuffer(&Collection{}, 2)

	// Init state
	assert.Equal(buffer.maxSize, 2)
	assert.Equal(len(buffer.waitings), 0)
	assert.Equal(len(buffer.items), 0)

	// Add one item
	item1 := &Item{}
	assert.True(buffer.Add("item1", item1))
	assert.Equal(buffer.maxSize, 2)
	assert.Equal(len(buffer.waitings), 1)
	assert.Equal(len(buffer.items), 1)

	assert.Equal([]*Item{item1}, buffer.GetCurrentList())

	// Add a 2nd item
	item2 := &Item{}
	assert.True(buffer.Add("item2", item2))
	assert.Equal(buffer.maxSize, 2)
	assert.Equal(len(buffer.waitings), 2)
	assert.Equal(len(buffer.items), 2)

	assert.Equal([]*Item{item1, item2}, buffer.GetCurrentList())

	// Add a 3rd item
	item3 := &Item{}
	assert.True(buffer.Add("item3", item3))
	assert.Equal(buffer.maxSize, 2)
	assert.Equal(len(buffer.waitings), 2)
	assert.Equal(len(buffer.items), 3)

	assert.Equal([]*Item{item1, item2},
		buffer.GetCurrentList())

	// Add a 4th item
	item4 := &Item{}
	assert.True(buffer.Add("item4", item4))
	assert.Equal(buffer.maxSize, 2)
	assert.Equal(len(buffer.waitings), 2)
	assert.Equal(len(buffer.items), 4)

	assert.Equal([]*Item{item1, item2},
		buffer.GetCurrentList())

	// Increase the size
	buffer.SetSize(3)

	assert.Equal([]*Item{item1, item2, item3},
		buffer.GetCurrentList())

	// Remove the 1st item
	assert.True(buffer.Remove("item1"))

	assert.Equal([]*Item{item2, item3, item4},
		buffer.GetCurrentList())

	// Decrease the size
	buffer.SetSize(1)

	assert.Equal([]*Item{item2},
		buffer.GetCurrentList())

	assert.True(buffer.Remove("item2"))
	assert.True(buffer.Remove("item4"))

	assert.Equal([]*Item{item3}, buffer.GetCurrentList())

	assert.True(buffer.Remove("item3"))

	assert.Equal([]*Item{}, buffer.GetCurrentList())
}
