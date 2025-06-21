package entities

// CellType represents what occupies a cell in the room grid
type CellType int

const (
	CellTypeEmpty CellType = iota
	CellMonster
	CellItem
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

// Monster represents a monster placed in the room
type Monster struct {
	ID       string   // UUID for this monster instance
	Key      string   // Reference key from the API
	Name     string   // Name of the monster
	CR       float64  // Challenge Rating of the monster
	Position Position // Position of the monster in the room (if grid is used)

}

// Room represents a rectangular room in a dungeon
type Room struct {
	Width       int        // Width of the room in grid units
	Height      int        // Height of the room in grid units
	LightLevel  LightLevel // Light level of the room
	Description string     // "bright", "dim", or "dark"
	Monsters    []Monster  // Monsters in the room
	Grid        [][]Cell   // Grid of cells in the room (if grid is used)
}

type LightLevel string

const (
	LightLevelBright LightLevel = "bright"
	LightLevelDim    LightLevel = "dim"
	LightLevelDark   LightLevel = "dark"
)
