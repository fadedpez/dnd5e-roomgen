package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// createTestItem creates an item with the given ID and position
func createTestItem(id string, x, y int) Item {
	return Item{
		ID:        id,
		Key:       "potion-healing",
		Name:      "Potion of Healing " + id,
		Type:      "equipment",
		Category:  "potion",
		Value:     50,
		ValueUnit: "gp",
		Weight:    1,
		Position:  Position{X: x, Y: y},
	}
}

func TestAddItem(t *testing.T) {
	// Create a room with a grid
	room := NewRoom(5, 5, LightLevelBright)

	// Create an item
	item := createTestItem("1", 2, 2)

	// Add the item to the room
	err := AddItem(room, item)

	// Verify the item was added successfully
	assert.NoError(t, err)
	assert.Len(t, room.Items, 1)
	assert.Equal(t, item, room.Items[0])

	// Verify the item was placed on the grid
	assert.Equal(t, CellItem, room.Grid[item.Position.Y][item.Position.X].Type)
	assert.Equal(t, item.ID, room.Grid[item.Position.Y][item.Position.X].EntityID)

	// Test adding item to occupied cell
	item2 := createTestItem("2", 2, 2)

	err = AddItem(room, item2)

	// Should return an error
	assert.Error(t, err)
	assert.Equal(t, ErrCellOccupied, err)
	assert.Len(t, room.Items, 1) // Item count should not change

	// Test adding item with invalid position
	item3 := createTestItem("3", 10, 10)

	err = AddItem(room, item3)

	// Should return error
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPosition, err)
	assert.Len(t, room.Items, 1) // Item count should not change

	// Test adding item to nil room
	var nilRoom *Room
	item4 := createTestItem("4", 2, 2)

	err = AddItem(nilRoom, item4)

	// Should return error
	assert.Error(t, err)
	assert.Equal(t, ErrNilRoom, err)
}

func TestRemoveItem(t *testing.T) {
	// Create a room with a grid
	room := NewRoom(5, 5, LightLevelBright)

	// Add multiple items
	item1 := createTestItem("item1", 1, 1)
	item2 := createTestItem("item2", 2, 2)
	item3 := createTestItem("item3", 3, 3)

	// Add items to room
	err := AddItem(room, item1)
	assert.NoError(t, err)

	err = AddItem(room, item2)
	assert.NoError(t, err)

	err = AddItem(room, item3)
	assert.NoError(t, err)

	// Verify items were added
	assert.Len(t, room.Items, 3)

	// Remove item2 by ID
	removed, err := RemoveItem(room, "item2")

	// Verify item was removed successfully
	assert.NoError(t, err)
	assert.True(t, removed)
	assert.Len(t, room.Items, 2)

	// Verify item is no longer in the list
	for _, i := range room.Items {
		assert.NotEqual(t, "item2", i.ID)
	}

	// Verify cell is now empty
	assert.Equal(t, CellTypeEmpty, room.Grid[2][2].Type)
	assert.Empty(t, room.Grid[2][2].EntityID)

	// Try to remove an item that doesn't exist
	removed, err = RemoveItem(room, "nonexistent")
	assert.NoError(t, err)
	assert.False(t, removed)
	assert.Len(t, room.Items, 2)

	// Test removing from nil room
	var nilRoom *Room
	removed, err = RemoveItem(nilRoom, "item1")
	assert.Error(t, err)
	assert.Equal(t, ErrNilRoom, err)
	assert.False(t, removed)
}
