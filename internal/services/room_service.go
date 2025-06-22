package services

import (
	"fmt"
	"math"

	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
	"github.com/fadedpez/dnd5e-roomgen/internal/repositories"
	"github.com/google/uuid"
)

// RoomService handles the business logic for room generation and management
type RoomService struct {
	monsterRepo repositories.MonsterRepository
	itemRepo    repositories.ItemRepository
	balancer    Balancer
}

// NewRoomService creates a new RoomService with the required dependencies
func NewRoomService() (*RoomService, error) {
	// Create a new API monster repository
	monsterRepo, err := repositories.NewAPIMonsterRepository()
	if err != nil {
		return nil, err
	}

	// Create a new API item repository
	itemRepo, err := repositories.NewAPIItemRepository()
	if err != nil {
		return nil, err
	}

	// Create a balancer with the same repository
	balancer := NewBalancer(monsterRepo)

	// Return the service with the repository interface
	return &RoomService{
		monsterRepo: monsterRepo,
		itemRepo:    itemRepo,
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

// ItemConfig contains parameters for item generation
type ItemConfig struct {
	Key         string             // Item key for lookup
	Name        string             // Item name for display
	Count       int                // Number of this item type to add
	RandomPlace bool               // Whether to place items randomly
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

	// We no longer force grid initialization
	// The placement interface will handle gridless rooms properly

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

	// We no longer force grid initialization
	// The placement interface will handle gridless rooms properly

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

// AddItemsToRoom adds items to a room based on the provided configuration
func (s *RoomService) AddItemsToRoom(room *entities.Room, itemConfigs []ItemConfig) error {
	if room == nil {
		return fmt.Errorf("room cannot be nil")
	}

	// We no longer force grid initialization
	// The placement interface will handle gridless rooms properly

	for _, config := range itemConfigs {
		for i := 0; i < config.Count; i++ {
			// Fetch item data from repository
			itemData, err := s.itemRepo.GetItemByKey(config.Key)
			if err != nil {
				return fmt.Errorf("failed to get item data for %s: %w", config.Key, err)
			}

			// Create a copy of the item with a unique ID
			item := entities.Item{
				ID:        uuid.NewString(),
				Key:       itemData.Key,
				Name:      itemData.Name,
				Type:      itemData.Type,
				Category:  itemData.Category,
				Value:     itemData.Value,
				ValueUnit: itemData.ValueUnit,
				Weight:    itemData.Weight,
			}

			// Place item either randomly or at a specific position
			if config.RandomPlace {
				position, err := entities.FindEmptyPosition(room)
				if err != nil {
					return fmt.Errorf("failed to place item %s: %w", config.Key, err)
				}
				item.Position = position
			} else if config.Position != nil {
				// Use the specified position
				item.Position = *config.Position
			} else {
				return fmt.Errorf("item %s must have a position when RandomPlace is false", config.Key)
			}

			// Add the item to the room
			if err := entities.AddItem(room, item); err != nil {
				return fmt.Errorf("failed to add item %s: %w", config.Key, err)
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

// PopulateRoomWithItems is a convenience method that creates a room and populates it with items
func (s *RoomService) PopulateRoomWithItems(roomConfig RoomConfig, itemConfigs []ItemConfig) (*entities.Room, error) {
	// Generate the room
	room, err := s.GenerateRoom(roomConfig)
	if err != nil {
		return nil, err
	}

	// Add items to it
	err = s.AddItemsToRoom(room, itemConfigs)
	if err != nil {
		return nil, err
	}

	return room, nil
}

// PopulateRoomWithMonstersAndItems is a convenience method that creates a room and populates it with monsters and items
func (s *RoomService) PopulateRoomWithMonstersAndItems(
	roomConfig RoomConfig,
	monsterConfigs []MonsterConfig,
	itemConfigs []ItemConfig) (*entities.Room, error) {

	// Generate the room
	room, err := s.GenerateRoom(roomConfig)
	if err != nil {
		return nil, err
	}

	// Add monsters to it
	err = s.AddMonstersToRoom(room, monsterConfigs)
	if err != nil {
		return nil, err
	}

	// Add items to it
	err = s.AddItemsToRoom(room, itemConfigs)
	if err != nil {
		return nil, err
	}

	return room, nil
}

// PopulateRoomWithAll is a convenience method that creates a room and populates it with monsters, players, and items
func (s *RoomService) PopulateRoomWithAll(
	roomConfig RoomConfig,
	monsterConfigs []MonsterConfig,
	playerConfigs []PlayerConfig,
	itemConfigs []ItemConfig) (*entities.Room, error) {

	// Generate the room
	room, err := s.GenerateRoom(roomConfig)
	if err != nil {
		return nil, err
	}

	// Add players first so they get priority placement
	err = s.AddPlayersToRoom(room, playerConfigs)
	if err != nil {
		return nil, err
	}

	// Then add monsters
	err = s.AddMonstersToRoom(room, monsterConfigs)
	if err != nil {
		return nil, err
	}

	// Finally add items
	err = s.AddItemsToRoom(room, itemConfigs)
	if err != nil {
		return nil, err
	}

	return room, nil
}

// PopulateTreasureRoom creates a room specifically designed to contain treasure items
func (s *RoomService) PopulateTreasureRoom(
	roomConfig RoomConfig,
	itemCount int,
	guardianMonsterConfigs []MonsterConfig) (*entities.Room, error) {

	// Generate the room
	room, err := s.GenerateRoom(roomConfig)
	if err != nil {
		return nil, err
	}

	// Set the room type to treasure
	room.RoomType = &entities.TreasureRoomType{}

	// Check if there's enough space for the items
	// We need at least one cell per item
	availableSpace := room.Width * room.Height

	// Account for guardian monsters if specified
	monsterCount := 0
	for _, config := range guardianMonsterConfigs {
		monsterCount += config.Count
	}

	availableSpace -= monsterCount

	if itemCount > availableSpace {
		return nil, fmt.Errorf("not enough space in room for %d items (available space: %d)",
			itemCount, availableSpace)
	}

	// Get random valuable items
	items, err := s.itemRepo.GetRandomItems(itemCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get random items: %w", err)
	}

	// Convert items to item configs
	itemConfigs := make([]ItemConfig, len(items))
	for i, item := range items {
		itemConfigs[i] = ItemConfig{
			Key:         item.Key,
			Name:        item.Name,
			Count:       1,
			RandomPlace: true,
		}
	}

	// Add items to the room
	err = s.AddItemsToRoom(room, itemConfigs)
	if err != nil {
		return nil, err
	}

	// Add guardian monsters if specified
	if len(guardianMonsterConfigs) > 0 {
		err = s.AddMonstersToRoom(room, guardianMonsterConfigs)
		if err != nil {
			return nil, err
		}
	}

	return room, nil
}

// PopulateRandomTreasureRoomWithParty creates a treasure room with loot scaled appropriately for the party
// and adds the party to the room. The room should be cleared after all treasure is collected.
func (s *RoomService) PopulateRandomTreasureRoomWithParty(
	roomConfig RoomConfig,
	party entities.Party,
	includeGuardian bool,
	difficulty entities.EncounterDifficulty) (*entities.Room, error) {

	if len(party.Members) == 0 {
		return nil, fmt.Errorf("party must have at least one member")
	}

	// Calculate average party level
	totalLevels := 0
	for _, member := range party.Members {
		totalLevels += member.Level
	}
	avgLevel := float64(totalLevels) / float64(len(party.Members))
	partySize := len(party.Members)

	// Scale item count based on party size and level
	// Base formula: 1 item per party member + 1 item per 3 levels (average)
	baseItemCount := partySize + int(avgLevel/3)

	// Adjust based on difficulty
	itemMultiplier := 1.0
	switch difficulty {
	case entities.EncounterDifficultyEasy:
		itemMultiplier = 0.75
	case entities.EncounterDifficultyMedium:
		itemMultiplier = 1.0
	case entities.EncounterDifficultyHard:
		itemMultiplier = 1.25
	case entities.EncounterDifficultyDeadly:
		itemMultiplier = 1.5
	}

	itemCount := int(float64(baseItemCount) * itemMultiplier)
	if itemCount < 1 {
		itemCount = 1 // Ensure at least one item
	}

	// Determine item categories based on party level
	var itemCategories []string

	// Add weapon categories based on level
	if avgLevel >= 5 {
		itemCategories = append(itemCategories, "martial-weapons")
	} else {
		itemCategories = append(itemCategories, "simple-weapons")
	}

	// Add armor categories based on level
	if avgLevel >= 10 {
		itemCategories = append(itemCategories, "heavy-armor")
	} else if avgLevel >= 5 {
		itemCategories = append(itemCategories, "medium-armor")
	} else {
		itemCategories = append(itemCategories, "light-armor")
	}

	// Always include potions and adventuring gear
	itemCategories = append(itemCategories, "potion", "adventuring-gear")

	// Add guardian monsters if requested
	var guardianMonsterConfigs []MonsterConfig
	if includeGuardian {
		// Create a guardian appropriate for the party's level
		guardianCR := math.Max(1.0, avgLevel-2) // CR slightly below party level

		// For higher difficulties, make the guardian tougher
		if difficulty == entities.EncounterDifficultyHard || difficulty == entities.EncounterDifficultyDeadly {
			guardianCR = math.Max(1.0, avgLevel) // CR equal to party level
		}

		guardianMonsterConfigs = []MonsterConfig{
			{
				Name:        "Guardian",
				CR:          guardianCR,
				Count:       1,
				RandomPlace: true,
			},
		}

		// Balance the guardian for the party if we have a balancer
		if s.balancer != nil {
			var err error
			guardianMonsterConfigs, err = s.balancer.AdjustMonsterSelection(guardianMonsterConfigs, party, difficulty)
			if err != nil {
				return nil, fmt.Errorf("failed to balance guardian monster: %w", err)
			}
		}
	}

	// Get items from different categories
	allItems := []*entities.Item{}
	itemsPerCategory := itemCount / len(itemCategories)
	if itemsPerCategory < 1 {
		itemsPerCategory = 1
	}

	for _, category := range itemCategories {
		items, err := s.itemRepo.GetRandomItemsByCategory(category, itemsPerCategory)
		if err != nil {
			// If we can't find items in a category, just continue
			continue
		}
		allItems = append(allItems, items...)
	}

	// If we didn't get enough items from categories, fill with random items
	if len(allItems) < itemCount {
		remainingCount := itemCount - len(allItems)
		randomItems, err := s.itemRepo.GetRandomItems(remainingCount)
		if err == nil {
			allItems = append(allItems, randomItems...)
		}
	}

	// Limit to the requested item count in case we got too many
	if len(allItems) > itemCount {
		allItems = allItems[:itemCount]
	}

	// Convert items to item configs
	itemConfigs := make([]ItemConfig, len(allItems))
	for i, item := range allItems {
		itemConfigs[i] = ItemConfig{
			Key:         item.Key,
			Name:        item.Name,
			Count:       1,
			RandomPlace: true,
		}
	}

	// Generate the room
	room, err := s.GenerateRoom(roomConfig)
	if err != nil {
		return nil, err
	}

	// Set the room type to treasure
	room.RoomType = &entities.TreasureRoomType{}

	// Check if there's enough space for items, guardians, and party members
	availableSpace := room.Width * room.Height

	// Account for guardian monsters
	monsterCount := 0
	for _, config := range guardianMonsterConfigs {
		monsterCount += config.Count
	}

	// Account for party members
	partyMemberCount := len(party.Members)

	availableSpace -= (monsterCount + partyMemberCount)

	if len(itemConfigs) > availableSpace {
		return nil, fmt.Errorf("not enough space in room for %d items, %d monsters, and %d party members (available space: %d)",
			len(itemConfigs), monsterCount, partyMemberCount, availableSpace)
	}

	// Add items to the room
	err = s.AddItemsToRoom(room, itemConfigs)
	if err != nil {
		return nil, err
	}

	// Add guardian monsters if specified
	if len(guardianMonsterConfigs) > 0 {
		err = s.AddMonstersToRoom(room, guardianMonsterConfigs)
		if err != nil {
			return nil, err
		}
	}

	// Add party members to the room
	// Create player configs from party members
	playerConfigs := make([]PlayerConfig, len(party.Members))
	for i, member := range party.Members {
		playerConfigs[i] = PlayerConfig{
			Name:        member.Name,
			Level:       member.Level,
			RandomPlace: true,
		}
	}

	err = s.AddPlayersToRoom(room, playerConfigs)
	if err != nil {
		return nil, err
	}

	// Add a note to the room description about clearing after treasure collection
	if room.Description == "" {
		room.Description = "A treasure room with valuable items. Remember to clear the room after collecting all treasure."
	} else {
		room.Description += " Remember to clear the room after collecting all treasure."
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
				// Get XP from the repository
				xp, err := s.monsterRepo.GetMonsterXP(monster.Key)
				if err != nil {
					// Log the error but continue
					fmt.Printf("Warning: failed to get XP for monster %s: %v\n", monster.Key, err)
				} else {
					totalXP += xp
				}
			}

			// Create a copy of monster IDs to avoid modification during iteration
			monsterIDs := make([]string, len(room.Monsters))
			for i, monster := range room.Monsters {
				monsterIDs[i] = monster.ID
			}

			// Remove each monster by ID
			for _, id := range monsterIDs {
				if !entities.RemoveEntity(room, id, entities.CellMonster) {
					notRemoved = append(notRemoved, id)
				}
			}
		} else {
			// Remove specific monsters by ID
			for _, monsterID := range entityIDs {
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

				// Use RemoveEntity to handle removal
				if !entities.RemoveEntity(room, monsterID, entities.CellMonster) {
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
				if !entities.RemoveEntity(room, id, entities.CellItem) {
					notRemoved = append(notRemoved, id)
				}
			}
		} else {
			// Remove specific items by ID
			for _, itemID := range entityIDs {
				if !entities.RemoveEntity(room, itemID, entities.CellItem) {
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
				if !entities.RemoveEntity(room, id, entities.CellPlayer) {
					notRemoved = append(notRemoved, id)
				}
			}
		} else {
			// Remove specific players by ID
			for _, playerID := range entityIDs {
				if !entities.RemoveEntity(room, playerID, entities.CellPlayer) {
					notRemoved = append(notRemoved, playerID)
				}
			}
		}

	default:
		return 0, nil, fmt.Errorf("unsupported entity type: %v", entityType)
	}

	return totalXP, notRemoved, nil
}
