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
