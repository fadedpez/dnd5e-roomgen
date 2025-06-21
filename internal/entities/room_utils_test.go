package entities

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createTestRoom creates a standard room with a grid for testing
func createTestRoom() *Room {
	room := NewRoom(5, 5, LightLevelBright)
	InitializeGrid(room)
	return room
}

// createTestRoomNoGrid creates a standard room without a grid for testing
func createTestRoomNoGrid() *Room {
	return NewRoom(5, 5, LightLevelBright)
}

// createTestMonster creates a monster with the given ID and position
func createTestMonster(id string, x, y int) Monster {
	return Monster{
		ID:       id,
		Key:      "goblin",
		Name:     "Goblin " + id,
		CR:       0.25,
		Position: Position{X: x, Y: y},
	}
}

func TestNewRoom(t *testing.T) {
	room := createTestRoom()

	assert.NotNil(t, room)
	assert.Equal(t, 5, room.Width)
	assert.Equal(t, 5, room.Height)
	assert.Equal(t, LightLevelBright, room.LightLevel)
	assert.Empty(t, room.Monsters)

}

func TestInitializeGrid(t *testing.T) {
	// Create a room
	room := createTestRoom()

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
			assert.Equal(t, CellTypeEmpty, room.Grid[i][j].Type)
		}
	}
}

func TestAddMonster(t *testing.T) {
	// Create a room with a grid
	room := createTestRoom()

	// Create a monster
	monster := createTestMonster("1", 2, 2)

	// Add the monster to the room
	err := AddMonster(room, monster)

	// Verify the monster was added successfully
	assert.NoError(t, err)
	assert.Len(t, room.Monsters, 1)
	assert.Equal(t, monster, room.Monsters[0])

	// Verify the monster was placed on the grid
	assert.Equal(t, CellMonster, room.Grid[monster.Position.Y][monster.Position.X].Type)
	assert.Equal(t, monster.ID, room.Grid[monster.Position.Y][monster.Position.X].EntityID)

	// Test adding monster to occupied cell
	monster2 := createTestMonster("2", 2, 2)

	err = AddMonster(room, monster2)

	// Should return an error
	assert.Error(t, err)
	assert.Len(t, room.Monsters, 1) // Monster count should not change

	// Test adding monster with invalid position
	monster3 := createTestMonster("3", 10, 10)

	err = AddMonster(room, monster3)

	// Should return error
	assert.Error(t, err)
	assert.Len(t, room.Monsters, 1) // Monster count should not change

	// Test adding monster to room without grid
	roomNoGrid := createTestRoomNoGrid()
	monster4 := createTestMonster("4", 2, 2)

	err = AddMonster(roomNoGrid, monster4)

	// Should return succeed (no grid to check)
	assert.NoError(t, err)
	assert.Len(t, roomNoGrid.Monsters, 1)
	assert.Equal(t, monster4, roomNoGrid.Monsters[0])
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
	assert.Equal(t, CellTypeEmpty, room.Grid[pos.Y][pos.X].Type)

	// Fill the room with monsters except for one cell
	for i := range room.Grid {
		for j := range room.Grid[i] {
			// Leave one cell empty
			if i == 2 && j == 3 {
				continue
			}

			monster := createTestMonster(fmt.Sprintf("%d-%d", i, j), j, i)
			err := AddMonster(room, monster)
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

	err = AddMonster(room, lastMonster)
	assert.NoError(t, err)

	// Try to find an empty position in a full room
	pos, err = FindEmptyPosition(room)

	// Should return error
	assert.Error(t, err)

	// Test with room with no grid
	roomNoGrid := createTestRoomNoGrid()
	pos, err = FindEmptyPosition(roomNoGrid)
	assert.Error(t, err)
}

func TestRemoveMonster(t *testing.T) {
	// Create a room with a grid
	room := createTestRoom()

	// Add multiple monsters
	monster1 := createTestMonster("monster1", 1, 1)
	monster2 := createTestMonster("monster2", 2, 2)
	monster3 := createTestMonster("monster3", 3, 3)

	// Add monsters to room
	err := AddMonster(room, monster1)
	assert.NoError(t, err)

	err = AddMonster(room, monster2)
	assert.NoError(t, err)

	err = AddMonster(room, monster3)
	assert.NoError(t, err)

	// Verify monsters were added
	assert.Len(t, room.Monsters, 3)

	// Remove monster2 by ID
	removed, err := RemoveMonster(room, "monster2")

	// Verify monster was removed successfully
	assert.NoError(t, err)
	assert.True(t, removed)
	assert.Len(t, room.Monsters, 2)

	// Verify monster is no longer in the list
	for _, m := range room.Monsters {
		assert.NotEqual(t, "monster2", m.ID)
	}

	// Verify cell is now empty
	assert.Equal(t, CellTypeEmpty, room.Grid[2][2].Type)
	assert.Empty(t, room.Grid[2][2].EntityID)

	// Try to remove a monster that doesn't exist
	removed, err = RemoveMonster(room, "nonexistent")
	assert.NoError(t, err)
	assert.False(t, removed)
	assert.Len(t, room.Monsters, 2)

	// Test removing from room without grid
	roomNoGrid := createTestRoomNoGrid()
	monster4 := createTestMonster("monster4", 1, 1)

	err = AddMonster(roomNoGrid, monster4)
	assert.NoError(t, err)

	removed, err = RemoveMonster(roomNoGrid, "monster4")
	assert.NoError(t, err)
	assert.True(t, removed)
	assert.Len(t, roomNoGrid.Monsters, 0)
}

func TestMoveEntity(t *testing.T) {
	// Create a room with a grid
	room := createTestRoom()

	// Add a monster to the room
	monster := createTestMonster("monster1", 1, 1)
	err := AddMonster(room, monster)
	assert.NoError(t, err)

	// Verify initial position
	assert.Equal(t, 1, room.Monsters[0].Position.X)
	assert.Equal(t, 1, room.Monsters[0].Position.Y)
	assert.Equal(t, CellMonster, room.Grid[1][1].Type)
	assert.Equal(t, monster.ID, room.Grid[1][1].EntityID)

	// Move the monster to a new position
	newPos := Position{X: 3, Y: 3}
	err = MoveEntity(room, monster.ID, newPos)
	assert.NoError(t, err)

	// Verify the monster was moved
	assert.Equal(t, 3, room.Monsters[0].Position.X)
	assert.Equal(t, 3, room.Monsters[0].Position.Y)

	// Verify the new cell has the monster
	assert.Equal(t, CellMonster, room.Grid[3][3].Type)
	assert.Equal(t, monster.ID, room.Grid[3][3].EntityID)

	// Verify the old cell is now empty
	assert.Equal(t, CellTypeEmpty, room.Grid[1][1].Type)
	assert.Empty(t, room.Grid[1][1].EntityID)

	// Test moving to an occupied cell
	monster2 := createTestMonster("monster2", 2, 2)
	err = AddMonster(room, monster2)
	assert.NoError(t, err)

	// Try to move monster1 to monster2's position
	err = MoveEntity(room, monster.ID, Position{X: 2, Y: 2})
	assert.Error(t, err) // Should fail as the cell is occupied

	// Verify monster1 is still at its previous position
	assert.Equal(t, 3, room.Monsters[0].Position.X)
	assert.Equal(t, 3, room.Monsters[0].Position.Y)

	// Test moving outside room bounds
	err = MoveEntity(room, monster.ID, Position{X: 15, Y: 15})
	assert.Error(t, err) // Should fail as the cell is out of bounds

	// Test moving non-existent monster
	err = MoveEntity(room, "nonexistent", Position{X: 2, Y: 2})
	assert.Error(t, err) // Should fail as the monster doesn't exist

	// Test moving in a room without grid
	roomNoGrid := createTestRoomNoGrid()
	InitializeGrid(roomNoGrid)
	monsterNoGrid := createTestMonster("monsterNoGrid", 1, 1)
	err = AddMonster(roomNoGrid, monsterNoGrid)
	assert.NoError(t, err)

	// Move the monster to a new position
	newPos = Position{X: 3, Y: 3}
	err = MoveEntity(roomNoGrid, monsterNoGrid.ID, newPos)
	assert.NoError(t, err)

	// Verify the monster was moved
	assert.Equal(t, 3, roomNoGrid.Monsters[0].Position.X)
	assert.Equal(t, 3, roomNoGrid.Monsters[0].Position.Y)
}

func TestCalculateDistance(t *testing.T) {
	testCases := []struct {
		name     string
		pos1     Position
		pos2     Position
		expected float64
	}{
		{
			name:     "Same position",
			pos1:     Position{X: 3, Y: 3},
			pos2:     Position{X: 3, Y: 3},
			expected: 0,
		},
		{
			name:     "Horizontal movement",
			pos1:     Position{X: 1, Y: 3},
			pos2:     Position{X: 5, Y: 3},
			expected: 4,
		},
		{
			name:     "Vertical movement",
			pos1:     Position{X: 3, Y: 1},
			pos2:     Position{X: 3, Y: 5},
			expected: 4,
		},
		{
			name:     "Diagonal movement",
			pos1:     Position{X: 1, Y: 1},
			pos2:     Position{X: 5, Y: 5},
			expected: 4, // In D&D 5e, diagonal movement is the same as the larger of horizontal or vertical
		},
		{
			name:     "L-shaped movement",
			pos1:     Position{X: 1, Y: 1},
			pos2:     Position{X: 5, Y: 3},
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
