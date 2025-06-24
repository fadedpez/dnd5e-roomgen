package entities

// Placeable represents an entity that can be placed in a room
type Placeable interface {
	GetID() string
	GetPosition() Position
	SetPosition(pos Position)
	GetCellType() CellType
}
