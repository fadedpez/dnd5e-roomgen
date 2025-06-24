package entities

// Placeable represents an entity that can be placed in a room
type Placeable interface {
	GetID() string
	GetPosition() Position
	SetPosition(pos Position)
	GetCellType() CellType
}

// RoomType defines the behavior of a specific type of room
// TODO: implement room types
type RoomType interface {
	// Type returns the string identifier for this room type
	Type() string

	// Description returns a human-readable description of the room type
	Description() string
}
