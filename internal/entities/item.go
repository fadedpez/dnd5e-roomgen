package entities

import (
	"errors"
)

// Error constants for item operations
var (
	ErrNilRoom         = errors.New("room is nil")
	ErrInvalidPosition = errors.New("position is outside room boundaries")
	ErrCellOccupied    = errors.New("cell is already occupied")
)

// Item represents a treasure item placed in the room
type Item struct {
	ID                  string   // UUID for this item instance
	Key                 string   // Reference key for the item in the API
	Name                string   // Name of the item
	Type                string   // Type of item (equipment, weapon, armor)
	Category            string   // Equipment category or weapon/armor category
	Value               int      // Gold value of the item (from Cost.Quantity)
	ValueUnit           string   // Currency unit (from Cost.Unit)
	Weight              int      // Weight of the item
	Position            Position // Position of the item in the room
	Properties          []string // Special properties (for weapons)
	DamageDice          string   // Damage dice (for weapons)
	DamageType          string   // Type of damage (for weapons)
	ArmorClass          int      // Base armor class (for armor)
	StealthDisadvantage bool     // Whether armor gives disadvantage on stealth checks
}

// GetID returns the unique identifier for this item
func (i *Item) GetID() string {
	return i.ID
}

// GetPosition returns the current position of this item in the room
func (i *Item) GetPosition() Position {
	return i.Position
}

// SetPosition updates the position of this item
func (i *Item) SetPosition(pos Position) {
	i.Position = pos
}

// GetCellType returns the type of cell this item occupies
func (i *Item) GetCellType() CellType {
	return CellItem
}

// ItemConfig represents configuration for an item to be placed in a room
type ItemConfig struct {
	Key      string   // Reference key from the API
	Position Position // Optional position override
}
