package services

import (
	"testing"

	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
	"github.com/fadedpez/dnd5e-roomgen/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestRoomConfig creates a standard room configuration for testing
func createTestRoomConfig(width, height int, lightLevel entities.LightLevel, useGrid bool) RoomConfig {
	return RoomConfig{
		Width:       width,
		Height:      height,
		LightLevel:  lightLevel,
		Description: "Test room",
		UseGrid:     useGrid,
	}
}

// createTestRoomConfigWithObstacles creates a room configuration with obstacles for testing
func createTestRoomConfigWithObstacles(width, height int, lightLevel entities.LightLevel, useGrid bool, obstacleConfigs []ObstacleConfig) (RoomConfig, []PlaceableConfig) {
	config := createTestRoomConfig(width, height, lightLevel, useGrid)

	// Convert ObstacleConfigs to PlaceableConfigs
	placeableConfigs := make([]PlaceableConfig, 0, len(obstacleConfigs))
	for _, obstacleConfig := range obstacleConfigs {
		placeableConfigs = append(placeableConfigs, obstacleConfig)
	}

	return config, placeableConfigs
}

// assertRoomProperties checks that a room has the expected properties
func assertRoomProperties(t *testing.T, room *entities.Room, width, height int, lightLevel entities.LightLevel, hasGrid bool) {
	assert.Equal(t, width, room.Width)
	assert.Equal(t, height, room.Height)
	assert.Equal(t, lightLevel, room.LightLevel)

	if hasGrid {
		assert.NotNil(t, room.Grid)
		assert.Equal(t, height, len(room.Grid))
		assert.Equal(t, width, len(room.Grid[0]))
	} else {
		assert.Nil(t, room.Grid)
	}
}

func TestGenerateRoom(t *testing.T) {
	testCases := []struct {
		name        string
		config      RoomConfig
		expectError bool
		checkFunc   func(*testing.T, *entities.Room)
	}{
		{
			name:        "Valid room with grid",
			config:      createTestRoomConfig(10, 15, entities.LightLevelDim, true),
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				assertRoomProperties(t, room, 10, 15, entities.LightLevelDim, true)
				assert.Equal(t, "Test room", room.Description)
			},
		},
		{
			name:        "Valid room without grid",
			config:      createTestRoomConfig(5, 8, entities.LightLevelDark, false),
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				assertRoomProperties(t, room, 5, 8, entities.LightLevelDark, false)
			},
		},
		{
			name: "Default light level",
			config: RoomConfig{
				Width:       7,
				Height:      7,
				Description: "A room with default lighting",
				UseGrid:     true,
			},
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				assert.Equal(t, entities.LightLevelBright, room.LightLevel)
				assert.NotNil(t, room.Grid)
			},
		},
		{
			name:        "Invalid width",
			config:      createTestRoomConfig(0, 10, entities.LightLevelBright, true),
			expectError: true,
		},
		{
			name:        "Invalid height",
			config:      createTestRoomConfig(10, -5, entities.LightLevelBright, true),
			expectError: true,
		},
	}

	service, err := NewRoomService()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			room, err := service.GenerateRoom(tc.config)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, room)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, room)
				if tc.checkFunc != nil {
					tc.checkFunc(t, room)
				}
			}
		})
	}
}

// createTestMonsterConfigWithRealData creates a monster configuration using real monster data from test files
func createTestMonsterConfigWithRealData(t *testing.T, key string, count int, randomPlace bool, position *entities.Position) MonsterConfig {
	// Load monster data
	monsters, err := testutil.LoadAllMonsters()
	require.NoError(t, err, "Failed to load monster data")

	// Get the specific monster
	monster, ok := monsters[key]
	require.True(t, ok, "Monster %s not found in test data", key)

	return MonsterConfig{
		Name:        monster.Name,
		Key:         key,
		CR:          monster.ChallengeRating,
		Count:       count,
		RandomPlace: randomPlace,
		Position:    position,
	}
}

// createTestItemConfigWithRealData creates an item configuration using real item data from test files
func createTestItemConfigWithRealData(t *testing.T, key string, randomPlace bool, position *entities.Position) ItemConfig {
	items, err := testutil.LoadAllEquipment()
	require.NoError(t, err, "Failed to load item data")

	item, ok := items[key]
	require.True(t, ok, "Item %s not found in test data", key)

	return ItemConfig{
		Key:         key,
		Name:        item.Name, // Set the name from the loaded item data
		RandomPlace: randomPlace,
		Position:    position,
	}
}

// createTestPlayerConfig creates a standard player configuration for testing
func createTestPlayerConfig(name string, level int, randomPlace bool, position *entities.Position) PlayerConfig {
	config := PlayerConfig{
		Name:        name,
		Level:       level,
		RandomPlace: randomPlace,
	}

	if position != nil {
		config.Position = position
	}

	return config
}

// createTestNPCConfig creates a standard NPC configuration for testing
func createTestNPCConfig(name string, level int, count int, randomPlace bool, position *entities.Position, inventory []entities.Item) NPCConfig {
	config := NPCConfig{
		Name:        name,
		Level:       level,
		Count:       count,
		Inventory:   inventory,
		RandomPlace: randomPlace,
	}

	if position != nil {
		config.Position = position
	}

	return config
}

// createTestObstacleConfig creates a standard obstacle configuration for testing
func createTestObstacleConfig(name string, key string, blocking bool, count int, randomPlace bool, position *entities.Position) ObstacleConfig {
	return ObstacleConfig{
		Name:        name,
		Key:         key,
		Blocking:    blocking,
		Count:       count,
		RandomPlace: randomPlace,
		Position:    position,
	}
}

func TestCleanupRoom(t *testing.T) {

	// Create a service with the mock repository
	service := &RoomService{}

	testCases := []struct {
		name          string
		setupRoom     func() *entities.Room
		entityType    entities.CellType
		entityIDs     []string
		expectedXP    int
		expectedCount int      // Number of entities that should remain after cleanup
		notRemovedIDs []string // IDs of entities that should not be removed
	}{
		{
			name: "Remove all monsters",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)

				// Add three different monsters
				goblin := entities.Monster{ID: "1", Key: "monster_goblin", Name: "Goblin", Position: entities.Position{X: 1, Y: 1}, XP: 50}
				banditcaptain := entities.Monster{ID: "2", Key: "monster_bandit-captain", Name: "Bandit Captain", Position: entities.Position{X: 3, Y: 3}, XP: 450}
				adultbluedragon := entities.Monster{ID: "3", Key: "monster_adult-blue-dragon", Name: "Adult Blue Dragon", Position: entities.Position{X: 5, Y: 5}, XP: 15000}

				PlaceEntity(room, &goblin)
				PlaceEntity(room, &banditcaptain)
				PlaceEntity(room, &adultbluedragon)

				return room
			},
			entityType:    entities.CellMonster,
			entityIDs:     []string{}, // Empty means remove all
			expectedXP:    15500,      // 50 (goblin) + 450 (bandit captain) + 15000 (adult blue dragon)
			expectedCount: 0,
			notRemovedIDs: []string{}, // All should be removed
		},
		{
			name: "Remove specific monsters",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)

				// Add three different monsters
				goblin := entities.Monster{ID: "1", Key: "monster_goblin", Name: "Goblin", Position: entities.Position{X: 1, Y: 1}, XP: 50}
				banditcaptain := entities.Monster{ID: "2", Key: "monster_bandit-captain", Name: "Bandit Captain", Position: entities.Position{X: 3, Y: 3}, XP: 450}
				adultbluedragon := entities.Monster{ID: "3", Key: "monster_adult-blue-dragon", Name: "Adult Blue Dragon", Position: entities.Position{X: 5, Y: 5}, XP: 15000}

				PlaceEntity(room, &goblin)
				PlaceEntity(room, &banditcaptain)
				PlaceEntity(room, &adultbluedragon)

				return room
			},
			entityType:    entities.CellMonster,
			entityIDs:     []string{"1", "3"}, // Remove goblin and dragon
			expectedXP:    15050,              // 50 (goblin) + 15000 (adult blue dragon)
			expectedCount: 1,                  // Bandit captain should remain
			notRemovedIDs: []string{},         // All requested monsters should be removed
		},
		{
			name: "Remove non-existent monster",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)

				// Add one monster
				banditcaptain := entities.Monster{ID: "2", Key: "monster_bandit-captain", Name: "Bandit Captain", Position: entities.Position{X: 3, Y: 3}, XP: 450}
				PlaceEntity(room, &banditcaptain)

				return room
			},
			entityType:    entities.CellMonster,
			entityIDs:     []string{"999"}, // Non-existent ID
			expectedXP:    0,               // No monsters removed
			expectedCount: 1,               // Monster should still be there
			notRemovedIDs: []string{"999"}, // This ID should be in the not-removed list
		},
		{
			name: "Remove all players",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)

				// Add two players
				player1 := entities.Player{ID: "p1", Name: "Aragorn", Level: 5, Position: entities.Position{X: 1, Y: 1}}
				player2 := entities.Player{ID: "p2", Name: "Gandalf", Level: 10, Position: entities.Position{X: 3, Y: 3}}

				PlaceEntity(room, &player1)
				PlaceEntity(room, &player2)

				return room
			},
			entityType:    entities.CellPlayer,
			entityIDs:     []string{}, // Empty means remove all
			expectedXP:    0,          // Players don't give XP
			expectedCount: 0,
			notRemovedIDs: []string{}, // All should be removed
		},
		{
			name: "Remove all items",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)

				// Add two items
				item1 := entities.Item{ID: "i1", Key: "equipment_potion-of-healing", Name: "Potion of Healing", Position: entities.Position{X: 1, Y: 1}}
				item2 := entities.Item{ID: "i2", Key: "equipment_longsword", Name: "Longsword", Position: entities.Position{X: 3, Y: 3}}

				PlaceEntity(room, &item1)
				PlaceEntity(room, &item2)

				return room
			},
			entityType:    entities.CellItem,
			entityIDs:     []string{}, // Empty means remove all
			expectedXP:    0,          // Items don't give XP
			expectedCount: 0,
			notRemovedIDs: []string{}, // All should be removed
		},
		{
			name: "Remove all NPCs",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)

				// Add two NPCs
				npc1 := entities.NPC{ID: "n1", Name: "Merchant", Position: entities.Position{X: 1, Y: 1}}
				npc2 := entities.NPC{ID: "n2", Name: "Guard", Position: entities.Position{X: 3, Y: 3}}

				PlaceEntity(room, &npc1)
				PlaceEntity(room, &npc2)

				return room
			},
			entityType:    entities.CellNPC,
			entityIDs:     []string{}, // Empty means remove all
			expectedXP:    0,          // NPCs don't give XP
			expectedCount: 0,
			notRemovedIDs: []string{}, // All should be removed
		},
		{
			name: "Remove specific NPCs",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)

				// Add two NPCs
				npc1 := entities.NPC{ID: "n1", Name: "Merchant", Position: entities.Position{X: 1, Y: 1}}
				npc2 := entities.NPC{ID: "n2", Name: "Guard", Position: entities.Position{X: 3, Y: 3}}

				PlaceEntity(room, &npc1)
				PlaceEntity(room, &npc2)

				return room
			},
			entityType:    entities.CellNPC,
			entityIDs:     []string{"n1"}, // Remove merchant
			expectedXP:    0,              // NPCs don't give XP
			expectedCount: 1,              // Guard should remain
			notRemovedIDs: []string{},     // All requested NPCs should be removed
		},
		{
			name: "Remove all obstacles",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)

				// Add two obstacles
				wall := entities.Obstacle{ID: "o1", Name: "Stone Wall", Key: "wall_stone", Blocking: true, Position: entities.Position{X: 1, Y: 1}}
				table := entities.Obstacle{ID: "o2", Name: "Wooden Table", Key: "furniture_table", Blocking: true, Position: entities.Position{X: 3, Y: 3}}

				PlaceEntity(room, &wall)
				PlaceEntity(room, &table)

				return room
			},
			entityType:    entities.CellObstacle,
			entityIDs:     []string{}, // Empty means remove all
			expectedXP:    0,          // Obstacles don't give XP
			expectedCount: 0,
			notRemovedIDs: []string{}, // All should be removed
		},
		{
			name: "Remove specific obstacles",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)

				// Add three obstacles
				wall := entities.Obstacle{ID: "o1", Name: "Stone Wall", Key: "wall_stone", Blocking: true, Position: entities.Position{X: 1, Y: 1}}
				table := entities.Obstacle{ID: "o2", Name: "Wooden Table", Key: "furniture_table", Blocking: true, Position: entities.Position{X: 3, Y: 3}}
				barrel := entities.Obstacle{ID: "o3", Name: "Barrel", Key: "furniture_barrel", Blocking: false, Position: entities.Position{X: 5, Y: 5}}

				PlaceEntity(room, &wall)
				PlaceEntity(room, &table)
				PlaceEntity(room, &barrel)

				return room
			},
			entityType:    entities.CellObstacle,
			entityIDs:     []string{"o1", "o3"}, // Remove wall and barrel
			expectedXP:    0,                    // Obstacles don't give XP
			expectedCount: 1,                    // Table should remain
			notRemovedIDs: []string{},           // All requested obstacles should be removed
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			room := tc.setupRoom()

			xp, notRemoved, err := service.CleanupRoom(room, tc.entityType, tc.entityIDs)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedXP, xp, "Expected XP doesn't match")
			assert.Equal(t, tc.notRemovedIDs, notRemoved, "Not removed IDs don't match")

			// Check the count of remaining entities based on entity type
			switch tc.entityType {
			case entities.CellMonster:
				assert.Equal(t, tc.expectedCount, len(room.Monsters), "Unexpected number of monsters remaining")
			case entities.CellPlayer:
				assert.Equal(t, tc.expectedCount, len(room.Players), "Unexpected number of players remaining")
			case entities.CellItem:
				assert.Equal(t, tc.expectedCount, len(room.Items), "Unexpected number of items remaining")
			case entities.CellNPC:
				assert.Equal(t, tc.expectedCount, len(room.NPCs), "Unexpected number of NPCs remaining")
			case entities.CellObstacle:
				assert.Equal(t, tc.expectedCount, len(room.Obstacles), "Unexpected number of obstacles remaining")
			}

			if tc.expectedCount > 0 && len(tc.entityIDs) > 0 {
				// If we expect entities to remain, make sure the right ones are there
				switch tc.entityType {
				case entities.CellMonster:
					for _, monster := range room.Monsters {
						for _, id := range tc.entityIDs {
							assert.NotEqual(t, id, monster.ID, "Monster with ID %s should have been removed", id)
						}
					}
				case entities.CellPlayer:
					for _, player := range room.Players {
						for _, id := range tc.entityIDs {
							assert.NotEqual(t, id, player.ID, "Player with ID %s should have been removed", id)
						}
					}
				case entities.CellItem:
					for _, item := range room.Items {
						for _, id := range tc.entityIDs {
							assert.NotEqual(t, id, item.ID, "Item with ID %s should have been removed", id)
						}
					}
				case entities.CellNPC:
					for _, npc := range room.NPCs {
						for _, id := range tc.entityIDs {
							assert.NotEqual(t, id, npc.ID, "NPC with ID %s should have been removed", id)
						}
					}
				case entities.CellObstacle:
					for _, obstacle := range room.Obstacles {
						for _, id := range tc.entityIDs {
							assert.NotEqual(t, id, obstacle.ID, "Obstacle with ID %s should have been removed", id)
						}
					}
				}
			}
		})
	}
}

func TestGridlessRoomEntityPlacement(t *testing.T) {
	// TODO: Update this test to use real monster data from the API instead of mock names
	// Currently, this test will produce 404 warnings when calculating XP because
	// it uses fictional monster names that don't exist in the API.
	// These warnings are expected and don't affect the test functionality.

	service, err := NewRoomService()
	if err != nil {
		t.Fatal(err)
	}

	// Create a gridless room
	roomConfig := createTestRoomConfig(10, 10, entities.LightLevelBright, false)
	room, err := service.GenerateRoom(roomConfig)
	if err != nil {
		t.Fatal(err)
	}

	// Verify room is gridless
	assert.Nil(t, room.Grid)
	assert.Equal(t, 10, room.Width)
	assert.Equal(t, 10, room.Height)

	// Test monster placement
	monsterConfigs := []MonsterConfig{
		createTestMonsterConfigWithRealData(t, "goblin", 1, true, nil),
		createTestMonsterConfigWithRealData(t, "goblin", 1, true, nil),
		createTestMonsterConfigWithRealData(t, "goblin", 1, true, nil),
		createTestMonsterConfigWithRealData(t, "bandit-captain", 1, false, &entities.Position{X: 5, Y: 5}),
		createTestMonsterConfigWithRealData(t, "bandit-captain", 1, false, &entities.Position{X: 5, Y: 5}),
	}

	// Convert MonsterConfigs to PlaceableConfigs
	placeableConfigs := make([]PlaceableConfig, 0)
	for _, config := range monsterConfigs {
		placeableConfigs = append(placeableConfigs, config)
	}

	err = service.AddPlaceablesToRoom(room, placeableConfigs)
	assert.NoError(t, err)
	assert.Len(t, room.Monsters, 5)

	// Verify specific positions for fixed-position monsters
	banditcaptainCount := 0
	for _, monster := range room.Monsters {
		if monster.Name == "Bandit Captain" {
			banditcaptainCount++
			assert.Equal(t, 5, monster.Position.X)
			assert.Equal(t, 5, monster.Position.Y)
		}
	}
	assert.Equal(t, 2, banditcaptainCount)

	// Test player placement
	playerConfigs := []PlayerConfig{
		createTestPlayerConfig("Aragorn", 5, true, nil),
		createTestPlayerConfig("Gandalf", 10, false, &entities.Position{X: 3, Y: 3}),
	}

	// Convert PlayerConfigs to PlaceableConfigs
	playerPlaceableConfigs := make([]PlaceableConfig, 0)
	for _, config := range playerConfigs {
		playerPlaceableConfigs = append(playerPlaceableConfigs, config)
	}

	err = service.AddPlaceablesToRoom(room, playerPlaceableConfigs)
	assert.NoError(t, err)
	assert.Len(t, room.Players, 2)

	// Verify specific position for fixed-position player
	var foundGandalf bool
	for _, player := range room.Players {
		if player.Name == "Gandalf" {
			foundGandalf = true
			assert.Equal(t, 3, player.Position.X)
			assert.Equal(t, 3, player.Position.Y)
		}
	}
	assert.True(t, foundGandalf, "Gandalf should be found in the players list")

	// Test item placement
	itemConfigs := []ItemConfig{
		createTestItemConfigWithRealData(t, "abacus", true, nil),
		createTestItemConfigWithRealData(t, "abacus", true, nil),
		createTestItemConfigWithRealData(t, "battleaxe", false, &entities.Position{X: 7, Y: 7}),
	}

	// Convert ItemConfigs to PlaceableConfigs
	itemPlaceableConfigs := make([]PlaceableConfig, 0)
	for _, config := range itemConfigs {
		itemPlaceableConfigs = append(itemPlaceableConfigs, config)
	}

	err = service.AddPlaceablesToRoom(room, itemPlaceableConfigs)
	assert.NoError(t, err)
	assert.Len(t, room.Items, 3)

	// Verify specific position for fixed-position item
	var foundBattleaxe bool
	for _, item := range room.Items {
		if item.Name == "Battleaxe" {
			foundBattleaxe = true
			assert.Equal(t, 7, item.Position.X)
			assert.Equal(t, 7, item.Position.Y)
		}
	}
	assert.True(t, foundBattleaxe, "Battleaxe should be found in the items list")

	// Test NPC placement
	initialItems := []entities.Item{
		{
			ID:       "item1",
			Key:      "equipment_potion-of-healing",
			Name:     "Potion of Healing",
			Position: entities.Position{X: 0, Y: 0},
		},
	}

	npcConfigs := []NPCConfig{
		createTestNPCConfig("Merchant", 3, 1, true, nil, initialItems),
		createTestNPCConfig("Guard", 2, 1, false, &entities.Position{X: 6, Y: 6}, nil),
	}

	// Convert NPCConfigs to PlaceableConfigs
	npcPlaceableConfigs := make([]PlaceableConfig, 0)
	for _, config := range npcConfigs {
		npcPlaceableConfigs = append(npcPlaceableConfigs, config)
	}

	err = service.AddPlaceablesToRoom(room, npcPlaceableConfigs)
	assert.NoError(t, err)
	assert.Len(t, room.NPCs, 2)

	// Verify specific position for fixed-position NPC
	var foundGuard bool
	for _, npc := range room.NPCs {
		if npc.Name == "Guard" {
			foundGuard = true
			assert.Equal(t, 6, npc.Position.X)
			assert.Equal(t, 6, npc.Position.Y)
		}
	}
	assert.True(t, foundGuard, "Guard should be found in the NPCs list")

	// Verify inventory for merchant
	var merchantID string
	for _, npc := range room.NPCs {
		if npc.Name == "Merchant" {
			merchantID = npc.ID
			assert.Len(t, npc.Inventory, 1)
			assert.Equal(t, "Potion of Healing", npc.Inventory[0].Name)
		}
	}
	assert.NotEmpty(t, merchantID, "Merchant should be found in the NPCs list")

	// Test obstacle placement
	obstacleConfigs := []ObstacleConfig{
		createTestObstacleConfig("Stone Wall", "wall_stone", true, 1, true, nil),
		createTestObstacleConfig("Wooden Table", "furniture_table", true, 1, false, &entities.Position{X: 8, Y: 8}),
	}

	// Convert ObstacleConfigs to PlaceableConfigs
	obstaclePlaceableConfigs := make([]PlaceableConfig, 0)
	for _, config := range obstacleConfigs {
		obstaclePlaceableConfigs = append(obstaclePlaceableConfigs, config)
	}

	err = service.AddPlaceablesToRoom(room, obstaclePlaceableConfigs)
	assert.NoError(t, err)
	assert.Len(t, room.Obstacles, 2)

	// Verify specific position for fixed-position obstacle
	var foundTable bool
	for _, obstacle := range room.Obstacles {
		if obstacle.Name == "Wooden Table" {
			foundTable = true
			assert.Equal(t, 8, obstacle.Position.X)
			assert.Equal(t, 8, obstacle.Position.Y)
			assert.True(t, obstacle.Blocking)
		}
	}
	assert.True(t, foundTable, "Wooden Table should be found in the obstacles list")

	// Verify grid is still nil after all entity placements
	assert.Nil(t, room.Grid)

	// Test entity removal
	// Remove all monsters
	_, notRemoved, err := service.CleanupRoom(room, entities.CellMonster, []string{})
	assert.NoError(t, err)
	assert.Empty(t, notRemoved)
	assert.Len(t, room.Monsters, 0)

	// Remove one player
	_, notRemoved, err = service.CleanupRoom(room, entities.CellPlayer, []string{room.Players[0].ID})
	assert.NoError(t, err)
	assert.Empty(t, notRemoved)
	assert.Len(t, room.Players, 1)

	// Remove one item
	_, notRemoved, err = service.CleanupRoom(room, entities.CellItem, []string{room.Items[0].ID})
	assert.NoError(t, err)
	assert.Empty(t, notRemoved)
	assert.Len(t, room.Items, 2)

	// Remove one NPC
	_, notRemoved, err = service.CleanupRoom(room, entities.CellNPC, []string{merchantID})
	assert.NoError(t, err)
	assert.Empty(t, notRemoved)
	assert.Len(t, room.NPCs, 1)
	assert.Equal(t, "Guard", room.NPCs[0].Name)

	// Remove one obstacle
	var tableID string
	for _, obstacle := range room.Obstacles {
		if obstacle.Name == "Wooden Table" {
			tableID = obstacle.ID
			break
		}
	}
	assert.NotEmpty(t, tableID, "Should have found the table ID")
	_, notRemoved, err = service.CleanupRoom(room, entities.CellObstacle, []string{tableID})
	assert.NoError(t, err)
	assert.Empty(t, notRemoved)
	assert.Len(t, room.Obstacles, 1)
	assert.NotEqual(t, "Wooden Table", room.Obstacles[0].Name)

	// Verify grid is still nil after entity removal
	assert.Nil(t, room.Grid)
}

func TestPlacementOnGrid(t *testing.T) {

	// Define test cases for different entity types
	testCases := []struct {
		name          string
		setupRoom     func() *entities.Room
		entity        entities.Placeable
		expectedError bool
		checkFunc     func(t *testing.T, room *entities.Room, entity entities.Placeable)
	}{
		{
			name: "Place player on empty cell",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)
				return room
			},
			entity: &entities.Player{
				ID:       "p1",
				Name:     "Aragorn",
				Level:    5,
				Position: entities.Position{X: 2, Y: 3},
			},
			expectedError: false,
			checkFunc: func(t *testing.T, room *entities.Room, entity entities.Placeable) {
				player := entity.(*entities.Player)

				// Verify player was added
				assert.Equal(t, 1, len(room.Players))
				assert.Equal(t, player.ID, room.Players[0].ID)
				assert.Equal(t, player.Position.X, room.Players[0].Position.X)
				assert.Equal(t, player.Position.Y, room.Players[0].Position.Y)

				// Verify grid cell is updated
				assert.Equal(t, entities.CellPlayer, room.Grid[player.Position.Y][player.Position.X].Type)
				assert.Equal(t, player.ID, room.Grid[player.Position.Y][player.Position.X].EntityID)
			},
		},
		{
			name: "Place monster on empty cell",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)
				return room
			},
			entity: &entities.Monster{
				ID:       "m1",
				Name:     "Goblin",
				CR:       0.25,
				Position: entities.Position{X: 4, Y: 5},
			},
			expectedError: false,
			checkFunc: func(t *testing.T, room *entities.Room, entity entities.Placeable) {
				monster := entity.(*entities.Monster)

				// Verify monster was added
				assert.Equal(t, 1, len(room.Monsters))
				assert.Equal(t, monster.ID, room.Monsters[0].ID)
				assert.Equal(t, monster.Position.X, room.Monsters[0].Position.X)
				assert.Equal(t, monster.Position.Y, room.Monsters[0].Position.Y)

				// Verify grid cell is updated
				assert.Equal(t, entities.CellMonster, room.Grid[monster.Position.Y][monster.Position.X].Type)
				assert.Equal(t, monster.ID, room.Grid[monster.Position.Y][monster.Position.X].EntityID)
			},
		},
		{
			name: "Place NPC on empty cell",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)
				return room
			},
			entity: &entities.NPC{
				ID:       "n1",
				Name:     "Merchant",
				Position: entities.Position{X: 6, Y: 7},
			},
			expectedError: false,
			checkFunc: func(t *testing.T, room *entities.Room, entity entities.Placeable) {
				npc := entity.(*entities.NPC)

				// Verify NPC was added
				assert.Equal(t, 1, len(room.NPCs))
				assert.Equal(t, npc.ID, room.NPCs[0].ID)
				assert.Equal(t, npc.Position.X, room.NPCs[0].Position.X)
				assert.Equal(t, npc.Position.Y, room.NPCs[0].Position.Y)

				// Verify grid cell is updated
				assert.Equal(t, entities.CellNPC, room.Grid[npc.Position.Y][npc.Position.X].Type)
				assert.Equal(t, npc.ID, room.Grid[npc.Position.Y][npc.Position.X].EntityID)
			},
		},
		{
			name: "Place item on empty cell",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)
				return room
			},
			entity: &entities.Item{
				ID:       "i1",
				Name:     "Potion of Healing",
				Key:      "potion_healing",
				Position: entities.Position{X: 8, Y: 1},
			},
			expectedError: false,
			checkFunc: func(t *testing.T, room *entities.Room, entity entities.Placeable) {
				item := entity.(*entities.Item)

				// Verify item was added
				assert.Equal(t, 1, len(room.Items))
				assert.Equal(t, item.ID, room.Items[0].ID)
				assert.Equal(t, item.Position.X, room.Items[0].Position.X)
				assert.Equal(t, item.Position.Y, room.Items[0].Position.Y)

				// Verify grid cell is updated
				assert.Equal(t, entities.CellItem, room.Grid[item.Position.Y][item.Position.X].Type)
				assert.Equal(t, item.ID, room.Grid[item.Position.Y][item.Position.X].EntityID)
			},
		},
		{
			name: "Place obstacle on empty cell",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)
				return room
			},
			entity: &entities.Obstacle{
				ID:       "o1",
				Name:     "Stone Wall",
				Key:      "wall_stone",
				Blocking: true,
				Position: entities.Position{X: 2, Y: 3},
			},
			expectedError: false,
			checkFunc: func(t *testing.T, room *entities.Room, entity entities.Placeable) {
				obstacle := entity.(*entities.Obstacle)

				// Verify obstacle was added
				assert.Equal(t, 1, len(room.Obstacles))
				assert.Equal(t, obstacle.ID, room.Obstacles[0].ID)
				assert.Equal(t, obstacle.Position.X, room.Obstacles[0].Position.X)
				assert.Equal(t, obstacle.Position.Y, room.Obstacles[0].Position.Y)

				// Verify grid cell is updated
				assert.Equal(t, entities.CellObstacle, room.Grid[obstacle.Position.Y][obstacle.Position.X].Type)
				assert.Equal(t, obstacle.ID, room.Grid[obstacle.Position.Y][obstacle.Position.X].EntityID)
			},
		},
		{
			name: "Place entity out of bounds",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)
				return room
			},
			entity: &entities.Player{
				ID:       "p2",
				Name:     "Legolas",
				Level:    4,
				Position: entities.Position{X: 15, Y: 3}, // Out of bounds X
			},
			expectedError: true,
			checkFunc: func(t *testing.T, room *entities.Room, entity entities.Placeable) {
				// Verify player was not added
				assert.Equal(t, 0, len(room.Players))
			},
		},
		{
			name: "Place entity on negative coordinates",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)
				return room
			},
			entity: &entities.Monster{
				ID:       "m2",
				Name:     "Orc",
				CR:       0.5,
				Position: entities.Position{X: -1, Y: 3}, // Negative X
			},
			expectedError: true,
			checkFunc: func(t *testing.T, room *entities.Room, entity entities.Placeable) {
				// Verify monster was not added
				assert.Equal(t, 0, len(room.Monsters))
			},
		},
		{
			name: "Place entity on occupied cell",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)

				// Pre-place a player
				player := &entities.Player{
					ID:       "p3",
					Name:     "Gimli",
					Level:    4,
					Position: entities.Position{X: 5, Y: 5},
				}
				PlaceEntity(room, player)

				return room
			},
			entity: &entities.Item{
				ID:       "i2",
				Name:     "Sword",
				Key:      "weapon_sword",
				Position: entities.Position{X: 5, Y: 5}, // Same position as pre-placed player
			},
			expectedError: true,
			checkFunc: func(t *testing.T, room *entities.Room, entity entities.Placeable) {
				// Verify item was not added
				assert.Equal(t, 0, len(room.Items))

				// Verify player is still there
				assert.Equal(t, 1, len(room.Players))
				assert.Equal(t, "p3", room.Players[0].ID)
			},
		},
		{
			name: "Place entity in gridless room",
			setupRoom: func() *entities.Room {
				// Create room without grid
				return NewRoom(10, 10, entities.LightLevelBright)
			},
			entity: &entities.Obstacle{
				ID:       "o2",
				Name:     "Wooden Table",
				Key:      "furniture_table",
				Blocking: true,
				Position: entities.Position{X: 7, Y: 8},
			},
			expectedError: false,
			checkFunc: func(t *testing.T, room *entities.Room, entity entities.Placeable) {
				obstacle := entity.(*entities.Obstacle)

				// Verify obstacle was added
				assert.Equal(t, 1, len(room.Obstacles))
				assert.Equal(t, obstacle.ID, room.Obstacles[0].ID)
				assert.Equal(t, obstacle.Position.X, room.Obstacles[0].Position.X)
				assert.Equal(t, obstacle.Position.Y, room.Obstacles[0].Position.Y)

				// Verify grid is nil
				assert.Nil(t, room.Grid)
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			room := tc.setupRoom()

			// Attempt to place the entity
			err := PlaceEntity(room, tc.entity)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tc.checkFunc(t, room, tc.entity)
			}
		})
	}
}

func TestGridlessRoomCleanup(t *testing.T) {
	// Create a service
	service := &RoomService{}

	// Create a gridless room with monsters
	room := NewRoom(10, 10, entities.LightLevelBright)
	// Explicitly not initializing grid
	assert.Nil(t, room.Grid)

	// Add monsters directly using the placement interface
	goblin := entities.Monster{ID: "1", Key: "monster_goblin", Name: "Goblin", Position: entities.Position{X: 1, Y: 1}, XP: 50}
	banditcaptain := entities.Monster{ID: "2", Key: "monster_bandit-captain", Name: "Bandit Captain", Position: entities.Position{X: 3, Y: 3}, XP: 450}
	adultbluedragon := entities.Monster{ID: "3", Key: "monster_adult-blue-dragon", Name: "Adult Blue Dragon", Position: entities.Position{X: 5, Y: 5}, XP: 15000}

	PlaceEntity(room, &goblin)
	PlaceEntity(room, &banditcaptain)
	PlaceEntity(room, &adultbluedragon)

	// Verify monsters were added
	assert.Len(t, room.Monsters, 3)
	assert.Nil(t, room.Grid)

	// Test removing specific monsters
	var notRemoved []string
	xp, notRemoved, err := service.CleanupRoom(room, entities.CellMonster, []string{"1", "3"})
	assert.NoError(t, err)
	assert.Empty(t, notRemoved)
	assert.Equal(t, 15050, xp) // 50 (goblin) + 15000 (adult blue dragon)
	assert.Len(t, room.Monsters, 1)

	// Verify the remaining monster is the banditcaptain
	assert.Equal(t, "2", room.Monsters[0].ID)
	assert.Equal(t, "Bandit Captain", room.Monsters[0].Name)

	// Verify grid is still nil after cleanup
	assert.Nil(t, room.Grid)

	// Add obstacles directly using the placement interface
	wall := entities.Obstacle{ID: "o1", Name: "Stone Wall", Key: "wall_stone", Blocking: true, Position: entities.Position{X: 1, Y: 1}}
	table := entities.Obstacle{ID: "o2", Name: "Wooden Table", Key: "furniture_table", Blocking: true, Position: entities.Position{X: 3, Y: 3}}
	barrel := entities.Obstacle{ID: "o3", Name: "Barrel", Key: "furniture_barrel", Blocking: false, Position: entities.Position{X: 5, Y: 5}}

	PlaceEntity(room, &wall)
	PlaceEntity(room, &table)
	PlaceEntity(room, &barrel)

	// Verify obstacles were added
	assert.Len(t, room.Obstacles, 3)
	assert.Nil(t, room.Grid)

	// Test removing specific obstacles
	xp, notRemoved, err = service.CleanupRoom(room, entities.CellObstacle, []string{"o1", "o3"})
	assert.NoError(t, err)
	assert.Empty(t, notRemoved)
	assert.Equal(t, 0, xp) // Obstacles don't give XP
	assert.Len(t, room.Obstacles, 1)

	// Verify the remaining obstacle is the table
	assert.Equal(t, "o2", room.Obstacles[0].ID)
	assert.Equal(t, "Wooden Table", room.Obstacles[0].Name)

	// Verify grid is still nil after obstacle cleanup
	assert.Nil(t, room.Grid)
}

func TestNPCInventoryManagement(t *testing.T) {
	service, err := NewRoomService()
	require.NoError(t, err)

	// Create a room
	roomConfig := createTestRoomConfig(10, 10, entities.LightLevelBright, true)
	room, err := service.GenerateRoom(roomConfig)
	require.NoError(t, err)

	// Create initial inventory
	initialItems := []entities.Item{
		{
			ID:       "item1",
			Key:      "equipment_potion-of-healing",
			Name:     "Potion of Healing",
			Position: entities.Position{X: 0, Y: 0},
		},
	}

	// Create an NPC with initial inventory
	npcConfig := createTestNPCConfig("Merchant", 3, 1, false, &entities.Position{X: 5, Y: 5}, initialItems)

	// Add the NPC to the room
	err = service.AddPlaceablesToRoom(room, []PlaceableConfig{npcConfig})
	require.NoError(t, err)

	// Get the NPC ID
	require.Equal(t, 1, len(room.NPCs))
	npcID := room.NPCs[0].ID

	// Test GetNPCInventory
	inventory, err := service.GetNPCInventory(room, npcID)
	require.NoError(t, err)
	require.Equal(t, 1, len(inventory))
	assert.Equal(t, "Potion of Healing", inventory[0].Name)

	// Test AddItemToNPCInventory
	newItem := entities.Item{
		Key:      "equipment_longsword",
		Name:     "Longsword",
		Position: entities.Position{X: 0, Y: 0},
	}
	err = service.AddItemToNPCInventory(room, npcID, newItem)
	require.NoError(t, err)

	// Verify item was added
	inventory, err = service.GetNPCInventory(room, npcID)
	require.NoError(t, err)
	require.Equal(t, 2, len(inventory))

	// Find the item we just added (it will have a new ID)
	var foundNewItem bool
	for _, item := range inventory {
		if item.Name == "Longsword" {
			foundNewItem = true
			break
		}
	}
	assert.True(t, foundNewItem, "Added item not found in inventory")

	// Test RemoveItemFromNPCInventory
	// First get the ID of the potion
	var potionID string
	for _, item := range inventory {
		if item.Name == "Potion of Healing" {
			potionID = item.ID
			break
		}
	}
	require.NotEmpty(t, potionID, "Could not find potion in inventory")

	// Remove the potion
	removedItem, err := service.RemoveItemFromNPCInventory(room, npcID, potionID)
	require.NoError(t, err)
	assert.Equal(t, "Potion of Healing", removedItem.Name)

	// Verify item was removed
	inventory, err = service.GetNPCInventory(room, npcID)
	require.NoError(t, err)
	require.Equal(t, 1, len(inventory))
	assert.Equal(t, "Longsword", inventory[0].Name)

	// Test error cases
	// Try to remove an item that doesn't exist
	_, err = service.RemoveItemFromNPCInventory(room, npcID, "nonexistent-item")
	assert.Error(t, err)

	// Try to access inventory of an NPC that doesn't exist
	_, err = service.GetNPCInventory(room, "nonexistent-npc")
	assert.Error(t, err)
}

func TestEntityPlacementPriority(t *testing.T) {
	service := &RoomService{}

	// Create a room with a grid
	room := NewRoom(10, 10, entities.LightLevelBright)
	InitializeGrid(room)

	// Create configs for different entity types
	// We'll create one of each type and all will try to place at the same position
	targetPos := &entities.Position{X: 2, Y: 2}

	// Create configs in reverse priority order (items, obstacles, NPCs, monsters, players)
	// This tests that the priority ordering in AddPlaceablesToRoom works correctly
	itemConfig := createTestItemConfigWithRealData(t, "abacus", false, targetPos)
	obstacleConfig := createTestObstacleConfig("Stone Wall", "wall_stone", true, 1, false, targetPos)
	npcConfig := createTestNPCConfig("Merchant", 3, 1, false, targetPos, nil)
	monsterConfig := createTestMonsterConfigWithRealData(t, "goblin", 1, false, targetPos)
	playerConfig := createTestPlayerConfig("Aragorn", 5, false, targetPos)

	// Combine all configs into a single slice
	placeableConfigs := []PlaceableConfig{
		itemConfig,
		obstacleConfig,
		npcConfig,
		monsterConfig,
		playerConfig,
	}

	// Add all entities to the room
	err := service.AddPlaceablesToRoom(room, placeableConfigs)
	assert.NoError(t, err)

	// Check if the player was placed (highest priority)
	assert.Equal(t, 1, len(room.Players))
	assert.Equal(t, 2, room.Players[0].Position.X)
	assert.Equal(t, 2, room.Players[0].Position.Y)

	// Check that the grid cell has the player
	assert.Equal(t, entities.CellPlayer, room.Grid[2][2].Type)

	// Now let's check other entities - they should have been placed elsewhere or discarded
	// We'll check each entity type and verify that if it was placed, it's not at the target position

	if len(room.Monsters) > 0 {
		assert.False(t, (room.Monsters[0].Position.X == 2 && room.Monsters[0].Position.Y == 2))
		assert.Equal(t, entities.CellMonster, room.Grid[room.Monsters[0].Position.Y][room.Monsters[0].Position.X].Type)
	}

	if len(room.NPCs) > 0 {
		assert.False(t, (room.NPCs[0].Position.X == 2 && room.NPCs[0].Position.Y == 2))
		assert.Equal(t, entities.CellNPC, room.Grid[room.NPCs[0].Position.Y][room.NPCs[0].Position.X].Type)
	}

	if len(room.Obstacles) > 0 {
		assert.False(t, (room.Obstacles[0].Position.X == 2 && room.Obstacles[0].Position.Y == 2))
		assert.Equal(t, entities.CellObstacle, room.Grid[room.Obstacles[0].Position.Y][room.Obstacles[0].Position.X].Type)
	}

	if len(room.Items) > 0 {
		assert.False(t, (room.Items[0].Position.X == 2 && room.Items[0].Position.Y == 2))
		assert.Equal(t, entities.CellItem, room.Grid[room.Items[0].Position.Y][room.Items[0].Position.X].Type)
	}

	// Verify all entities that were placed are within bounds
	for _, monster := range room.Monsters {
		assert.GreaterOrEqual(t, monster.Position.X, 0)
		assert.Less(t, monster.Position.X, 5)
		assert.GreaterOrEqual(t, monster.Position.Y, 0)
		assert.Less(t, monster.Position.Y, 5)
	}

	for _, npc := range room.NPCs {
		assert.GreaterOrEqual(t, npc.Position.X, 0)
		assert.Less(t, npc.Position.X, 5)
		assert.GreaterOrEqual(t, npc.Position.Y, 0)
		assert.Less(t, npc.Position.Y, 5)
	}

	for _, obstacle := range room.Obstacles {
		assert.GreaterOrEqual(t, obstacle.Position.X, 0)
		assert.Less(t, obstacle.Position.X, 5)
		assert.GreaterOrEqual(t, obstacle.Position.Y, 0)
		assert.Less(t, obstacle.Position.Y, 5)
	}

	for _, item := range room.Items {
		assert.GreaterOrEqual(t, item.Position.X, 0)
		assert.Less(t, item.Position.X, 5)
		assert.GreaterOrEqual(t, item.Position.Y, 0)
		assert.Less(t, item.Position.Y, 5)
	}
}
