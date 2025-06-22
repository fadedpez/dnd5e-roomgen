package entities

import (
	"fmt"
	"math"
)

// NewRoom creates a new room with the specified dimensions
func NewRoom(width, height int, lightLevel LightLevel) *Room {
	room := &Room{
		Width:      width,
		Height:     height,
		LightLevel: lightLevel,
		Monsters:   make([]Monster, 0),
		Players:    make([]Player, 0),
		Items:      make([]Item, 0),
	}

	return room
}

// InitializeGrid creates and initializes the grid for a room
// All cells are initialized as empty
func InitializeGrid(room *Room) {
	if room == nil {
		return
	}

	room.Grid = make([][]Cell, room.Height)
	for i := range room.Grid {
		room.Grid[i] = make([]Cell, room.Width)
		for j := range room.Grid[i] {
			room.Grid[i][j] = Cell{Type: CellTypeEmpty}
		}
	}
}

// AddMonster adds a monster to the room and places it on the grid if available
func AddMonster(room *Room, monster Monster) error {
	if room == nil {
		return ErrNilRoom
	}

	return PlaceEntity(room, &monster)
}

// RemoveMonster removes a monster from the room by its ID
// Returns true if the monster was found and removed, false otherwise
// If the room has a grid, the cell where the monster was is cleared
func RemoveMonster(room *Room, monsterID string) (bool, error) {
	if room == nil {
		return false, ErrNilRoom
	}

	// Use the generic RemoveEntity function but adapt the return value
	removed := RemoveEntity(room, monsterID, CellMonster)
	return removed, nil
}

// MoveEntity moves an entity (like a player or monster) from its current position to a new position
// Returns error if:
// - Room is nil
// - Entity is not found in the room
// - Target position is out of bounds
// - Target position is already occupied
func MoveEntity(room *Room, entityID string, newPosition Position) error {
	if room == nil {
		return fmt.Errorf("cannot move entity in nil room")
	}

	// Find the entity in the room
	entityIndex := -1
	for i, monster := range room.Monsters {
		if monster.ID == entityID {
			entityIndex = i
			break
		}
	}

	if entityIndex == -1 {
		return fmt.Errorf("entity with ID %s not found in room", entityID)
	}

	// Store the old position
	oldPos := room.Monsters[entityIndex].Position

	// If room has no grid, just update position
	if room.Grid == nil {

		room.Monsters[entityIndex].Position = newPosition
		return nil
	}

	// Validate new position is within bounds
	if newPosition.X < 0 || newPosition.X >= room.Width ||
		newPosition.Y < 0 || newPosition.Y >= room.Height {
		return fmt.Errorf("new position (%d, %d) is outside room bounds (%d, %d)",
			newPosition.X, newPosition.Y, room.Width, room.Height)
	}

	// Check if target cell is empty
	if room.Grid[newPosition.Y][newPosition.X].Type != CellTypeEmpty {
		return fmt.Errorf("cell (%d, %d) is already occupied", newPosition.X, newPosition.Y)
	}

	// Move the entity
	room.Monsters[entityIndex].Position = newPosition

	// Clear the old cell
	room.Grid[oldPos.Y][oldPos.X] = Cell{Type: CellTypeEmpty}

	// Update the entity's position
	room.Monsters[entityIndex].Position = newPosition

	// Update the grid
	room.Grid[newPosition.Y][newPosition.X] = Cell{Type: CellMonster, EntityID: entityID}

	return nil
}

// CalculateDistance calculates the distance between two positions using D&D 5d rules
// In D&D 5e, diagonal movement counts the same as orthognonal movement (Chebyshev distance)
// Returns the distance in grid units
func CalculateDistance(pos1, pos2 Position) float64 {
	dx := float64(pos2.X - pos1.X)
	dy := float64(pos2.Y - pos1.Y)

	// D&D 5e movement (diagonal counts as 1)
	return math.Max(math.Abs(dx), math.Abs(dy))
}

// AddPlayer adds a player to the room and places it on the grid if available
func AddPlayer(room *Room, player Player) error {
	if room == nil {
		return ErrNilRoom
	}

	return PlaceEntity(room, &player)
}

// RemovePlayer removes a player from the room by their ID
// Returns true if the player was found and removed, false otherwise
// If the room has a grid, the cell where the player was is cleared
func RemovePlayer(room *Room, playerID string) (bool, error) {
	if room == nil {
		return false, ErrNilRoom
	}

	// Use the generic RemoveEntity function but adapt the return value
	removed := RemoveEntity(room, playerID, CellPlayer)
	return removed, nil
}
