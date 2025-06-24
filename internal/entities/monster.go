package entities

// Monster represents a monster placed in the room
type Monster struct {
	ID       string   // UUID for this monster instance
	Key      string   // Reference key from the API
	Name     string   // Name of the monster
	CR       float64  // Challenge Rating of the monster
	XP       int      // Experience points awarded when defeated
	Position Position // Position of the monster in the room (if grid is used)
}

// GetID returns the unique identifier for this monster
func (m *Monster) GetID() string {
	return m.ID
}

// GetPosition returns the current position of this monster in the room
func (m *Monster) GetPosition() Position {
	return m.Position
}

// SetPosition updates the position of this monster
func (m *Monster) SetPosition(pos Position) {
	m.Position = pos
}

// GetCellType returns the type of cell this monster occupies
func (m *Monster) GetCellType() CellType {
	return CellMonster
}
