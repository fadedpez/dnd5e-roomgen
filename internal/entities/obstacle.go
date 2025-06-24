package entities

// Obstacle represents any physical obstacle in a room (wall, furniture, etc.)
type Obstacle struct {
	ID       string   // Unique identifier
	Name     string   // Descriptive name of the obstacle
	Key      string   // Key for identifying the obstacle type
	Position Position // Position in the room
	Blocking bool     // Whether the obstacle blocks movement
}

// GetID implements Placeable for Obstacle
func (o *Obstacle) GetID() string {
	return o.ID
}

// GetPosition implements Placeable for Obstacle
func (o *Obstacle) GetPosition() Position {
	return o.Position
}

// SetPosition implements Placeable for Obstacle
func (o *Obstacle) SetPosition(pos Position) {
	o.Position = pos
}

// GetCellType implements Placeable for Obstacle
func (o *Obstacle) GetCellType() CellType {
	return CellObstacle
}
