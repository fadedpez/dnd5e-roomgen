package services

import (
	"fmt"
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
func createTestItemConfigWithRealData(t *testing.T, key string, count int, randomPlace bool, position *entities.Position) ItemConfig {
	items, err := testutil.LoadAllEquipment()
	require.NoError(t, err, "Failed to load item data")

	_, ok := items[key]
	require.True(t, ok, "Item %s not found in test data", key)

	return ItemConfig{
		Key:         key,
		Count:       count,
		RandomPlace: randomPlace,
		Position:    position,
	}
}

func TestAddMonstersToRoom(t *testing.T) {
	testCases := []struct {
		name           string
		roomSetup      func() *entities.Room
		monsterConfigs []MonsterConfig
		expectError    bool
		checkFunc      func(*testing.T, *entities.Room)
	}{
		{
			name: "Single monster with random placement",
			roomSetup: func() *entities.Room {
				room := entities.NewRoom(10, 10, entities.LightLevelBright)
				entities.InitializeGrid(room)
				return room
			},
			monsterConfigs: []MonsterConfig{
				createTestMonsterConfigWithRealData(t, "goblin", 1, true, nil),
			},
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				assert.Len(t, room.Monsters, 1)
				assert.Equal(t, "Goblin", room.Monsters[0].Name)
				assert.Equal(t, 0.25, room.Monsters[0].CR)
			},
		},
		{
			name: "Multiple monsters with random placement",
			roomSetup: func() *entities.Room {
				room := entities.NewRoom(10, 10, entities.LightLevelBright)
				entities.InitializeGrid(room)
				return room
			},
			monsterConfigs: []MonsterConfig{
				createTestMonsterConfigWithRealData(t, "goblin", 3, true, nil),
			},
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				assert.Len(t, room.Monsters, 3)
				for _, monster := range room.Monsters {
					assert.Equal(t, "Goblin", monster.Name)
				}
			},
		},
		{
			name: "Monster with specific position",
			roomSetup: func() *entities.Room {
				room := entities.NewRoom(10, 10, entities.LightLevelBright)
				entities.InitializeGrid(room)
				return room
			},
			monsterConfigs: []MonsterConfig{
				createTestMonsterConfigWithRealData(t, "bandit-captain", 1, false, &entities.Position{X: 3, Y: 4}),
			},
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				assert.Len(t, room.Monsters, 1)
				assert.Equal(t, "Bandit Captain", room.Monsters[0].Name)
				assert.Equal(t, 3, room.Monsters[0].Position.X)
				assert.Equal(t, 4, room.Monsters[0].Position.Y)
			},
		},
		{
			name: "Room without initialized grid",
			roomSetup: func() *entities.Room {
				// Create room without grid
				return entities.NewRoom(10, 10, entities.LightLevelBright)
			},
			monsterConfigs: []MonsterConfig{
				createTestMonsterConfigWithRealData(t, "goblin", 1, true, nil),
			},
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				// Grid should remain nil as we no longer force initialization
				assert.Nil(t, room.Grid)
				assert.Len(t, room.Monsters, 1)
				// Monster should still have a position
				assert.GreaterOrEqual(t, room.Monsters[0].Position.X, 0)
				assert.Less(t, room.Monsters[0].Position.X, room.Width)
				assert.GreaterOrEqual(t, room.Monsters[0].Position.Y, 0)
				assert.Less(t, room.Monsters[0].Position.Y, room.Height)
			},
		},
		{
			name:      "Nil room",
			roomSetup: func() *entities.Room { return nil },
			monsterConfigs: []MonsterConfig{
				createTestMonsterConfigWithRealData(t, "goblin", 1, true, nil),
			},
			expectError: true,
		},
		{
			name: "Non-random placement without position",
			roomSetup: func() *entities.Room {
				room := entities.NewRoom(10, 10, entities.LightLevelBright)
				entities.InitializeGrid(room)
				return room
			},
			monsterConfigs: []MonsterConfig{
				func() MonsterConfig {
					config := createTestMonsterConfigWithRealData(t, "goblin", 1, false, nil)
					return config
				}(),
			},
			expectError: true,
		},
	}

	service, err := NewRoomService()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			room := tc.roomSetup()
			err := service.AddMonstersToRoom(room, tc.monsterConfigs)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, room)
				}
			}
		})
	}
}

func TestPopulateRoomWithMonsters(t *testing.T) {
	testCases := []struct {
		name           string
		roomConfig     RoomConfig
		monsterConfigs []MonsterConfig
		expectError    bool
		checkFunc      func(*testing.T, *entities.Room)
	}{
		{
			name:       "Valid room with monsters",
			roomConfig: createTestRoomConfig(10, 10, entities.LightLevelBright, true),
			monsterConfigs: []MonsterConfig{
				createTestMonsterConfigWithRealData(t, "goblin", 2, true, nil),
			},
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				// Check room properties
				assertRoomProperties(t, room, 10, 10, entities.LightLevelBright, true)
				// Check monsters
				assert.Len(t, room.Monsters, 2)
				assert.Equal(t, "Goblin", room.Monsters[0].Name)
			},
		},
		{
			name:       "Invalid room dimensions",
			roomConfig: createTestRoomConfig(-5, 10, entities.LightLevelBright, true),
			monsterConfigs: []MonsterConfig{
				createTestMonsterConfigWithRealData(t, "goblin", 1, true, nil),
			},
			expectError: true,
		},
		{
			name:       "Invalid monster placement",
			roomConfig: createTestRoomConfig(10, 10, entities.LightLevelBright, true),
			monsterConfigs: []MonsterConfig{
				func() MonsterConfig {
					config := createTestMonsterConfigWithRealData(t, "goblin", 1, false, nil)
					return config
				}(),
			},
			expectError: true,
		},
	}

	service, err := NewRoomService()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			room, err := service.PopulateRoomWithMonsters(tc.roomConfig, tc.monsterConfigs)

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

func TestCleanupRoom(t *testing.T) {
	// Create a mock monster repository for testing
	mockRepo := NewMockMonsterRepositoryWithTestData(t)

	// Create a service with the mock repository
	service := &RoomService{
		monsterRepo: mockRepo,
	}

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
				room := entities.NewRoom(10, 10, entities.LightLevelBright)
				entities.InitializeGrid(room)

				// Add three different monsters
				goblin := entities.Monster{ID: "1", Key: "monster_goblin", Name: "Goblin", Position: entities.Position{X: 1, Y: 1}}
				banditcaptain := entities.Monster{ID: "2", Key: "monster_bandit-captain", Name: "Bandit Captain", Position: entities.Position{X: 3, Y: 3}}
				adultbluedragon := entities.Monster{ID: "3", Key: "monster_adult-blue-dragon", Name: "Adult Blue Dragon", Position: entities.Position{X: 5, Y: 5}}

				entities.AddMonster(room, goblin)
				entities.AddMonster(room, banditcaptain)
				entities.AddMonster(room, adultbluedragon)

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
				room := entities.NewRoom(10, 10, entities.LightLevelBright)
				entities.InitializeGrid(room)

				// Add three different monsters
				goblin := entities.Monster{ID: "1", Key: "monster_goblin", Name: "Goblin", Position: entities.Position{X: 1, Y: 1}}
				banditcaptain := entities.Monster{ID: "2", Key: "monster_bandit-captain", Name: "Bandit Captain", Position: entities.Position{X: 3, Y: 3}}
				adultbluedragon := entities.Monster{ID: "3", Key: "monster_adult-blue-dragon", Name: "Adult Blue Dragon", Position: entities.Position{X: 5, Y: 5}}

				entities.AddMonster(room, goblin)
				entities.AddMonster(room, banditcaptain)
				entities.AddMonster(room, adultbluedragon)

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
				room := entities.NewRoom(10, 10, entities.LightLevelBright)
				entities.InitializeGrid(room)

				// Add one monster
				goblin := entities.Monster{ID: "1", Key: "monster_goblin", Name: "Goblin", Position: entities.Position{X: 1, Y: 1}}
				entities.AddMonster(room, goblin)

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

func TestAddPlayersToRoom(t *testing.T) {
	testCases := []struct {
		name          string
		roomSetup     func() *entities.Room
		playerConfigs []PlayerConfig
		expectError   bool
		checkFunc     func(*testing.T, *entities.Room)
	}{
		{
			name: "Single player with random placement",
			roomSetup: func() *entities.Room {
				room := entities.NewRoom(10, 10, entities.LightLevelBright)
				entities.InitializeGrid(room)
				return room
			},
			playerConfigs: []PlayerConfig{
				createTestPlayerConfig("Aragorn", 5, true, nil),
			},
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				assert.Len(t, room.Players, 1)
				assert.Equal(t, "Aragorn", room.Players[0].Name)
				assert.Equal(t, 5, room.Players[0].Level)
			},
		},
		{
			name: "Multiple players with specific positions",
			roomSetup: func() *entities.Room {
				room := entities.NewRoom(10, 10, entities.LightLevelBright)
				entities.InitializeGrid(room)
				return room
			},
			playerConfigs: []PlayerConfig{
				createTestPlayerConfig("Gandalf", 10, false, &entities.Position{X: 1, Y: 1}),
				createTestPlayerConfig("Frodo", 3, false, &entities.Position{X: 2, Y: 2}),
			},
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				assert.Len(t, room.Players, 2)

				// Find players by name
				var gandalf, frodo *entities.Player
				for i := range room.Players {
					if room.Players[i].Name == "Gandalf" {
						gandalf = &room.Players[i]
					} else if room.Players[i].Name == "Frodo" {
						frodo = &room.Players[i]
					}
				}

				assert.NotNil(t, gandalf)
				assert.NotNil(t, frodo)

				assert.Equal(t, 10, gandalf.Level)
				assert.Equal(t, 3, frodo.Level)

				assert.Equal(t, 1, gandalf.Position.X)
				assert.Equal(t, 1, gandalf.Position.Y)
				assert.Equal(t, 2, frodo.Position.X)
				assert.Equal(t, 2, frodo.Position.Y)
			},
		},
		{
			name: "Player with invalid position",
			roomSetup: func() *entities.Room {
				room := entities.NewRoom(5, 5, entities.LightLevelBright)
				entities.InitializeGrid(room)
				return room
			},
			playerConfigs: []PlayerConfig{
				createTestPlayerConfig("Legolas", 7, false, &entities.Position{X: 10, Y: 10}),
			},
			expectError: true,
		},
		{
			name: "Player with no position and not random",
			roomSetup: func() *entities.Room {
				room := entities.NewRoom(5, 5, entities.LightLevelBright)
				entities.InitializeGrid(room)
				return room
			},
			playerConfigs: []PlayerConfig{
				func() PlayerConfig {
					config := createTestPlayerConfig("Gimli", 6, false, nil)
					return config
				}(),
			},
			expectError: true,
		},
		{
			name: "Add player to room without grid",
			roomSetup: func() *entities.Room {
				return entities.NewRoom(5, 5, entities.LightLevelBright)
			},
			playerConfigs: []PlayerConfig{
				createTestPlayerConfig("Boromir", 5, true, nil),
			},
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				assert.Len(t, room.Players, 1)
				assert.NotNil(t, room.Grid) // Grid should be initialized
			},
		},
	}

	service, err := NewRoomService()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			room := tc.roomSetup()
			err := service.AddPlayersToRoom(room, tc.playerConfigs)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, room)
				}
			}
		})
	}
}

func TestPopulateRoomWithMonstersAndPlayers(t *testing.T) {
	testCases := []struct {
		name           string
		roomConfig     RoomConfig
		playerConfigs  []PlayerConfig
		monsterConfigs []MonsterConfig
		expectError    bool
		checkFunc      func(*testing.T, *entities.Room)
	}{
		{
			name:       "Room with players and monsters",
			roomConfig: createTestRoomConfig(10, 10, entities.LightLevelBright, true),
			playerConfigs: []PlayerConfig{
				createTestPlayerConfig("Aragorn", 5, true, nil),
				createTestPlayerConfig("Legolas", 5, true, nil),
			},
			monsterConfigs: []MonsterConfig{
				createTestMonsterConfigWithRealData(t, "goblin", 2, true, nil),
				createTestMonsterConfigWithRealData(t, "bandit-captain", 1, false, &entities.Position{X: 5, Y: 5}),
			},
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				// Check room properties
				assertRoomProperties(t, room, 10, 10, entities.LightLevelBright, true)
				// Check players
				assert.Len(t, room.Players, 2)

				// Check monsters
				assert.Len(t, room.Monsters, 3) // 2 goblins + 1 banditcaptain

				// Find the banditcaptain (should be at position 5,5)
				var banditcaptainFound bool
				for _, monster := range room.Monsters {
					if monster.Name == "Bandit Captain" {
						assert.Equal(t, 5, monster.Position.X)
						assert.Equal(t, 5, monster.Position.Y)
						banditcaptainFound = true
						break
					}
				}
				assert.True(t, banditcaptainFound, "Bandit Captain should be found at the specified position")
			},
		},
		{
			name:       "Invalid room config",
			roomConfig: createTestRoomConfig(-5, 10, entities.LightLevelBright, true),
			playerConfigs: []PlayerConfig{
				createTestPlayerConfig("Aragorn", 5, true, nil),
			},
			monsterConfigs: []MonsterConfig{
				createTestMonsterConfigWithRealData(t, "goblin", 1, true, nil),
			},
			expectError: true,
		},
		{
			name:       "Invalid player config",
			roomConfig: createTestRoomConfig(10, 10, entities.LightLevelBright, true),
			playerConfigs: []PlayerConfig{
				createTestPlayerConfig("Aragorn", 5, false, nil), // Missing position
			},
			monsterConfigs: []MonsterConfig{
				createTestMonsterConfigWithRealData(t, "goblin", 1, true, nil),
			},
			expectError: true,
		},
		{
			name:       "Invalid monster config",
			roomConfig: createTestRoomConfig(10, 10, entities.LightLevelBright, true),
			playerConfigs: []PlayerConfig{
				createTestPlayerConfig("Aragorn", 5, true, nil),
			},
			monsterConfigs: []MonsterConfig{
				func() MonsterConfig {
					config := createTestMonsterConfigWithRealData(t, "goblin", 1, false, nil)
					return config
				}(),
			},
			expectError: true,
		},
	}

	service, err := NewRoomService()
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			room, err := service.PopulateRoomWithMonstersAndPlayers(tc.roomConfig, tc.monsterConfigs, tc.playerConfigs)

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

// MockMonsterRepository is a mock implementation of the MonsterRepository interface for testing
type MockMonsterRepository struct {
	xpValues map[string]int
}

func (m *MockMonsterRepository) GetMonsterXP(monsterKey string) (int, error) {
	if xp, ok := m.xpValues[monsterKey]; ok {
		return xp, nil
	}
	return 0, fmt.Errorf("monster not found: %s", monsterKey)
}

// NewMockMonsterRepositoryWithTestData creates a MockMonsterRepository with data from test JSON files
func NewMockMonsterRepositoryWithTestData(t *testing.T) *MockMonsterRepository {
	xpValues, err := testutil.CreateTestMonsterRepository()
	require.NoError(t, err, "Failed to create test monster repository")

	// Ensure we have at least the basic monsters needed for tests
	requiredMonsters := []string{"goblin", "bandit-captain", "adult-blue-dragon"}
	for _, monster := range requiredMonsters {
		monsterKey := "monster_" + monster
		_, exists := xpValues[monsterKey]
		require.True(t, exists, "Required monster %s not found in test data", monsterKey)
	}

	return &MockMonsterRepository{
		xpValues: xpValues,
	}
}

// NewMockItemRepositoryWithTestData creates a mock item repository with real test data
func NewMockItemRepositoryWithTestData(t *testing.T, requiredItems []string) *MockItemRepository {
	items, err := testutil.LoadAllEquipment()
	require.NoError(t, err, "Failed to load item data")

	// Verify that all required items are available
	for _, item := range requiredItems {
		itemKey := "item_" + item
		_, exists := items[itemKey]
		require.True(t, exists, "Required item %s not found in test data", itemKey)
	}

	return &MockItemRepository{
		items: items,
	}
}

func TestBalanceMonsterConfigs(t *testing.T) {
	// Create a mock monster repository
	mockRepo := NewMockMonsterRepositoryWithTestData(t)

	// Create a room service with the mock repository
	roomService := &RoomService{
		monsterRepo: mockRepo,
		balancer:    NewBalancer(mockRepo),
	}

	// Create a test party
	party := entities.Party{
		Members: []entities.PartyMember{
			{Name: "Player1", Level: 5},
			{Name: "Player2", Level: 5},
			{Name: "Player3", Level: 5},
			{Name: "Player4", Level: 5},
		},
	}

	// Create monster configs
	monsterConfigs := []MonsterConfig{
		{Name: "Goblin", Key: "goblin", CR: 0.25, Count: 2, RandomPlace: true},
		{Name: "Bandit Captain", Key: "bandit-captain", CR: 0.5, Count: 1, RandomPlace: true},
	}

	// Test cases
	testCases := []struct {
		name        string
		configs     []MonsterConfig
		party       entities.Party
		difficulty  entities.EncounterDifficulty
		expectError bool
		checkFunc   func(t *testing.T, configs []MonsterConfig)
	}{
		{
			name:        "Balance for medium difficulty",
			configs:     monsterConfigs,
			party:       party,
			difficulty:  entities.EncounterDifficultyMedium,
			expectError: false,
			checkFunc: func(t *testing.T, configs []MonsterConfig) {
				// Verify configs were adjusted
				assert.Len(t, configs, 2)
				totalCount := 0
				for _, config := range configs {
					totalCount += config.Count
				}
				// Original total was 3 (2 goblins + 1 banditcaptain)
				// For a level 5 party of 4, medium difficulty should scale this up
				assert.GreaterOrEqual(t, totalCount, 3)
			},
		},
		{
			name:        "Empty monster configs",
			configs:     []MonsterConfig{},
			party:       party,
			difficulty:  entities.EncounterDifficultyEasy,
			expectError: false,
			checkFunc: func(t *testing.T, configs []MonsterConfig) {
				assert.Empty(t, configs)
			},
		},
		{
			name:        "Empty party",
			configs:     monsterConfigs,
			party:       entities.Party{},
			difficulty:  entities.EncounterDifficultyEasy,
			expectError: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configs, err := roomService.BalanceMonsterConfigs(tc.configs, tc.party, tc.difficulty)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, configs)
				}
			}
		})
	}
}

func TestPopulateRoomWithBalancedMonsters(t *testing.T) {
	// Create a mock monster repository
	mockRepo := NewMockMonsterRepositoryWithTestData(t)

	// Create a room service with the mock repository
	roomService := &RoomService{
		monsterRepo: mockRepo,
		balancer:    NewBalancer(mockRepo),
	}

	// Create a test party
	party := entities.Party{
		Members: []entities.PartyMember{
			{Name: "Player1", Level: 3},
			{Name: "Player2", Level: 3},
		},
	}

	// Create room config
	roomConfig := RoomConfig{
		Width:       10,
		Height:      10,
		LightLevel:  entities.LightLevelBright,
		Description: "Test Room",
		UseGrid:     true,
	}

	// Create monster configs
	monsterConfigs := []MonsterConfig{
		{Name: "Goblin", Key: "goblin", CR: 0.25, Count: 2, RandomPlace: true},
	}

	// Test cases
	testCases := []struct {
		name           string
		roomConfig     RoomConfig
		monsterConfigs []MonsterConfig
		party          entities.Party
		difficulty     entities.EncounterDifficulty
		expectError    bool
		checkFunc      func(t *testing.T, room *entities.Room)
	}{
		{
			name:           "Generate room with balanced monsters",
			roomConfig:     roomConfig,
			monsterConfigs: monsterConfigs,
			party:          party,
			difficulty:     entities.EncounterDifficultyHard,
			expectError:    false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				// Verify room was created with correct dimensions
				assert.Equal(t, 10, room.Width)
				assert.Equal(t, 10, room.Height)

				// Verify monsters were added
				assert.NotEmpty(t, room.Monsters)

				// For a level 3 party of 2, hard difficulty should have scaled up the monsters
				assert.GreaterOrEqual(t, len(room.Monsters), 2)
			},
		},
		{
			name:           "Empty party",
			roomConfig:     roomConfig,
			monsterConfigs: monsterConfigs,
			party:          entities.Party{},
			difficulty:     entities.EncounterDifficultyEasy,
			expectError:    true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			room, err := roomService.PopulateRoomWithBalancedMonsters(tc.roomConfig, tc.monsterConfigs, tc.party, tc.difficulty)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, room)
				}
			}
		})
	}
}

func TestDetermineRoomDifficulty(t *testing.T) {
	// Create a mock monster repository
	mockRepo := NewMockMonsterRepositoryWithTestData(t)

	// Create a room service with the mock repository
	roomService := &RoomService{
		monsterRepo: mockRepo,
		balancer:    NewBalancer(mockRepo),
	}

	// Create a test party
	party := entities.Party{
		Members: []entities.PartyMember{
			{Name: "Player1", Level: 5},
			{Name: "Player2", Level: 5},
			{Name: "Player3", Level: 5},
			{Name: "Player4", Level: 5},
		},
	}

	// Test cases
	testCases := []struct {
		name         string
		setupRoom    func() *entities.Room
		party        entities.Party
		expectedDiff entities.EncounterDifficulty
		expectError  bool
	}{
		{
			name: "Empty room",
			setupRoom: func() *entities.Room {
				return &entities.Room{
					Width:      10,
					Height:     10,
					LightLevel: entities.LightLevelBright,
					Monsters:   []entities.Monster{},
				}
			},
			party:        party,
			expectedDiff: entities.EncounterDifficultyEasy,
			expectError:  false,
		},
		{
			name: "Easy encounter",
			setupRoom: func() *entities.Room {
				return &entities.Room{
					Width:      10,
					Height:     10,
					LightLevel: entities.LightLevelBright,
					Monsters: []entities.Monster{
						{ID: "1", Name: "Goblin", CR: 0.25},
						{ID: "2", Name: "Goblin", CR: 0.25},
					},
				}
			},
			party:        party,
			expectedDiff: entities.EncounterDifficultyEasy,
			expectError:  false,
		},
		{
			name: "Hard encounter",
			setupRoom: func() *entities.Room {
				return &entities.Room{
					Width:      10,
					Height:     10,
					LightLevel: entities.LightLevelBright,
					Monsters: []entities.Monster{
						{ID: "1", Name: "Adult Blue Dragon", CR: 5},
						{ID: "2", Name: "Adult Blue Dragon", CR: 5},
					},
				}
			},
			party:        party,
			expectedDiff: entities.EncounterDifficultyDeadly,
			expectError:  false,
		},
		{
			name:        "Nil room",
			setupRoom:   func() *entities.Room { return nil },
			party:       party,
			expectError: true,
		},
		{
			name: "Empty party",
			setupRoom: func() *entities.Room {
				return &entities.Room{
					Width:      10,
					Height:     10,
					LightLevel: entities.LightLevelBright,
					Monsters: []entities.Monster{
						{ID: "1", Name: "Goblin", CR: 0.25},
					},
				}
			},
			party:       entities.Party{},
			expectError: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			room := tc.setupRoom()
			difficulty, err := roomService.DetermineRoomDifficulty(room, tc.party)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDiff, difficulty)
			}
		})
	}
}

// MockItemRepository is a mock implementation of the ItemRepository interface for testing
type MockItemRepository struct {
	items map[string]*entities.Item
}

func (m *MockItemRepository) GetItemByKey(key string) (*entities.Item, error) {
	item, exists := m.items[key]
	if !exists {
		return nil, fmt.Errorf("item with key %s not found", key)
	}
	return item, nil
}

func (m *MockItemRepository) GetRandomItems(count int) ([]*entities.Item, error) {
	if len(m.items) == 0 {
		return nil, fmt.Errorf("no items available")
	}

	result := make([]*entities.Item, 0, count)

	// Convert map to slice for easier iteration
	itemSlice := make([]*entities.Item, 0, len(m.items))
	for _, item := range m.items {
		itemSlice = append(itemSlice, item)
	}

	// Duplicate items as needed to meet the requested count
	for i := 0; i < count; i++ {
		// Use modulo to cycle through available items
		itemIndex := i % len(itemSlice)
		// Create a copy of the item with a unique ID
		itemCopy := *itemSlice[itemIndex]
		itemCopy.ID = fmt.Sprintf("%s_%d", itemCopy.ID, i)
		result = append(result, &itemCopy)
	}

	return result, nil
}

func (m *MockItemRepository) GetRandomItemsByCategory(category string, count int) ([]*entities.Item, error) {
	if len(m.items) == 0 {
		return nil, fmt.Errorf("no items available")
	}

	// Find items matching the category
	matchingItems := make([]*entities.Item, 0)
	for _, item := range m.items {
		if item.Category == category {
			matchingItems = append(matchingItems, item)
		}
	}

	if len(matchingItems) == 0 {
		return nil, fmt.Errorf("no items found for category %s", category)
	}

	result := make([]*entities.Item, 0, count)

	// Duplicate items as needed to meet the requested count
	for i := 0; i < count; i++ {
		// Use modulo to cycle through available items
		itemIndex := i % len(matchingItems)
		// Create a copy of the item with a unique ID
		itemCopy := *matchingItems[itemIndex]
		itemCopy.ID = fmt.Sprintf("%s_%d", itemCopy.ID, i)
		result = append(result, &itemCopy)
	}

	return result, nil
}

func TestPopulateTreasureRoom(t *testing.T) {
	// Create mock item repository with real test data
	mockItemRepo := NewMockItemRepositoryWithTestData(t, []string{"abacus", "battleaxe", "studded-leather-armor"})

	// Create a room service with the mock repository
	roomService := &RoomService{
		itemRepo: mockItemRepo,
	}

	// Test cases
	testCases := []struct {
		name                   string
		roomConfig             RoomConfig
		itemCount              int
		guardianMonsterConfigs []MonsterConfig
		expectError            bool
		checkFunc              func(t *testing.T, room *entities.Room)
	}{
		{
			name:        "Basic treasure room",
			roomConfig:  createTestRoomConfig(10, 10, entities.LightLevelBright, true),
			itemCount:   3,
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				// Check room type
				assert.NotNil(t, room.RoomType)
				assert.Equal(t, "treasure", room.RoomType.Type())

				// Check items were added
				assert.Len(t, room.Items, 3)

				// Verify grid has items placed
				if room.Grid != nil {
					itemCellCount := 0
					for y := 0; y < room.Height; y++ {
						for x := 0; x < room.Width; x++ {
							if room.Grid[y][x].Type == entities.CellItem {
								itemCellCount++
							}
						}
					}
					assert.Equal(t, 3, itemCellCount)
				}
			},
		},
		{
			name:       "Treasure room with guardian",
			roomConfig: createTestRoomConfig(15, 15, entities.LightLevelDim, true),
			itemCount:  2,
			guardianMonsterConfigs: []MonsterConfig{
				createTestMonsterConfigWithRealData(t, "adult-blue-dragon", 1, true, nil),
			},
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				// Check room type
				assert.NotNil(t, room.RoomType)
				assert.Equal(t, "treasure", room.RoomType.Type())

				// Check items were added
				assert.Len(t, room.Items, 2)

				// Check guardian was added
				assert.Len(t, room.Monsters, 1)
				assert.Equal(t, "Adult Blue Dragon", room.Monsters[0].Name)
			},
		},
		{
			name:        "Too many items for room size",
			roomConfig:  createTestRoomConfig(2, 2, entities.LightLevelBright, true),
			itemCount:   10, // Too many for a 2x2 room
			expectError: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			room, err := roomService.PopulateTreasureRoom(tc.roomConfig, tc.itemCount, tc.guardianMonsterConfigs)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, room)
				}
			}
		})
	}
}

func TestPopulateRandomTreasureRoomWithParty(t *testing.T) {
	// Create mock item repository with test items
	mockItemRepo := &MockItemRepository{
		items: map[string]*entities.Item{
			"item_sword": {
				ID:         "1",
				Key:        "item_sword",
				Name:       "Longsword",
				Type:       "weapon",
				Category:   "martial-weapons",
				Value:      15,
				ValueUnit:  "gp",
				Weight:     3,
				DamageDice: "1d8",
				DamageType: "slashing",
			},
			"item_armor": {
				ID:                  "2",
				Key:                 "item_armor",
				Name:                "Chain Mail",
				Type:                "armor",
				Category:            "heavy-armor",
				Value:               75,
				ValueUnit:           "gp",
				Weight:              55,
				ArmorClass:          16,
				StealthDisadvantage: true,
			},
			"item_potion": {
				ID:        "3",
				Key:       "item_potion",
				Name:      "Potion of Healing",
				Type:      "potion",
				Category:  "potion",
				Value:     50,
				ValueUnit: "gp",
				Weight:    1,
			},
		},
	}

	// Create a mock balancer for monster balancing
	mockBalancer := &MockBalancer{
		adjustFunc: func(configs []MonsterConfig, party entities.Party, difficulty entities.EncounterDifficulty) ([]MonsterConfig, error) {
			// Just return the same configs for testing
			return configs, nil
		},
	}

	// Create a room service with the mock repositories
	roomService := &RoomService{
		itemRepo: mockItemRepo,
		balancer: mockBalancer,
	}

	// Test cases
	testCases := []struct {
		name            string
		roomConfig      RoomConfig
		party           entities.Party
		includeGuardian bool
		difficulty      entities.EncounterDifficulty
		expectError     bool
		checkFunc       func(t *testing.T, room *entities.Room)
	}{
		{
			name:       "Small party, easy difficulty, no guardian",
			roomConfig: createTestRoomConfig(10, 10, entities.LightLevelBright, true),
			party: entities.Party{
				Members: []entities.PartyMember{
					{Name: "Player1", Level: 1},
					{Name: "Player2", Level: 2},
				},
			},
			includeGuardian: false,
			difficulty:      entities.EncounterDifficultyEasy,
			expectError:     false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				// Check room type
				assert.NotNil(t, room.RoomType)
				assert.Equal(t, "treasure", room.RoomType.Type())

				// Check items were added (for small party, low level, easy difficulty)
				// Expected: 2 players + 1 item per 3 levels (avg 1.5 level) = ~2.5 items
				// With 0.75 multiplier for easy difficulty = ~1.9 items, rounded to at least 1
				assert.GreaterOrEqual(t, len(room.Items), 1)

				// Check party members were added
				assert.Len(t, room.Players, 2)

				// Check no monsters (guardian disabled)
				assert.Empty(t, room.Monsters)

				// Check description contains reminder about clearing
				assert.Contains(t, room.Description, "clear the room after collecting all treasure")
			},
		},
		{
			name:       "Medium party, medium difficulty, with guardian",
			roomConfig: createTestRoomConfig(15, 15, entities.LightLevelDim, true),
			party: entities.Party{
				Members: []entities.PartyMember{
					{Name: "Player1", Level: 3},
					{Name: "Player2", Level: 4},
					{Name: "Player3", Level: 5},
					{Name: "Player4", Level: 4},
				},
			},
			includeGuardian: true,
			difficulty:      entities.EncounterDifficultyMedium,
			expectError:     false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				// Check room type
				assert.NotNil(t, room.RoomType)
				assert.Equal(t, "treasure", room.RoomType.Type())

				// Check items were added
				// Expected: 4 players + 1 item per 3 levels (avg 4 level) = ~5.3 items
				// With 1.0 multiplier for medium difficulty = ~5.3 items
				assert.GreaterOrEqual(t, len(room.Items), 5)

				// Check party members were added
				assert.Len(t, room.Players, 4)

				// Check guardian was added
				assert.NotEmpty(t, room.Monsters)

				// Check description contains reminder about clearing
				assert.Contains(t, room.Description, "clear the room after collecting all treasure")
			},
		},
		{
			name:       "Large party, hard difficulty, with guardian",
			roomConfig: createTestRoomConfig(20, 20, entities.LightLevelDark, true),
			party: entities.Party{
				Members: []entities.PartyMember{
					{Name: "Player1", Level: 8},
					{Name: "Player2", Level: 9},
					{Name: "Player3", Level: 10},
					{Name: "Player4", Level: 9},
					{Name: "Player5", Level: 8},
					{Name: "Player6", Level: 10},
				},
			},
			includeGuardian: true,
			difficulty:      entities.EncounterDifficultyHard,
			expectError:     false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				// Check room type
				assert.NotNil(t, room.RoomType)
				assert.Equal(t, "treasure", room.RoomType.Type())

				// Check items were added
				// Expected: 6 players + 1 item per 3 levels (avg 9 level) = ~9 items
				// With 1.25 multiplier for hard difficulty = ~11.25 items
				assert.GreaterOrEqual(t, len(room.Items), 11)

				// Check party members were added
				assert.Len(t, room.Players, 6)

				// Check guardian was added
				assert.NotEmpty(t, room.Monsters)

				// Check description contains reminder about clearing
				assert.Contains(t, room.Description, "clear the room after collecting all treasure")
			},
		},
		{
			name:       "Too small room for party and items",
			roomConfig: createTestRoomConfig(2, 2, entities.LightLevelBright, true),
			party: entities.Party{
				Members: []entities.PartyMember{
					{Name: "Player1", Level: 10},
					{Name: "Player2", Level: 10},
				},
			},
			includeGuardian: true,
			difficulty:      entities.EncounterDifficultyDeadly,
			expectError:     true, // Room is too small for party + items + guardian
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			room, err := roomService.PopulateRandomTreasureRoomWithParty(tc.roomConfig, tc.party, tc.includeGuardian, tc.difficulty)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, room)
				}
			}
		})
	}
}

// MockBalancer is a mock implementation of the monster balancer
type MockBalancer struct {
	adjustFunc func(configs []MonsterConfig, party entities.Party, difficulty entities.EncounterDifficulty) ([]MonsterConfig, error)
}

func (m *MockBalancer) AdjustMonsterSelection(configs []MonsterConfig, party entities.Party, difficulty entities.EncounterDifficulty) ([]MonsterConfig, error) {
	if m.adjustFunc != nil {
		return m.adjustFunc(configs, party, difficulty)
	}
	return configs, nil
}

func (m *MockBalancer) DetermineEncounterDifficulty(monsters []entities.Monster, party entities.Party) (entities.EncounterDifficulty, error) {
	// Simple mock implementation - just return medium difficulty
	return entities.EncounterDifficultyMedium, nil
}

func (m *MockBalancer) CalculateTargetCR(party entities.Party, difficulty entities.EncounterDifficulty) (float64, error) {
	// Calculate average party level as CR
	totalLevels := 0
	for _, member := range party.Members {
		totalLevels += member.Level
	}
	avgLevel := float64(totalLevels) / float64(len(party.Members))
	return avgLevel, nil
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
		createTestMonsterConfigWithRealData(t, "goblin", 3, true, nil),
		createTestMonsterConfigWithRealData(t, "bandit-captain", 2, false, &entities.Position{X: 5, Y: 5}),
	}
	err = service.AddMonstersToRoom(room, monsterConfigs)
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
	err = service.AddPlayersToRoom(room, playerConfigs)
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
		createTestItemConfigWithRealData(t, "abacus", 2, true, nil),
		createTestItemConfigWithRealData(t, "battleaxe", 1, false, &entities.Position{X: 7, Y: 7}),
	}

	// Set up the item repository with real test data
	service.itemRepo = NewMockItemRepositoryWithTestData(t, []string{"abacus", "battleaxe"})

	err = service.AddItemsToRoom(room, itemConfigs)
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
	// Create a mock monster repository for testing
	mockRepo := NewMockMonsterRepositoryWithTestData(t)

	// Create a service with the mock repository
	service := &RoomService{
		monsterRepo: mockRepo,
	}

	// Create a gridless room with monsters
	room := entities.NewRoom(10, 10, entities.LightLevelBright)
	// Explicitly not initializing grid
	assert.Nil(t, room.Grid)

	// Add monsters directly using the placement interface
	goblin := entities.Monster{ID: "1", Key: "monster_goblin", Name: "Goblin", Position: entities.Position{X: 1, Y: 1}}
	banditcaptain := entities.Monster{ID: "2", Key: "monster_bandit-captain", Name: "Bandit Captain", Position: entities.Position{X: 3, Y: 3}}
	adultbluedragon := entities.Monster{ID: "3", Key: "monster_adult-blue-dragon", Name: "Adult Blue Dragon", Position: entities.Position{X: 5, Y: 5}}

	entities.AddMonster(room, goblin)
	entities.AddMonster(room, banditcaptain)
	entities.AddMonster(room, adultbluedragon)

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
