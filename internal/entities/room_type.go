package entities

// RoomType defines the behavior of a specific type of room
type RoomType interface {
	// Type returns the string identifier for this room type
	Type() string

	// Description returns a human-readable description of the room type
	Description() string
}

// CombatRoomType represents a room with monsters
type CombatRoomType struct{}

func (r *CombatRoomType) Type() string {
	return "combat"
}

func (r *CombatRoomType) Description() string {
	return "A room with monsters for combat encounters"
}

// TreasureRoomType represents a room with treasure items
type TreasureRoomType struct{}

func (r *TreasureRoomType) Type() string {
	return "treasure"
}

func (r *TreasureRoomType) Description() string {
	return "A room with treasure and valuable items"
}

// DefaultRoomType returns the default room type (combat)
func DefaultRoomType() RoomType {
	return &CombatRoomType{}
}
