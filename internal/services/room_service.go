package services

import (
	"fmt"

	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
	"github.com/fadedpez/dnd5e-roomgen/internal/repositories"
	"github.com/google/uuid"
)

// RoomService handles the business logic for room generation and management
type RoomService struct {
	monsterRepo repositories.MonsterRepository
	balancer    Balancer
}

// NewRoomService creates a new RoomService with the required dependencies
func NewRoomService() (*RoomService, error) {
	// Create a new API monster repository
	monsterRepo, err := repositories.NewAPIMonsterRepository()
	if err != nil {
		return nil, err
	}

	// Create a balancer with the same repository
	balancer := NewBalancer(monsterRepo)

	// Return the service with the repository interface
	return &RoomService{
		monsterRepo: monsterRepo,
		balancer:    balancer,
	}, nil
}

// RoomConfig contains all the parameters for room generation
type RoomConfig struct {
	Width       int
	Height      int
	LightLevel  entities.LightLevel
	Description string
	UseGrid     bool
}

// MonsterConfig contains parameters for monster generation
type MonsterConfig struct {
	Name        string
	Key         string
	CR          float64
	Count       int                // Number of this monster type to add
	RandomPlace bool               // Whether to place monsters randomly
	Position    *entities.Position // Optional specific position (only used if RandomPlace is false)
}

// PlayerConfig contains parameters for player character placement
type PlayerConfig struct {
	Name        string
	Level       int                // Character level
	RandomPlace bool               // Whether to place player randomly
	Position    *entities.Position // Optional specific position (only used if RandomPlace is false)
}

// GenerateRoom creates a new room based on the provided configuration
func (s *RoomService) GenerateRoom(config RoomConfig) (*entities.Room, error) {
	if config.Width <= 0 || config.Height <= 0 {
		return nil, fmt.Errorf("room dimensions must be positive")
	}

	// Set default light level if not specified
	lightLevel := config.LightLevel
	if lightLevel == "" {
		lightLevel = entities.LightLevelBright
	}

	// Create the room
	room := entities.NewRoom(config.Width, config.Height, lightLevel)
	room.Description = config.Description

	// Initialize grid if requested
	if config.UseGrid {
		entities.InitializeGrid(room)
	}

	return room, nil
}

// AddMonstersToRoom adds monsters to a room based on the provided configuration
func (s *RoomService) AddMonstersToRoom(room *entities.Room, monsterConfigs []MonsterConfig) error {
	if room == nil {
		return fmt.Errorf("room cannot be nil")
	}

	// Check if we need to initialize the grid for monster placement
	if room.Grid == nil && len(monsterConfigs) > 0 {
		entities.InitializeGrid(room)
	}

	for _, config := range monsterConfigs {
		for i := 0; i < config.Count; i++ {
			monster := entities.Monster{
				ID:   uuid.NewString(),
				Name: config.Name,
				Key:  config.Key,
				CR:   config.CR,
			}

			// Place monster either randomly or at a specific position
			if config.RandomPlace {
				position, err := entities.FindEmptyPosition(room)
				if err != nil {
					return fmt.Errorf("failed to place monster %s: %w", config.Name, err)
				}
				monster.Position = position
			} else if config.Position != nil {
				// Use the specified position
				monster.Position = *config.Position
			} else {
				return fmt.Errorf("monster %s must have a position when RandomPlace is false", config.Name)
			}

			// Add the monster to the room
			if err := entities.AddMonster(room, monster); err != nil {
				return fmt.Errorf("failed to add monster %s: %w", config.Name, err)
			}
		}
	}

	return nil
}

// AddPlayersToRoom adds players to a room based on the provided configuration
func (s *RoomService) AddPlayersToRoom(room *entities.Room, playerConfigs []PlayerConfig) error {
	if room == nil {
		return fmt.Errorf("room cannot be nil")
	}

	// Check if we need to initialize the grid for player placement
	if room.Grid == nil && len(playerConfigs) > 0 {
		entities.InitializeGrid(room)
	}

	for _, config := range playerConfigs {
		player := entities.Player{
			ID:    uuid.NewString(),
			Name:  config.Name,
			Level: config.Level,
		}

		// Place player either randomly or at a specific position
		if config.RandomPlace {
			position, err := entities.FindEmptyPosition(room)
			if err != nil {
				return fmt.Errorf("failed to place player %s: %w", config.Name, err)
			}
			player.Position = position
		} else if config.Position != nil {
			// Use the specified position
			player.Position = *config.Position
		} else {
			return fmt.Errorf("player %s must have a position when RandomPlace is false", config.Name)
		}

		// Add the player to the room
		if err := entities.AddPlayer(room, player); err != nil {
			return fmt.Errorf("failed to add player %s: %w", config.Name, err)
		}
	}

	return nil
}

// PopulateRoomWithMonsters is a convenience method that creates a room and populates it with monsters
func (s *RoomService) PopulateRoomWithMonsters(roomConfig RoomConfig, monsterConfigs []MonsterConfig) (*entities.Room, error) {
	// First generate the room
	room, err := s.GenerateRoom(roomConfig)
	if err != nil {
		return nil, err
	}

	// Then add monsters to it
	err = s.AddMonstersToRoom(room, monsterConfigs)
	if err != nil {
		return nil, err
	}

	return room, nil
}

// PopulateRoomWithMonstersAndPlayers is a convenience method that creates a room and populates it with monsters and players
func (s *RoomService) PopulateRoomWithMonstersAndPlayers(
	roomConfig RoomConfig,
	monsterConfigs []MonsterConfig,
	playerConfigs []PlayerConfig) (*entities.Room, error) {

	// First generate the room
	room, err := s.GenerateRoom(roomConfig)
	if err != nil {
		return nil, err
	}

	// Add players to the room
	err = s.AddPlayersToRoom(room, playerConfigs)
	if err != nil {
		return nil, err
	}

	// Then add monsters to it
	err = s.AddMonstersToRoom(room, monsterConfigs)
	if err != nil {
		return nil, err
	}

	return room, nil
}

// BalanceMonsterConfigs adjusts monster configurations based on party composition and desired difficulty
func (s *RoomService) BalanceMonsterConfigs(monsterConfigs []MonsterConfig, party entities.Party, difficulty entities.EncounterDifficulty) ([]MonsterConfig, error) {
	return s.balancer.AdjustMonsterSelection(monsterConfigs, party, difficulty)
}

// PopulateRoomWithBalancedMonsters generates a room and populates it with monsters balanced for the party and difficulty
func (s *RoomService) PopulateRoomWithBalancedMonsters(roomConfig RoomConfig, monsterConfigs []MonsterConfig, party entities.Party, difficulty entities.EncounterDifficulty) (*entities.Room, error) {
	// Generate the room
	room, err := s.GenerateRoom(roomConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate room: %w", err)
	}

	// Balance the monster configurations
	balancedConfigs, err := s.BalanceMonsterConfigs(monsterConfigs, party, difficulty)
	if err != nil {
		return nil, fmt.Errorf("failed to balance monster configurations: %w", err)
	}

	// Add the balanced monsters to the room
	if err := s.AddMonstersToRoom(room, balancedConfigs); err != nil {
		return nil, fmt.Errorf("failed to add monsters to room: %w", err)
	}

	return room, nil
}

// DetermineRoomDifficulty determines the difficulty of a room's monster encounter for a given party
func (s *RoomService) DetermineRoomDifficulty(room *entities.Room, party entities.Party) (entities.EncounterDifficulty, error) {
	if room == nil {
		return "", fmt.Errorf("room cannot be nil")
	}

	return s.balancer.DetermineEncounterDifficulty(room.Monsters, party)
}

// CleanupRoom removes monsters from a room and returns XP gained
// If monsterIDs is empty, all monsters are removed
// Returns the total XP gained, a slice of monster IDs that weren't removed, and any error encountered
func (s *RoomService) CleanupRoom(room *entities.Room, monsterIDs []string) (int, []string, error) {
	if room == nil {
		return 0, nil, fmt.Errorf("room cannot be nil")
	}

	// Track total XP gained and monsters not removed
	totalXP := 0
	notRemoved := []string{}

	// If monsterIDs is empty, remove all monsters
	if len(monsterIDs) == 0 {
		// Process all monsters to calculate XP
		for _, monster := range room.Monsters {
			// Get XP from the repository
			xp, err := s.monsterRepo.GetMonsterXP(monster.Key)
			if err != nil {
				// Log the error but continue
				fmt.Printf("Warning: failed to get XP for monster %s: %v\n", monster.Key, err)
			} else {
				totalXP += xp
			}
		}

		// Clear the grid of monsters
		if room.Grid != nil {
			for y := 0; y < room.Height; y++ {
				for x := 0; x < room.Width; x++ {
					if room.Grid[y][x].Type == entities.CellMonster {
						room.Grid[y][x] = entities.Cell{Type: entities.CellTypeEmpty}
					}
				}
			}
		}

		// Clear the monsters slice
		room.Monsters = make([]entities.Monster, 0)
	} else {
		// Remove specific monsters by ID
		for _, monsterID := range monsterIDs {
			// Find the monster to get its key
			var monsterKey string
			for _, monster := range room.Monsters {
				if monster.ID == monsterID {
					monsterKey = monster.Key
					break
				}
			}

			if monsterKey != "" {
				// Get XP from the repository
				xp, err := s.monsterRepo.GetMonsterXP(monsterKey)
				if err != nil {
					// Log the error but continue
					fmt.Printf("Warning: failed to get XP for monster %s: %v\n", monsterKey, err)
				} else {
					totalXP += xp
				}
			}

			// Use the existing RemoveMonster function to handle removal
			removed, err := entities.RemoveMonster(room, monsterID)
			if err != nil {
				return totalXP, notRemoved, fmt.Errorf("error removing monster %s: %w", monsterID, err)
			}

			// If monster wasn't found, add to notRemoved slice
			if !removed {
				notRemoved = append(notRemoved, monsterID)
			}
		}
	}

	return totalXP, notRemoved, nil
}
