package services

import (
	"fmt"
	"strings"

	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
	"github.com/google/uuid"
)

// RoomService handles the business logic for room generation and management
type RoomService struct {
	balancer Balancer
}

// NewRoomService creates a new RoomService with the required dependencies
func NewRoomService() (*RoomService, error) {
	// Create a balancer with the same repository
	balancer := NewBalancer()

	// Return the service with the repository interface
	return &RoomService{
		balancer: balancer,
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

// ItemConfig contains parameters for item generation
type ItemConfig struct {
	Key         string             // Item key for lookup
	Name        string             // Item name for display
	Count       int                // Number of this item type to add
	RandomPlace bool               // Whether to place items randomly
	Position    *entities.Position // Optional specific position (only used if RandomPlace is false)
}

// PlaceableConfig defines the interface for any placeable entity configuration
type PlaceableConfig interface {
	// CreatePlaceable creates a new placeable entity from this configuration
	// It may use the repository to fetch additional data if needed
	CreatePlaceable(s *RoomService) (entities.Placeable, error)

	// ShouldPlaceRandomly returns whether the entity should be placed at a random position
	ShouldPlaceRandomly() bool

	// GetPosition returns the specified position for the entity, if any
	GetPosition() *entities.Position

	// GetName returns a display name for error messages
	GetName() string
}

// Ensure our config types implement PlaceableConfig
var _ PlaceableConfig = (*MonsterConfig)(nil)
var _ PlaceableConfig = (*PlayerConfig)(nil)
var _ PlaceableConfig = (*ItemConfig)(nil)

// CreatePlaceable implements PlaceableConfig for MonsterConfig
func (c MonsterConfig) CreatePlaceable(s *RoomService) (entities.Placeable, error) {
	monster := &entities.Monster{
		ID:   uuid.NewString(),
		Name: c.Name,
		Key:  c.Key,
		CR:   c.CR,
	}
	return monster, nil
}

// ShouldPlaceRandomly implements PlaceableConfig for MonsterConfig
func (c MonsterConfig) ShouldPlaceRandomly() bool {
	return c.RandomPlace
}

// GetPosition implements PlaceableConfig for MonsterConfig
func (c MonsterConfig) GetPosition() *entities.Position {
	return c.Position
}

// GetName implements PlaceableConfig for MonsterConfig
func (c MonsterConfig) GetName() string {
	return c.Name
}

// CreatePlaceable implements PlaceableConfig for PlayerConfig
func (c PlayerConfig) CreatePlaceable(s *RoomService) (entities.Placeable, error) {
	player := &entities.Player{
		ID:    uuid.NewString(),
		Name:  c.Name,
		Level: c.Level,
	}
	return player, nil
}

// ShouldPlaceRandomly implements PlaceableConfig for PlayerConfig
func (c PlayerConfig) ShouldPlaceRandomly() bool {
	return c.RandomPlace
}

// GetPosition implements PlaceableConfig for PlayerConfig
func (c PlayerConfig) GetPosition() *entities.Position {
	return c.Position
}

// GetName implements PlaceableConfig for PlayerConfig
func (c PlayerConfig) GetName() string {
	return c.Name
}

// CreatePlaceable implements PlaceableConfig for ItemConfig
func (c ItemConfig) CreatePlaceable(s *RoomService) (entities.Placeable, error) {
	item := &entities.Item{
		ID:   uuid.NewString(),
		Key:  c.Key,
		Name: c.Name,
	}

	return item, nil
}

// ShouldPlaceRandomly implements PlaceableConfig for ItemConfig
func (c ItemConfig) ShouldPlaceRandomly() bool {
	return c.RandomPlace
}

// GetPosition implements PlaceableConfig for ItemConfig
func (c ItemConfig) GetPosition() *entities.Position {
	return c.Position
}

// GetName implements PlaceableConfig for ItemConfig
func (c ItemConfig) GetName() string {
	return c.Name
}

// AddPlaceablesToRoom adds any placeable entities to a room based on their configurations
// Players will always be placed first. If the room becomes full, monsters and items may be discarded
// with a warning message rather than causing an error.
func (s *RoomService) AddPlaceablesToRoom(room *entities.Room, configs []PlaceableConfig) error {
	if room == nil {
		return fmt.Errorf("room cannot be nil")
	}

	if len(configs) == 0 {
		return fmt.Errorf("at least one placeable entity must be provided")
	}

	// Group configs by entity type for prioritization
	playerConfigs := []PlaceableConfig{}
	monsterConfigs := []PlaceableConfig{}
	itemConfigs := []PlaceableConfig{}
	otherConfigs := []PlaceableConfig{}

	// First pass: categorize configs without creating entities
	for _, config := range configs {
		// We can identify PlayerConfig, MonsterConfig, and ItemConfig by their type
		switch config.(type) {
		case PlayerConfig:
			playerConfigs = append(playerConfigs, config)
		case MonsterConfig:
			monsterConfigs = append(monsterConfigs, config)
		case ItemConfig:
			itemConfigs = append(itemConfigs, config)
		default:
			otherConfigs = append(otherConfigs, config)
		}
	}

	// Combine in priority order: players, monsters, items, others
	prioritizedConfigs := append(playerConfigs, monsterConfigs...)
	prioritizedConfigs = append(prioritizedConfigs, itemConfigs...)
	prioritizedConfigs = append(prioritizedConfigs, otherConfigs...)

	// Track which entities couldn't be placed
	var discardedEntities []string

	for _, config := range prioritizedConfigs {
		// Create the placeable entity
		entity, err := config.CreatePlaceable(s)
		if err != nil {
			return err
		}

		// Get the entity type for logging
		entityType := "entity"
		switch entity.GetCellType() {
		case entities.CellPlayer:
			entityType = "player"
		case entities.CellMonster:
			entityType = "monster"
		case entities.CellItem:
			entityType = "item"
		}

		// Place entity either randomly or at a specific position
		if config.ShouldPlaceRandomly() {
			position, err := FindEmptyPosition(room)
			if err != nil {
				// For players, this is a critical error
				if entity.GetCellType() == entities.CellPlayer {
					return fmt.Errorf("failed to place %s (player): %w", config.GetName(), err)
				}

				// For monsters and items, just log and continue
				discardedEntities = append(discardedEntities, fmt.Sprintf("%s (%s)", config.GetName(), entityType))
				continue
			}
			entity.SetPosition(position)
		} else if pos := config.GetPosition(); pos != nil {
			// Use the specified position
			entity.SetPosition(*pos)
		} else {
			return fmt.Errorf("%s must have a position when RandomPlace is false", config.GetName())
		}

		// Add the entity to the room using the interface-based method
		if err := PlaceEntity(room, entity); err != nil {
			// For players with specific positions, this is a critical error
			if entity.GetCellType() == entities.CellPlayer {
				return fmt.Errorf("failed to add %s (player): %w", config.GetName(), err)
			}

			// For monsters and items, just log and continue
			discardedEntities = append(discardedEntities, fmt.Sprintf("%s (%s)", config.GetName(), entityType))
			continue
		}
	}

	// Log a warning if any entities were discarded
	if len(discardedEntities) > 0 {
		fmt.Printf("Warning: Could not place %d entities in the room because it was full: %s\n",
			len(discardedEntities), strings.Join(discardedEntities, ", "))
	}

	return nil
}

// GenerateAndPopulateRoom creates a room and populates it with any combination of monsters, players, and items
// At least one entity type must be provided
// If party is provided, monsters will be balanced according to the specified difficulty (defaults to Medium)
func (s *RoomService) GenerateAndPopulateRoom(
	roomConfig RoomConfig,
	monsterConfigs []MonsterConfig,
	playerConfigs []PlayerConfig,
	itemConfigs []ItemConfig,
	party *entities.Party,
	difficulty entities.EncounterDifficulty,
) (*entities.Room, error) {
	// Check that at least one entity type is provided
	if len(monsterConfigs) == 0 && len(playerConfigs) == 0 && len(itemConfigs) == 0 {
		return nil, fmt.Errorf("at least one entity type (monster, player, or item) must be provided")
	}

	// If party is provided and we have monsters, balance them according to difficulty
	if party != nil && len(monsterConfigs) > 0 {
		// Default to medium difficulty if not specified
		if difficulty == "" {
			difficulty = entities.EncounterDifficultyMedium
		}

		// Balance monster configurations
		var err error
		monsterConfigs, err = s.balancer.AdjustMonsterSelection(monsterConfigs, *party, difficulty)
		if err != nil {
			return nil, fmt.Errorf("failed to balance monsters: %w", err)
		}
	}

	// First generate the room
	room, err := s.GenerateRoom(roomConfig)
	if err != nil {
		return nil, err
	}

	// Collect all placeable configs into a single slice
	placeableConfigs := []PlaceableConfig{}

	// Add monster configs
	for _, config := range monsterConfigs {
		for i := 0; i < config.Count; i++ {
			instanceConfig := config
			placeableConfigs = append(placeableConfigs, instanceConfig)
		}
	}

	// Add player configs
	for _, config := range playerConfigs {
		placeableConfigs = append(placeableConfigs, config)
	}

	// Add item configs
	for _, config := range itemConfigs {
		for i := 0; i < config.Count; i++ {
			instanceConfig := config
			placeableConfigs = append(placeableConfigs, instanceConfig)
		}
	}

	// If we have entities to add, add them all at once
	if len(placeableConfigs) > 0 {
		if err := s.AddPlaceablesToRoom(room, placeableConfigs); err != nil {
			return nil, err
		}
	}

	return room, nil
}

// CleanupRoom removes entities from a room and returns XP gained for monsters
// If entityIDs is empty for a type, all entities of that type are removed
// Returns the total XP gained, a slice of entity IDs that weren't removed, and any error encountered
func (s *RoomService) CleanupRoom(room *entities.Room, entityType entities.CellType, entityIDs []string) (int, []string, error) {
	if room == nil {
		return 0, nil, fmt.Errorf("room cannot be nil")
	}

	// Track total XP gained and entities not removed
	totalXP := 0
	notRemoved := []string{}

	switch entityType {
	case entities.CellMonster:
		// If entityIDs is empty, remove all monsters
		if len(entityIDs) == 0 {
			// First calculate XP for all monsters
			for _, monster := range room.Monsters {
				// Use the explicit XP value if provided, otherwise calculate based on CR
				if monster.XP > 0 {
					totalXP += monster.XP
				} else {
					// Use CR-based estimate
					totalXP += int(monster.CR * 100)
				}
			}

			// Create a copy of monster IDs to avoid modification during iteration
			monsterIDs := make([]string, len(room.Monsters))
			for i, monster := range room.Monsters {
				monsterIDs[i] = monster.ID
			}

			// Remove each monster by ID
			for _, id := range monsterIDs {
				// Find the monster entity
				var monster *entities.Monster
				for i := range room.Monsters {
					if room.Monsters[i].ID == id {
						monster = &room.Monsters[i]
						break
					}
				}

				if monster != nil {
					removed, err := RemovePlaceable(room, monster)
					if !removed || err != nil {
						notRemoved = append(notRemoved, id)
					}
				} else {
					notRemoved = append(notRemoved, id)
				}
			}
		} else {
			// Remove specific monsters by ID
			for _, monsterID := range entityIDs {
				// Find the monster entity
				var monster *entities.Monster
				for i := range room.Monsters {
					if room.Monsters[i].ID == monsterID {
						monster = &room.Monsters[i]
						break
					}
				}

				if monster != nil {
					// Use the explicit XP value if provided, otherwise calculate based on CR
					if monster.XP > 0 {
						totalXP += monster.XP
					} else {
						// Use CR-based estimate
						totalXP += int(monster.CR * 100)
					}

					removed, err := RemovePlaceable(room, monster)
					if !removed || err != nil {
						notRemoved = append(notRemoved, monsterID)
					}
				} else {
					notRemoved = append(notRemoved, monsterID)
				}
			}
		}

	case entities.CellItem:
		// If entityIDs is empty, remove all items
		if len(entityIDs) == 0 {
			// Create a copy of item IDs to avoid modification during iteration
			itemIDs := make([]string, len(room.Items))
			for i, item := range room.Items {
				itemIDs[i] = item.ID
			}

			// Remove each item by ID
			for _, id := range itemIDs {
				// Find the item entity
				var item *entities.Item
				for i := range room.Items {
					if room.Items[i].ID == id {
						item = &room.Items[i]
						break
					}
				}

				if item != nil {
					removed, err := RemovePlaceable(room, item)
					if !removed || err != nil {
						notRemoved = append(notRemoved, id)
					}
				} else {
					notRemoved = append(notRemoved, id)
				}
			}
		} else {
			// Remove specific items by ID
			for _, itemID := range entityIDs {
				// Find the item entity
				var item *entities.Item
				for i := range room.Items {
					if room.Items[i].ID == itemID {
						item = &room.Items[i]
						break
					}
				}

				if item != nil {
					removed, err := RemovePlaceable(room, item)
					if !removed || err != nil {
						notRemoved = append(notRemoved, itemID)
					}
				} else {
					notRemoved = append(notRemoved, itemID)
				}
			}
		}

	case entities.CellPlayer:
		// If entityIDs is empty, remove all players
		if len(entityIDs) == 0 {
			// Create a copy of player IDs to avoid modification during iteration
			playerIDs := make([]string, len(room.Players))
			for i, player := range room.Players {
				playerIDs[i] = player.ID
			}

			// Remove each player by ID
			for _, id := range playerIDs {
				// Find the player entity
				var player *entities.Player
				for i := range room.Players {
					if room.Players[i].ID == id {
						player = &room.Players[i]
						break
					}
				}

				if player != nil {
					removed, err := RemovePlaceable(room, player)
					if !removed || err != nil {
						notRemoved = append(notRemoved, id)
					}
				} else {
					notRemoved = append(notRemoved, id)
				}
			}
		} else {
			// Remove specific players by ID
			for _, playerID := range entityIDs {
				// Find the player entity
				var player *entities.Player
				for i := range room.Players {
					if room.Players[i].ID == playerID {
						player = &room.Players[i]
						break
					}
				}

				if player != nil {
					removed, err := RemovePlaceable(room, player)
					if !removed || err != nil {
						notRemoved = append(notRemoved, playerID)
					}
				} else {
					notRemoved = append(notRemoved, playerID)
				}
			}
		}

	default:
		return 0, nil, fmt.Errorf("unsupported entity type: %v", entityType)
	}

	return totalXP, notRemoved, nil
}

// MoveEntity moves a placeable entity from its current position to a new position
// Returns an error if the move cannot be completed
func (s *RoomService) MoveEntity(room *entities.Room, entity entities.Placeable, newPosition entities.Position) error {
	return MovePlaceable(room, entity, newPosition)
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
	room := NewRoom(config.Width, config.Height, lightLevel)
	room.Description = config.Description

	// Initialize grid if requested
	if config.UseGrid {
		InitializeGrid(room)
	}

	return room, nil
}
