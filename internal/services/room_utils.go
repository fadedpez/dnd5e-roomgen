package services

import (
	"fmt"
	"math"

	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
)

// NewRoom creates a new room with the specified dimensions
func NewRoom(width, height int, lightLevel entities.LightLevel) *entities.Room {
	room := &entities.Room{
		Width:      width,
		Height:     height,
		LightLevel: lightLevel,
		Monsters:   make([]entities.Monster, 0),
		Players:    make([]entities.Player, 0),
		Items:      make([]entities.Item, 0),
	}

	return room
}

// InitializeGrid creates and initializes the grid for a room
// All cells are initialized as empty
func InitializeGrid(room *entities.Room) {
	if room == nil {
		return
	}

	room.Grid = make([][]entities.Cell, room.Height)
	for i := range room.Grid {
		room.Grid[i] = make([]entities.Cell, room.Width)
		for j := range room.Grid[i] {
			room.Grid[i][j] = entities.Cell{Type: entities.CellTypeEmpty}
		}
	}
}

// MovePlaceable moves any placeable entity from its current position to a new position
// Returns error if:
// - Room is nil
// - Entity is nil
// - Entity is not found in the room
// - Target position is out of bounds
// - Target position is already occupied
func MovePlaceable(room *entities.Room, entity entities.Placeable, newPosition entities.Position) error {
	if room == nil {
		return entities.ErrNilRoom
	}

	if entity == nil {
		return fmt.Errorf("entity cannot be nil")
	}

	entityID := entity.GetID()
	cellType := entity.GetCellType()
	oldPosition := entity.GetPosition()

	// If room has no grid, just update the entity's position in the appropriate slice
	if room.Grid == nil {
		// Find and update the entity in the appropriate slice
		switch cellType {
		case entities.CellMonster:
			for i := range room.Monsters {
				if room.Monsters[i].ID == entityID {
					room.Monsters[i].Position = newPosition
					// Also update the passed entity
					entity.SetPosition(newPosition)
					return nil
				}
			}
		case entities.CellPlayer:
			for i := range room.Players {
				if room.Players[i].ID == entityID {
					room.Players[i].Position = newPosition
					// Also update the passed entity
					entity.SetPosition(newPosition)
					return nil
				}
			}
		case entities.CellItem:
			for i := range room.Items {
				if room.Items[i].ID == entityID {
					room.Items[i].Position = newPosition
					// Also update the passed entity
					entity.SetPosition(newPosition)
					return nil
				}
			}
		}
		return fmt.Errorf("entity with ID %s not found in room", entityID)
	}

	// For rooms with a grid, validate the new position

	// Check if new position is within bounds
	if newPosition.X < 0 || newPosition.X >= room.Width ||
		newPosition.Y < 0 || newPosition.Y >= room.Height {
		return fmt.Errorf("new position (%d, %d) is outside room bounds (%d, %d)",
			newPosition.X, newPosition.Y, room.Width, room.Height)
	}

	// Check if target cell is empty or is the entity's current position
	if room.Grid[newPosition.Y][newPosition.X].Type != entities.CellTypeEmpty &&
		(newPosition.X != oldPosition.X || newPosition.Y != oldPosition.Y) {
		return fmt.Errorf("cell (%d, %d) is already occupied", newPosition.X, newPosition.Y)
	}

	// Find and update the entity in the appropriate slice
	entityFound := false
	switch cellType {
	case entities.CellMonster:
		for i := range room.Monsters {
			if room.Monsters[i].ID == entityID {
				room.Monsters[i].Position = newPosition
				entityFound = true
				break
			}
		}
	case entities.CellPlayer:
		for i := range room.Players {
			if room.Players[i].ID == entityID {
				room.Players[i].Position = newPosition
				entityFound = true
				break
			}
		}
	case entities.CellItem:
		for i := range room.Items {
			if room.Items[i].ID == entityID {
				room.Items[i].Position = newPosition
				entityFound = true
				break
			}
		}
	}

	if !entityFound {
		return fmt.Errorf("entity with ID %s not found in room", entityID)
	}

	// Update the grid
	// Clear old position
	if oldPosition.X >= 0 && oldPosition.X < room.Width &&
		oldPosition.Y >= 0 && oldPosition.Y < room.Height {
		room.Grid[oldPosition.Y][oldPosition.X] = entities.Cell{Type: entities.CellTypeEmpty}
	}

	// Set new position
	room.Grid[newPosition.Y][newPosition.X] = entities.Cell{
		Type:     cellType,
		EntityID: entityID,
	}

	// Also update the passed entity
	entity.SetPosition(newPosition)

	return nil
}

// CalculateDistance calculates the distance between two positions using D&D 5d rules
// In D&D 5e, diagonal movement counts the same as orthognonal movement (Chebyshev distance)
// Returns the distance in grid units
func CalculateDistance(pos1, pos2 entities.Position) float64 {
	dx := float64(pos2.X - pos1.X)
	dy := float64(pos2.Y - pos1.Y)

	// D&D 5e movement (diagonal counts as 1)
	return math.Max(math.Abs(dx), math.Abs(dy))
}

// RemovePlaceable removes any placeable entity from the room
// Returns true if the entity was found and removed, false otherwise
func RemovePlaceable(room *entities.Room, entity entities.Placeable) (bool, error) {
	if room == nil {
		return false, entities.ErrNilRoom
	}
	if entity == nil {
		return false, fmt.Errorf("entity cannot be nil")
	}

	return removeEntity(room, entity.GetID(), entity.GetCellType()), nil
}
