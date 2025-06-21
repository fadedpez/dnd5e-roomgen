package entities

import (
	"fmt"
	"math"
	"math/rand"
)

// NewRoom creates a new room with the specified dimensions
func NewRoom(width, height int, lightLevel LightLevel) *Room {
	room := &Room{
		Width:      width,
		Height:     height,
		LightLevel: lightLevel,
		Monsters:   make([]Monster, 0),
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
		return fmt.Errorf("cannot add monster to nil room")
	}

	if room.Grid != nil {
		pos := monster.Position

		// Check if position is valid
		if pos.X < 0 || pos.X >= room.Width || pos.Y < 0 || pos.Y >= room.Height {
			return fmt.Errorf("monster position (%d, %d) is outside room bounds (%d, %d)",
				pos.X, pos.Y, room.Width, room.Height)
		}

		// Check if cell is already occupied
		if room.Grid[pos.Y][pos.X].Type != CellTypeEmpty {
			return fmt.Errorf("cell (%d, %d) is already occupied", pos.X, pos.Y)
		}

		// Place monster on grid
		room.Grid[pos.Y][pos.X] = Cell{
			Type:     CellMonster,
			EntityID: monster.ID,
		}
	}

	// Add to monsters slice
	room.Monsters = append(room.Monsters, monster)

	return nil
}

// FindEmptyPosition finds a random empty position in the room
// Returns an error if no empty positions are available or if the room has no grid
// Only to be used when the grid is initialized, not used if consumer elects to not use the grid feature
func FindEmptyPosition(room *Room) (Position, error) {
	if room == nil || room.Grid == nil {
		return Position{}, fmt.Errorf("cannot find empty position in nil room or room with no grid")

	}

	// Count empty cells
	emptyCells := 0
	for i := range room.Grid {
		for j := range room.Grid[i] {
			if room.Grid[i][j].Type == CellTypeEmpty {
				emptyCells++
			}
		}
	}

	// If no empty cells, return error
	if emptyCells == 0 {
		return Position{}, fmt.Errorf("no empty positions available")
	}

	// Pick a random empty cell
	targetCell := rand.Intn(emptyCells)
	currentCell := 0

	for i := range room.Grid {
		for j := range room.Grid[i] {
			if room.Grid[i][j].Type == CellTypeEmpty {
				if currentCell == targetCell {
					return Position{X: j, Y: i}, nil
				}
				currentCell++
			}
		}
	}

	// This should never happen if our counting logic is correct
	return Position{}, fmt.Errorf("no empty positions available. counted %d empty cells, but found none", emptyCells)
}

// RemoveMonster removes a monster from the room by its ID
// Returns true if the monster was found and removed, false otherwise
// If the room has a grid, the cell where the monster was is cleared
func RemoveMonster(room *Room, monsterID string) (bool, error) {
	if room == nil {
		return false, fmt.Errorf("cannot remove monster from nil room")
	}

	// Find the monster in the room's monster slice
	monsterIndex := -1
	var monsterToRemove Monster

	for i, monster := range room.Monsters {
		if monster.ID == monsterID {
			monsterIndex = i
			monsterToRemove = monster
			break
		}

	}

	// If monster not found, return false
	if monsterIndex == -1 {
		return false, nil
	}

	// If room has a grid, clear the cell where the monster was
	if room.Grid != nil {
		pos := monsterToRemove.Position

		// Check if position is valid (should always be valid if monster was added correctly)
		if pos.X >= 0 && pos.X < room.Width && pos.Y >= 0 && pos.Y < room.Height {
			// Only clear the cell if it contains this monster
			if room.Grid[pos.Y][pos.X].Type == CellMonster && room.Grid[pos.Y][pos.X].EntityID == monsterID {
				room.Grid[pos.Y][pos.X] = Cell{Type: CellTypeEmpty}
			}
		}
	}

	// Remove the monster from the room's monster slice
	room.Monsters = append(room.Monsters[:monsterIndex], room.Monsters[monsterIndex+1:]...)

	return true, nil
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
