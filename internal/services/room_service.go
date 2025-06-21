package services

import (
	"fmt"

	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
	"github.com/google/uuid"
)

// RoomService handles the business logic for room generation and management
type RoomService struct {
	// we can add dependencies here later (monster repos, etc)
}

func NewRoomService() *RoomService {
	return &RoomService{}
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
