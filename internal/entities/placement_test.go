package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockPlaceable implements the Placeable interface for testing
type MockPlaceable struct {
	id       string
	position Position
	cellType CellType
}

func (m *MockPlaceable) GetID() string {
	return m.id
}

func (m *MockPlaceable) GetPosition() Position {
	return m.position
}

func (m *MockPlaceable) SetPosition(pos Position) {
	m.position = pos
}

func (m *MockPlaceable) GetCellType() CellType {
	return m.cellType
}

func TestPlaceEntity(t *testing.T) {
	// Create a room with a grid
	room := NewRoom(5, 5, LightLevelBright)

	// Test placing different entity types
	testCases := []struct {
		name        string
		entity      Placeable
		expectedErr error
	}{
		{
			name: "Place monster",
			entity: &Monster{
				ID:       "monster1",
				Key:      "goblin",
				Name:     "Test Goblin",
				CR:       0.25,
				Position: Position{X: 1, Y: 1},
			},
			expectedErr: nil,
		},
		{
			name: "Place player",
			entity: &Player{
				ID:       "player1",
				Name:     "Test Player",
				Level:    5,
				Position: Position{X: 2, Y: 2},
			},
			expectedErr: nil,
		},
		{
			name: "Place item",
			entity: &Item{
				ID:       "item1",
				Key:      "potion",
				Name:     "Test Potion",
				Type:     "equipment",
				Position: Position{X: 3, Y: 3},
			},
			expectedErr: nil,
		},
		{
			name: "Place at invalid position",
			entity: &MockPlaceable{
				id:       "mock1",
				position: Position{X: 10, Y: 10},
				cellType: CellMonster,
			},
			expectedErr: ErrInvalidPosition,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := PlaceEntity(room, tc.entity)

			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)

				// Verify entity was placed on grid
				pos := tc.entity.GetPosition()
				assert.Equal(t, tc.entity.GetCellType(), room.Grid[pos.Y][pos.X].Type)
				assert.Equal(t, tc.entity.GetID(), room.Grid[pos.Y][pos.X].EntityID)
			}
		})
	}

	// Test placing in occupied cell
	occupiedPos := Position{X: 1, Y: 1}
	mockEntity := &MockPlaceable{
		id:       "mock2",
		position: occupiedPos,
		cellType: CellMonster,
	}

	err := PlaceEntity(room, mockEntity)
	assert.Equal(t, ErrCellOccupied, err)
}

func TestRemoveEntity(t *testing.T) {
	// Create a room with a grid
	room := NewRoom(5, 5, LightLevelBright)

	// Add entities of different types
	monster := &Monster{
		ID:       "monster1",
		Key:      "goblin",
		Name:     "Test Goblin",
		CR:       0.25,
		Position: Position{X: 1, Y: 1},
	}

	player := &Player{
		ID:       "player1",
		Name:     "Test Player",
		Level:    5,
		Position: Position{X: 2, Y: 2},
	}

	item := &Item{
		ID:       "item1",
		Key:      "potion",
		Name:     "Test Potion",
		Type:     "equipment",
		Position: Position{X: 3, Y: 3},
	}

	// Place entities
	err := PlaceEntity(room, monster)
	assert.NoError(t, err)

	err = PlaceEntity(room, player)
	assert.NoError(t, err)

	err = PlaceEntity(room, item)
	assert.NoError(t, err)

	// Test removing entities
	testCases := []struct {
		name          string
		entityID      string
		cellType      CellType
		expectRemoved bool
	}{
		{
			name:          "Remove monster",
			entityID:      "monster1",
			cellType:      CellMonster,
			expectRemoved: true,
		},
		{
			name:          "Remove player",
			entityID:      "player1",
			cellType:      CellPlayer,
			expectRemoved: true,
		},
		{
			name:          "Remove item",
			entityID:      "item1",
			cellType:      CellItem,
			expectRemoved: true,
		},
		{
			name:          "Remove non-existent entity",
			entityID:      "nonexistent",
			cellType:      CellMonster,
			expectRemoved: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			removed := RemoveEntity(room, tc.entityID, tc.cellType)
			assert.Equal(t, tc.expectRemoved, removed)

			if tc.expectRemoved {
				// Verify entity is no longer in the grid
				var found bool
				for y := range room.Grid {
					for x := range room.Grid[y] {
						if room.Grid[y][x].EntityID == tc.entityID {
							found = true
						}
					}
				}
				assert.False(t, found, "Entity should be removed from grid")
			}
		})
	}
}

func TestFindEmptyPositionWithFullRoom(t *testing.T) {
	// Create a room with a grid
	room := NewRoom(3, 3, LightLevelBright)

	// Fill the room completely
	for y := 0; y < room.Height; y++ {
		for x := 0; x < room.Width; x++ {
			entity := &MockPlaceable{
				id:       "entity-" + string(rune('A'+y)) + string(rune('1'+x)),
				position: Position{X: x, Y: y},
				cellType: CellMonster,
			}
			err := PlaceEntity(room, entity)
			assert.NoError(t, err)
		}
	}

	// Try to find an empty position
	_, err := FindEmptyPosition(room)
	assert.Equal(t, ErrNoEmptyPositions, err)
}
