package services

import (
	"errors"
	"math/rand"

	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
)

// Error constants for placement operations
var (
	ErrNoEmptyPositions = errors.New("no empty positions available in room")
)

// PlaceEntity adds a placeable entity to a room at its current position
// If the position is invalid or the cell is occupied, returns an error
// For gridless rooms (room.Grid == nil), position validation is skipped
func PlaceEntity(room *entities.Room, entity entities.Placeable) error {
	if room == nil {
		return entities.ErrNilRoom
	}

	// For rooms with a grid, validate position before adding to slices
	if room.Grid != nil {
		pos := entity.GetPosition()

		// Check if position is within room boundaries
		if pos.X < 0 || pos.X >= room.Width ||
			pos.Y < 0 || pos.Y >= room.Height {
			return entities.ErrInvalidPosition
		}

		// Check if cell is already occupied
		if room.Grid[pos.Y][pos.X].Type != entities.CellTypeEmpty {
			return entities.ErrCellOccupied
		}
	}

	// Add entity to the appropriate slice based on its type
	switch entity.GetCellType() {
	case entities.CellMonster:
		if monster, ok := entity.(*entities.Monster); ok {
			room.Monsters = append(room.Monsters, *monster)
		}
	case entities.CellPlayer:
		if player, ok := entity.(*entities.Player); ok {
			room.Players = append(room.Players, *player)
		}
	case entities.CellItem:
		if item, ok := entity.(*entities.Item); ok {
			room.Items = append(room.Items, *item)
		}
	}

	// If this is a gridless room, we're done
	if room.Grid == nil {
		return nil
	}

	pos := entity.GetPosition()

	// Update grid
	room.Grid[pos.Y][pos.X] = entities.Cell{
		Type:     entity.GetCellType(),
		EntityID: entity.GetID(),
	}

	return nil
}

// removeEntity removes a placeable entity from a room by ID and cell type
// Returns true if the entity was found and removed, false otherwise
// For gridless rooms (room.Grid == nil), grid updates are skipped
func removeEntity(room *entities.Room, entityID string, cellType entities.CellType) bool {
	if room == nil {
		return false
	}

	// Find and remove the entity based on its type
	switch cellType {
	case entities.CellMonster:
		for i, monster := range room.Monsters {
			if monster.ID == entityID {
				// Clear grid cell if grid exists
				if room.Grid != nil {
					pos := monster.Position
					room.Grid[pos.Y][pos.X] = entities.Cell{
						Type:     entities.CellTypeEmpty,
						EntityID: "",
					}
				}

				// Remove monster from slice
				room.Monsters = append(room.Monsters[:i], room.Monsters[i+1:]...)
				return true
			}
		}
	case entities.CellPlayer:
		for i, player := range room.Players {
			if player.ID == entityID {
				// Clear grid cell if grid exists
				if room.Grid != nil {
					pos := player.Position
					room.Grid[pos.Y][pos.X] = entities.Cell{
						Type:     entities.CellTypeEmpty,
						EntityID: "",
					}
				}

				// Remove player from slice
				room.Players = append(room.Players[:i], room.Players[i+1:]...)
				return true
			}
		}
	case entities.CellItem:
		for i, item := range room.Items {
			if item.ID == entityID {
				// Clear grid cell if grid exists
				if room.Grid != nil {
					pos := item.Position
					room.Grid[pos.Y][pos.X] = entities.Cell{
						Type:     entities.CellTypeEmpty,
						EntityID: "",
					}
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
// For gridless rooms (room.Grid == nil), returns a random position within room dimensions
func FindEmptyPosition(room *entities.Room) (entities.Position, error) {
	if room == nil {
		return entities.Position{}, entities.ErrNilRoom
	}

	// For gridless rooms, return a random position within room dimensions
	if room.Grid == nil {
		return entities.Position{
			X: rand.Intn(room.Width),
			Y: rand.Intn(room.Height),
		}, nil
	}

	// Try to find an empty position
	emptyCells := []entities.Position{}
	for y := 0; y < room.Height; y++ {
		for x := 0; x < room.Width; x++ {
			if room.Grid[y][x].Type == entities.CellTypeEmpty {
				emptyCells = append(emptyCells, entities.Position{X: x, Y: y})
			}
		}
	}

	if len(emptyCells) == 0 {
		return entities.Position{}, ErrNoEmptyPositions
	}

	// Return a random empty position
	return emptyCells[rand.Intn(len(emptyCells))], nil
}
