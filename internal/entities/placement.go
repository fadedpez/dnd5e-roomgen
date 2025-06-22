package entities

import (
	"errors"
	"math/rand"
)

// Error constants for placement operations
var (
	ErrNoEmptyPositions = errors.New("no empty positions available in room")
)

// PlaceEntity adds a placeable entity to a room at its current position
// If the position is invalid or the cell is occupied, returns an error
func PlaceEntity(room *Room, entity Placeable) error {
	if room == nil {
		return ErrNilRoom
	}

	pos := entity.GetPosition()

	// Check if position is within room boundaries
	if pos.X < 0 || pos.X >= room.Width ||
		pos.Y < 0 || pos.Y >= room.Height {
		return ErrInvalidPosition
	}

	// Check if cell is already occupied
	if room.Grid[pos.Y][pos.X].Type != CellTypeEmpty {
		return ErrCellOccupied
	}

	// Update grid
	room.Grid[pos.Y][pos.X] = Cell{
		Type:     entity.GetCellType(),
		EntityID: entity.GetID(),
	}

	// Add entity to the appropriate slice based on its type
	switch entity.GetCellType() {
	case CellMonster:
		if monster, ok := entity.(*Monster); ok {
			room.Monsters = append(room.Monsters, *monster)
		}
	case CellPlayer:
		if player, ok := entity.(*Player); ok {
			room.Players = append(room.Players, *player)
		}
	case CellItem:
		if item, ok := entity.(*Item); ok {
			room.Items = append(room.Items, *item)
		}
	}

	return nil
}

// RemoveEntity removes a placeable entity from a room by ID and cell type
// Returns true if the entity was found and removed, false otherwise
func RemoveEntity(room *Room, entityID string, cellType CellType) bool {
	if room == nil {
		return false
	}

	// Find and remove the entity based on its type
	switch cellType {
	case CellMonster:
		for i, monster := range room.Monsters {
			if monster.ID == entityID {
				// Clear grid cell
				pos := monster.Position
				room.Grid[pos.Y][pos.X] = Cell{
					Type:     CellTypeEmpty,
					EntityID: "",
				}

				// Remove monster from slice
				room.Monsters = append(room.Monsters[:i], room.Monsters[i+1:]...)
				return true
			}
		}
	case CellPlayer:
		for i, player := range room.Players {
			if player.ID == entityID {
				// Clear grid cell
				pos := player.Position
				room.Grid[pos.Y][pos.X] = Cell{
					Type:     CellTypeEmpty,
					EntityID: "",
				}

				// Remove player from slice
				room.Players = append(room.Players[:i], room.Players[i+1:]...)
				return true
			}
		}
	case CellItem:
		for i, item := range room.Items {
			if item.ID == entityID {
				// Clear grid cell
				pos := item.Position
				room.Grid[pos.Y][pos.X] = Cell{
					Type:     CellTypeEmpty,
					EntityID: "",
				}

				// Remove item from slice
				room.Items = append(room.Items[:i], room.Items[i+1:]...)
				return true
			}
		}
	}

	return false
}

// FindEmptyPosition finds an empty position in the room
// Returns the position and nil error if successful, or an error if no empty position is found
func FindEmptyPosition(room *Room) (Position, error) {
	if room == nil {
		return Position{}, ErrNilRoom
	}

	// Try to find an empty position
	emptyCells := []Position{}
	for y := 0; y < room.Height; y++ {
		for x := 0; x < room.Width; x++ {
			if room.Grid[y][x].Type == CellTypeEmpty {
				emptyCells = append(emptyCells, Position{X: x, Y: y})
			}
		}
	}

	if len(emptyCells) == 0 {
		return Position{}, ErrNoEmptyPositions
	}

	// Return a random empty position
	return emptyCells[rand.Intn(len(emptyCells))], nil
}
