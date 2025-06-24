package entities

// Player represents a player character placed in the room
type Player struct {
	ID       string   // UUID for this player instance
	Name     string   // Name of the player character
	Level    int      // Level of the player character
	Position Position // Position of the player in the room (if grid is used)
}

// GetID returns the unique identifier for this player
func (p *Player) GetID() string {
	return p.ID
}

// GetPosition returns the current position of this player in the room
func (p *Player) GetPosition() Position {
	return p.Position
}

// SetPosition updates the position of this player
func (p *Player) SetPosition(pos Position) {
	p.Position = pos
}

// GetCellType returns the type of cell this player occupies
func (p *Player) GetCellType() CellType {
	return CellPlayer
}
