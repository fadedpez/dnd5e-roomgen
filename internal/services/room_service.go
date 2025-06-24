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

// NPCConfig contains parameters for NPC placement
type NPCConfig struct {
	Name        string
	Level       int                // Character level
	Count       int                // Number of this NPC type to add
	Inventory   []entities.Item    // Items in the NPC's inventory
	RandomPlace bool               // Whether to place NPC randomly
	Position    *entities.Position // Optional specific position (only used if RandomPlace is false)
}

// ObstacleConfig contains parameters for obstacle placement
type ObstacleConfig struct {
	Name        string             // Name of the obstacle
	Key         string             // Key for identifying the obstacle type
	Blocking    bool               // Whether the obstacle blocks movement
	Count       int                // Number of this obstacle type to add
	RandomPlace bool               // Whether to place obstacle randomly
	Position    *entities.Position // Optional specific position (only used if RandomPlace is false)
}

// ShouldPlaceRandomly implements PlaceableConfig for NPCConfig
func (c NPCConfig) ShouldPlaceRandomly() bool {
	return c.RandomPlace
}

// GetPosition implements PlaceableConfig for NPCConfig
func (c NPCConfig) GetPosition() *entities.Position {
	return c.Position
}

// GetName implements PlaceableConfig for NPCConfig
func (c NPCConfig) GetName() string {
	return c.Name
}

// CreatePlaceable implements PlaceableConfig for NPCConfig
func (c NPCConfig) CreatePlaceable(s *RoomService) (entities.Placeable, error) {
	npc := &entities.NPC{
		ID:        uuid.NewString(),
		Name:      c.Name,
		Inventory: c.Inventory,
	}
	return npc, nil
}

// GetCellType implements PlaceableConfig for NPCConfig
func (c NPCConfig) GetCellType() entities.CellType {
	return entities.CellNPC
}

// ShouldPlaceRandomly implements PlaceableConfig for ObstacleConfig
func (c ObstacleConfig) ShouldPlaceRandomly() bool {
	return c.RandomPlace
}

// GetPosition implements PlaceableConfig for ObstacleConfig
func (c ObstacleConfig) GetPosition() *entities.Position {
	return c.Position
}

// GetName implements PlaceableConfig for ObstacleConfig
func (c ObstacleConfig) GetName() string {
	return c.Name
}

// CreatePlaceable implements PlaceableConfig for ObstacleConfig
func (c ObstacleConfig) CreatePlaceable(s *RoomService) (entities.Placeable, error) {
	obstacle := &entities.Obstacle{
		ID:       uuid.NewString(),
		Name:     c.Name,
		Key:      c.Key,
		Blocking: c.Blocking,
	}
	return obstacle, nil
}

// GetCellType implements PlaceableConfig for ObstacleConfig
func (c ObstacleConfig) GetCellType() entities.CellType {
	return entities.CellObstacle
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

	// GetCellType returns the cell type for the entity
	GetCellType() entities.CellType
}

// Ensure our config types implement PlaceableConfig
var _ PlaceableConfig = (*MonsterConfig)(nil)
var _ PlaceableConfig = (*PlayerConfig)(nil)
var _ PlaceableConfig = (*ItemConfig)(nil)
var _ PlaceableConfig = (*NPCConfig)(nil)
var _ PlaceableConfig = (*ObstacleConfig)(nil)

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

// GetCellType implements PlaceableConfig for MonsterConfig
func (c MonsterConfig) GetCellType() entities.CellType {
	return entities.CellMonster
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

// GetCellType implements PlaceableConfig for PlayerConfig
func (c PlayerConfig) GetCellType() entities.CellType {
	return entities.CellPlayer
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

// GetCellType implements PlaceableConfig for ItemConfig
func (c ItemConfig) GetCellType() entities.CellType {
	return entities.CellItem
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
	npcConfigs := []PlaceableConfig{}
	obstacleConfigs := []PlaceableConfig{}
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
		case NPCConfig:
			npcConfigs = append(npcConfigs, config)
		case ObstacleConfig:
			obstacleConfigs = append(obstacleConfigs, config)
		default:
			otherConfigs = append(otherConfigs, config)
		}
	}

	// Combine in priority order: players, monsters, NPCs, obstacles, items, others
	prioritizedConfigs := append(playerConfigs, monsterConfigs...)
	prioritizedConfigs = append(prioritizedConfigs, npcConfigs...)
	prioritizedConfigs = append(prioritizedConfigs, obstacleConfigs...)
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
		case entities.CellNPC:
			entityType = "npc"
		case entities.CellObstacle:
			entityType = "obstacle"
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
	npcConfigs []NPCConfig,
	obstacleConfigs []ObstacleConfig,
	party *entities.Party,
	difficulty entities.EncounterDifficulty,
) (*entities.Room, error) {
	// Check that at least one entity type is provided
	if len(monsterConfigs) == 0 && len(playerConfigs) == 0 && len(itemConfigs) == 0 && len(npcConfigs) == 0 && len(obstacleConfigs) == 0 {
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

	// Add NPC configs
	for _, config := range npcConfigs {
		for i := 0; i < config.Count; i++ {
			instanceConfig := config
			placeableConfigs = append(placeableConfigs, instanceConfig)
		}
	}

	// Add Obstacle configs
	for _, config := range obstacleConfigs {
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

	case entities.CellNPC:
		// If entityIDs is empty, remove all NPCs
		if len(entityIDs) == 0 {
			// Create a copy of NPC IDs to avoid modification during iteration
			npcIDs := make([]string, len(room.NPCs))
			for i, npc := range room.NPCs {
				npcIDs[i] = npc.ID
			}

			// Remove each NPC by ID
			for _, id := range npcIDs {
				// Find the NPC entity
				var npc *entities.NPC
				for i := range room.NPCs {
					if room.NPCs[i].ID == id {
						npc = &room.NPCs[i]
						break
					}
				}

				if npc != nil {
					removed, err := RemovePlaceable(room, npc)
					if !removed || err != nil {
						notRemoved = append(notRemoved, id)
					}
				} else {
					notRemoved = append(notRemoved, id)
				}
			}
		} else {
			// Remove specific NPCs by ID
			for _, npcID := range entityIDs {
				// Find the NPC entity
				var npc *entities.NPC
				for i := range room.NPCs {
					if room.NPCs[i].ID == npcID {
						npc = &room.NPCs[i]
						break
					}
				}

				if npc != nil {
					removed, err := RemovePlaceable(room, npc)
					if !removed || err != nil {
						notRemoved = append(notRemoved, npcID)
					}
				} else {
					notRemoved = append(notRemoved, npcID)
				}
			}
		}

	case entities.CellObstacle:
		// If entityIDs is empty, remove all obstacles
		if len(entityIDs) == 0 {
			// Create a copy of obstacle IDs to avoid modification during iteration
			obstacleIDs := make([]string, len(room.Obstacles))
			for i, obstacle := range room.Obstacles {
				obstacleIDs[i] = obstacle.ID
			}

			// Remove each obstacle by ID
			for _, id := range obstacleIDs {
				// Find the obstacle entity
				var obstacle *entities.Obstacle
				for i := range room.Obstacles {
					if room.Obstacles[i].ID == id {
						obstacle = &room.Obstacles[i]
						break
					}
				}

				if obstacle != nil {
					removed, err := RemovePlaceable(room, obstacle)
					if !removed || err != nil {
						notRemoved = append(notRemoved, id)
					}
				} else {
					notRemoved = append(notRemoved, id)
				}
			}
		} else {
			// Remove specific obstacles by ID
			for _, obstacleID := range entityIDs {
				// Find the obstacle entity
				var obstacle *entities.Obstacle
				for i := range room.Obstacles {
					if room.Obstacles[i].ID == obstacleID {
						obstacle = &room.Obstacles[i]
						break
					}
				}

				if obstacle != nil {
					removed, err := RemovePlaceable(room, obstacle)
					if !removed || err != nil {
						notRemoved = append(notRemoved, obstacleID)
					}
				} else {
					notRemoved = append(notRemoved, obstacleID)
				}
			}
		}

	default:
		return 0, notRemoved, fmt.Errorf("unsupported entity type: %d", entityType)
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

// AddItemToNPCInventory adds an item to an NPC's inventory in the room
// Returns an error if the NPC is not found
func (s *RoomService) AddItemToNPCInventory(room *entities.Room, npcID string, item entities.Item) error {
	if room == nil {
		return entities.ErrNilRoom
	}

	npc, _ := FindNPCByID(room, npcID)
	if npc == nil {
		return fmt.Errorf("NPC with ID %s not found in room", npcID)
	}

	// Create a copy of the item with a new ID to ensure uniqueness
	itemCopy := item
	itemCopy.ID = uuid.NewString()

	npc.AddItemToInventory(itemCopy)
	return nil
}

// GetNPCInventory returns all items in an NPC's inventory
// Returns an error if the NPC is not found
func (s *RoomService) GetNPCInventory(room *entities.Room, npcID string) ([]entities.Item, error) {
	if room == nil {
		return nil, entities.ErrNilRoom
	}

	npc, _ := FindNPCByID(room, npcID)
	if npc == nil {
		return nil, fmt.Errorf("NPC with ID %s not found in room", npcID)
	}

	return npc.GetInventory(), nil
}

// RemoveItemFromNPCInventory removes an item from an NPC's inventory by ID
// Returns the removed item and an error if any occurred
func (s *RoomService) RemoveItemFromNPCInventory(room *entities.Room, npcID string, itemID string) (entities.Item, error) {
	if room == nil {
		return entities.Item{}, entities.ErrNilRoom
	}

	npc, _ := FindNPCByID(room, npcID)
	if npc == nil {
		return entities.Item{}, fmt.Errorf("NPC with ID %s not found in room", npcID)
	}

	item, success := npc.RemoveItemFromInventory(itemID)
	if !success {
		fmt.Printf("Warning: Item with ID %s not found in NPC %s's inventory\n", itemID, npc.Name)
		return entities.Item{}, fmt.Errorf("item with ID %s not found in NPC's inventory", itemID)
	}

	return item, nil
}
