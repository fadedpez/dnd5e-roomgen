package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
)

// createTestRoom creates a standard room with a grid for testing
func createTestRoom() *entities.Room {
	room := NewRoom(5, 5, entities.LightLevelBright)
	InitializeGrid(room) // Initialize the grid for testing
	return room
}

// createTestRoomNoGrid creates a standard room without a grid for testing
func createTestRoomNoGrid() *entities.Room {
	room := &entities.Room{
		Width:      5,
		Height:     5,
		LightLevel: entities.LightLevelBright,
		RoomType:   entities.DefaultRoomType(),
		Monsters:   make([]entities.Monster, 0),
		Players:    make([]entities.Player, 0),
		Items:      make([]entities.Item, 0),
	}
	return room
}

// createTestMonster creates a monster with the given ID and position
func createTestMonster(id string, x, y int) entities.Monster {
	return entities.Monster{
		ID:       id,
		Key:      "goblin",
		Name:     "Goblin " + id,
		CR:       0.25,
		Position: entities.Position{X: x, Y: y},
	}
}

// createTestPlayer creates a player with the given ID, level, and position
func createTestPlayer(id string, level int, x, y int) entities.Player {
	return entities.Player{
		ID:       id,
		Name:     "Player " + id,
		Level:    level,
		Position: entities.Position{X: x, Y: y},
	}
}

func TestNewRoom(t *testing.T) {
	room := NewRoom(5, 5, entities.LightLevelBright)

	assert.NotNil(t, room)
	assert.Equal(t, 5, room.Width)
	assert.Equal(t, 5, room.Height)
	assert.Equal(t, entities.LightLevelBright, room.LightLevel)
	assert.Empty(t, room.Monsters)
	assert.Empty(t, room.Players)
	assert.Empty(t, room.Items)

	// Grid is optional, so we don't check for it
	// RoomType is not set by NewRoom, so we don't check for it
}

func TestInitializeGrid(t *testing.T) {
	// Create a room
	room := createTestRoomNoGrid()

	// Grid should be nil initially
	assert.Nil(t, room.Grid)

	// Initialize the grid
	InitializeGrid(room)

	assert.NotNil(t, room.Grid)
	assert.Equal(t, 5, len(room.Grid))
	assert.Equal(t, 5, len(room.Grid[0]))

	// Check that all cells are initialized as empty
	for i := range room.Grid {
		for j := range room.Grid[i] {
			assert.Equal(t, entities.CellTypeEmpty, room.Grid[i][j].Type)
		}
	}
}

func TestFindEmptyPosition(t *testing.T) {
	// Create a room with a grid
	room := createTestRoom()

	// Find an empty position
	pos, err := FindEmptyPosition(room)

	// Should succeed
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, pos.X, 0)
	assert.Less(t, pos.X, room.Width)
	assert.GreaterOrEqual(t, pos.Y, 0)
	assert.Less(t, pos.Y, room.Height)

	// Verify the position is empty
	assert.Equal(t, entities.CellTypeEmpty, room.Grid[pos.Y][pos.X].Type)

	// Fill the room with monsters except for one cell
	for i := range room.Grid {
		for j := range room.Grid[i] {
			// Leave one cell empty
			if i == 2 && j == 3 {
				continue
			}

			monster := createTestMonster(fmt.Sprintf("%d-%d", i, j), j, i)
			err := PlaceEntity(room, &monster)
			assert.NoError(t, err)
		}
	}

	// Find an empty position again
	pos, err = FindEmptyPosition(room)

	// Should find the one empty cell at i 2, j 3
	assert.NoError(t, err)
	assert.Equal(t, 3, pos.X)
	assert.Equal(t, 2, pos.Y)

	// Fill the last empty cell
	lastMonster := createTestMonster("last", pos.X, pos.Y)

	err = PlaceEntity(room, &lastMonster)
	assert.NoError(t, err)

	// Try to find an empty position in a full room
	pos, err = FindEmptyPosition(room)

	// Should return error
	assert.Error(t, err)

	// Test with room with no grid
	roomNoGrid := createTestRoomNoGrid()
	pos, err = FindEmptyPosition(roomNoGrid)
	assert.NoError(t, err)             // Should succeed for gridless rooms
	assert.GreaterOrEqual(t, pos.X, 0) // Position should be within room dimensions
	assert.Less(t, pos.X, roomNoGrid.Width)
	assert.GreaterOrEqual(t, pos.Y, 0)
	assert.Less(t, pos.Y, roomNoGrid.Height)
}

func TestCalculateDistance(t *testing.T) {
	testCases := []struct {
		name     string
		pos1     entities.Position
		pos2     entities.Position
		expected float64
	}{
		{
			name:     "Same position",
			pos1:     entities.Position{X: 3, Y: 3},
			pos2:     entities.Position{X: 3, Y: 3},
			expected: 0,
		},
		{
			name:     "Horizontal movement",
			pos1:     entities.Position{X: 1, Y: 3},
			pos2:     entities.Position{X: 5, Y: 3},
			expected: 4,
		},
		{
			name:     "Vertical movement",
			pos1:     entities.Position{X: 3, Y: 1},
			pos2:     entities.Position{X: 3, Y: 5},
			expected: 4,
		},
		{
			name:     "Diagonal movement",
			pos1:     entities.Position{X: 1, Y: 1},
			pos2:     entities.Position{X: 5, Y: 5},
			expected: 4, // In D&D 5e, diagonal movement is the same as the larger of horizontal or vertical
		},
		{
			name:     "L-shaped movement",
			pos1:     entities.Position{X: 1, Y: 1},
			pos2:     entities.Position{X: 5, Y: 3},
			expected: 4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			distance := CalculateDistance(tc.pos1, tc.pos2)
			assert.Equal(t, tc.expected, distance, "Distance calculation incorrect")

			// Distance should be the same in reverse
			reverseDistance := CalculateDistance(tc.pos2, tc.pos1)
			assert.Equal(t, tc.expected, reverseDistance, "Distance calculation incorrect in reverse")
		})
	}
}

func TestMovePlaceable(t *testing.T) {
	// Create a room with a grid
	room := createTestRoom()
	InitializeGrid(room)

	// Create and place a monster
	monster := createTestMonster("monster1", 1, 1)
	err := PlaceEntity(room, &monster)
	assert.NoError(t, err)

	// Create and place a player
	player := createTestPlayer("player1", 5, 3, 3)
	err = PlaceEntity(room, &player)
	assert.NoError(t, err)

	// Test cases
	testCases := []struct {
		name        string
		room        *entities.Room
		entity      entities.Placeable
		newPosition entities.Position
		expectError bool
	}{
		{
			name:        "Move monster to empty cell",
			room:        room,
			entity:      &monster,
			newPosition: entities.Position{X: 2, Y: 2},
			expectError: false,
		},
		{
			name:        "Move to occupied cell",
			room:        room,
			entity:      &monster,
			newPosition: entities.Position{X: 3, Y: 3}, // Where player is
			expectError: true,
		},
		{
			name:        "Move out of bounds",
			room:        room,
			entity:      &monster,
			newPosition: entities.Position{X: 10, Y: 10},
			expectError: true,
		},
		{
			name:        "Nil room",
			room:        nil,
			entity:      &monster,
			newPosition: entities.Position{X: 2, Y: 2},
			expectError: true,
		},
		{
			name:        "Nil entity",
			room:        room,
			entity:      nil,
			newPosition: entities.Position{X: 2, Y: 2},
			expectError: true,
		},
		{
			name:        "Entity not in room",
			room:        room,
			entity:      &entities.Monster{ID: "nonexistent"},
			newPosition: entities.Position{X: 2, Y: 2},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := MovePlaceable(tc.room, tc.entity, tc.newPosition)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify entity position was updated
				if tc.entity != nil {
					assert.Equal(t, tc.newPosition, tc.entity.GetPosition())

					// If room has grid, verify grid was updated
					if tc.room != nil && tc.room.Grid != nil {
						// Check new position has entity
						assert.Equal(t, tc.entity.GetCellType(), tc.room.Grid[tc.newPosition.Y][tc.newPosition.X].Type)
						assert.Equal(t, tc.entity.GetID(), tc.room.Grid[tc.newPosition.Y][tc.newPosition.X].EntityID)

						// Check old position is empty
						oldPos := entities.Position{X: 1, Y: 1} // Original position of monster
						if tc.entity == &monster && !tc.expectError {
							assert.Equal(t, entities.CellTypeEmpty, tc.room.Grid[oldPos.Y][oldPos.X].Type)
						}
					}
				}
			}
		})
	}

	// Test moving in room without grid
	roomNoGrid := createTestRoomNoGrid()
	monsterNoGrid := createTestMonster("monster-nogrid", 1, 1)
	err = PlaceEntity(roomNoGrid, &monsterNoGrid)
	assert.NoError(t, err)

	newPos := entities.Position{X: 3, Y: 3}
	err = MovePlaceable(roomNoGrid, &monsterNoGrid, newPos)
	assert.NoError(t, err)
	assert.Equal(t, newPos, monsterNoGrid.Position)
}

func TestRemovePlaceable(t *testing.T) {
	// Create a room with a grid
	room := createTestRoom()
	InitializeGrid(room)

	// Create and place entities
	monster := createTestMonster("monster1", 1, 1)
	err := PlaceEntity(room, &monster)
	assert.NoError(t, err)

	player := createTestPlayer("player1", 5, 3, 3)
	err = PlaceEntity(room, &player)
	assert.NoError(t, err)

	// Create an item
	item := entities.Item{
		ID:       "item1",
		Key:      "potion",
		Name:     "Health Potion",
		Position: entities.Position{X: 2, Y: 2},
	}
	err = PlaceEntity(room, &item)
	assert.NoError(t, err)

	// Test cases
	testCases := []struct {
		name          string
		room          *entities.Room
		entity        entities.Placeable
		expectRemoved bool
		expectError   bool
	}{
		{
			name:          "Remove monster",
			room:          room,
			entity:        &monster,
			expectRemoved: true,
			expectError:   false,
		},
		{
			name:          "Remove player",
			room:          room,
			entity:        &player,
			expectRemoved: true,
			expectError:   false,
		},
		{
			name:          "Remove item",
			room:          room,
			entity:        &item,
			expectRemoved: true,
			expectError:   false,
		},
		{
			name:          "Remove non-existent entity",
			room:          room,
			entity:        &entities.Monster{ID: "nonexistent"},
			expectRemoved: false,
			expectError:   false,
		},
		{
			name:          "Nil room",
			room:          nil,
			entity:        &monster,
			expectRemoved: false,
			expectError:   true,
		},
		{
			name:          "Nil entity",
			room:          room,
			entity:        nil,
			expectRemoved: false,
			expectError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			removed, err := RemovePlaceable(tc.room, tc.entity)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectRemoved, removed)

			// If entity was removed, verify it's no longer in the room
			if removed {
				switch tc.entity.GetCellType() {
				case entities.CellMonster:
					// Check monster was removed from slice
					for _, m := range tc.room.Monsters {
						assert.NotEqual(t, tc.entity.GetID(), m.ID)
					}
				case entities.CellPlayer:
					// Check player was removed from slice
					for _, p := range tc.room.Players {
						assert.NotEqual(t, tc.entity.GetID(), p.ID)
					}
				case entities.CellItem:
					// Check item was removed from slice
					for _, i := range tc.room.Items {
						assert.NotEqual(t, tc.entity.GetID(), i.ID)
					}
				}

				// If room has grid, verify cell is now empty
				if tc.room != nil && tc.room.Grid != nil {
					pos := tc.entity.GetPosition()
					assert.Equal(t, entities.CellTypeEmpty, tc.room.Grid[pos.Y][pos.X].Type)
					assert.Empty(t, tc.room.Grid[pos.Y][pos.X].EntityID)
				}
			}
		})
	}

	// Test removing from room without grid
	roomNoGrid := createTestRoomNoGrid()
	monsterNoGrid := createTestMonster("monster-nogrid", 1, 1)
	err = PlaceEntity(roomNoGrid, &monsterNoGrid)
	assert.NoError(t, err)

	removed, err := RemovePlaceable(roomNoGrid, &monsterNoGrid)
	assert.NoError(t, err)
	assert.True(t, removed)
	assert.Empty(t, roomNoGrid.Monsters)
}
