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

func TestCleanupRoom(t *testing.T) {

	// Create a service with the mock repository
	service := &RoomService{}

	testCases := []struct {
		name          string
		setupRoom     func() *entities.Room
		monsterIDs    []string
		expectedXP    int
		expectedCount int      // Number of monsters that should remain after cleanup
		notRemovedIDs []string // IDs of monsters that should not be removed
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
			monsterIDs:    []string{}, // Empty means remove all
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
			monsterIDs:    []string{"1", "3"}, // Remove goblin and adultbluedragon
			expectedXP:    15050,              // 50 (goblin) + 15000 (adult blue dragon)
			expectedCount: 1,                  // Only banditcaptain should remain
			notRemovedIDs: []string{},         // All specified monsters should be removed
		},
		{
			name: "Remove non-existent monsters",
			setupRoom: func() *entities.Room {
				room := NewRoom(10, 10, entities.LightLevelBright)
				InitializeGrid(room)

				// Add one monster
				goblin := entities.Monster{ID: "1", Key: "monster_goblin", Name: "Goblin", Position: entities.Position{X: 1, Y: 1}, XP: 50}
				PlaceEntity(room, &goblin)

				return room
			},
			monsterIDs:    []string{"999"}, // Non-existent ID
			expectedXP:    0,               // No monsters removed
			expectedCount: 1,               // Monster should still be there
			notRemovedIDs: []string{"999"}, // This ID should be in the not-removed list
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			room := tc.setupRoom()

			xp, notRemoved, err := service.CleanupRoom(room, entities.CellMonster, tc.monsterIDs)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedXP, xp, "Expected XP doesn't match")
			assert.Equal(t, tc.notRemovedIDs, notRemoved, "Not removed IDs don't match")
			assert.Equal(t, tc.expectedCount, len(room.Monsters), "Unexpected number of monsters remaining")

			if tc.expectedCount > 0 {
				// If we expect monsters to remain, make sure the right ones are there
				if len(tc.monsterIDs) > 0 {
					for _, monster := range room.Monsters {
						for _, id := range tc.monsterIDs {
							assert.NotEqual(t, id, monster.ID, "Monster with ID %s should have been removed", id)
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

	// Verify grid is still nil after entity removal
	assert.Nil(t, room.Grid)
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
}
