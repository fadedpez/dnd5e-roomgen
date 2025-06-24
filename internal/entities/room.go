package entities

// CellType represents what occupies a cell in the room grid
type CellType int

const (
	CellTypeEmpty CellType = iota
	CellMonster
	CellItem
	CellPlayer
)

// Cell represents a single cell in the room grid
type Cell struct {
	Type     CellType // What occupies the cell
	EntityID string   // ID of the entity occupying the cell
}

// Position represents a position in the room grid
type Position struct {
	X int // X coordinate
	Y int // Y coordinate
}

// Placeable represents an entity that can be placed in a room
type Placeable interface {
	GetID() string
	GetPosition() Position
	SetPosition(pos Position)
	GetCellType() CellType
}

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

// Room represents a rectangular room in a dungeon
type Room struct {
	Width       int        // Width of the room in grid units
	Height      int        // Height of the room in grid units
	LightLevel  LightLevel // Light level of the room
	Description string     // room description
	RoomType    RoomType   // type of room
	Monsters    []Monster  // Monsters in the room
	Players     []Player   // Players in the room
	Items       []Item     // Items in the room
	Grid        [][]Cell   // Grid of cells in the room (if grid is used)
}

type LightLevel string

const (
	LightLevelBright LightLevel = "bright"
	LightLevelDim    LightLevel = "dim"
	LightLevelDark   LightLevel = "dark"
)

// EncounterDifficulty represents the difficulty level of an encounter
type EncounterDifficulty string

const (
	// EncounterDifficultyEasy represents an easy encounter
	EncounterDifficultyEasy EncounterDifficulty = "easy"

	// EncounterDifficultyMedium represents a medium encounter
	EncounterDifficultyMedium EncounterDifficulty = "medium"

	// EncounterDifficultyHard represents a hard encounter
	EncounterDifficultyHard EncounterDifficulty = "hard"

	// EncounterDifficultyDeadly represents a deadly encounter
	EncounterDifficultyDeadly EncounterDifficulty = "deadly"
)
