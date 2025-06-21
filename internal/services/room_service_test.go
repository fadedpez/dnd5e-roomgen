package services

import (
	"fmt"
	"testing"

	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
	"github.com/stretchr/testify/assert"
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

// createTestMonsterConfig creates a standard monster configuration for testing
func createTestMonsterConfig(name string, cr float64, count int, randomPlace bool, position *entities.Position) MonsterConfig {
	config := MonsterConfig{
		Name:        name,
		Key:         "monster_" + name,
		CR:          cr,
		Count:       count,
		RandomPlace: randomPlace,
	}

	if position != nil {
		config.Position = position
	}

	return config
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
				createTestMonsterConfig("Goblin", 0.25, 1, true, nil),
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
				createTestMonsterConfig("Goblin", 0.25, 3, true, nil),
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
				createTestMonsterConfig("Orc", 0.5, 1, false, &entities.Position{X: 3, Y: 4}),
			},
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				assert.Len(t, room.Monsters, 1)
				assert.Equal(t, "Orc", room.Monsters[0].Name)
				assert.Equal(t, 3, room.Monsters[0].Position.X)
				assert.Equal(t, 4, room.Monsters[0].Position.Y)
			},
		},
		{
			name: "Auto-initialize grid",
			roomSetup: func() *entities.Room {
				// Create room without grid
				return entities.NewRoom(10, 10, entities.LightLevelBright)
			},
			monsterConfigs: []MonsterConfig{
				createTestMonsterConfig("Goblin", 0.25, 1, true, nil),
			},
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				assert.NotNil(t, room.Grid)
				assert.Len(t, room.Monsters, 1)
			},
		},
		{
			name:      "Nil room",
			roomSetup: func() *entities.Room { return nil },
			monsterConfigs: []MonsterConfig{
				createTestMonsterConfig("Goblin", 0.25, 1, true, nil),
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
					config := createTestMonsterConfig("Goblin", 0.25, 1, false, nil)
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
				createTestMonsterConfig("Goblin", 0.25, 2, true, nil),
			},
			expectError: false,
			checkFunc: func(t *testing.T, room *entities.Room) {
				assertRoomProperties(t, room, 10, 10, entities.LightLevelBright, true)
				assert.Len(t, room.Monsters, 2)
				assert.Equal(t, "Goblin", room.Monsters[0].Name)
			},
		},
		{
			name:       "Invalid room dimensions",
			roomConfig: createTestRoomConfig(-5, 10, entities.LightLevelBright, true),
			monsterConfigs: []MonsterConfig{
				createTestMonsterConfig("Goblin", 0.25, 1, true, nil),
			},
			expectError: true,
		},
		{
			name:       "Invalid monster placement",
			roomConfig: createTestRoomConfig(10, 10, entities.LightLevelBright, true),
			monsterConfigs: []MonsterConfig{
				func() MonsterConfig {
					config := createTestMonsterConfig("Goblin", 0.25, 1, false, nil)
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
	mockRepo := &MockMonsterRepository{
		xpValues: map[string]int{
			"monster_Goblin": 50,
			"monster_Orc":    100,
			"monster_Troll":  450,
		},
	}

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
				goblin := entities.Monster{ID: "1", Key: "monster_Goblin", Name: "Goblin", Position: entities.Position{X: 1, Y: 1}}
				orc := entities.Monster{ID: "2", Key: "monster_Orc", Name: "Orc", Position: entities.Position{X: 3, Y: 3}}
				troll := entities.Monster{ID: "3", Key: "monster_Troll", Name: "Troll", Position: entities.Position{X: 5, Y: 5}}

				entities.AddMonster(room, goblin)
				entities.AddMonster(room, orc)
				entities.AddMonster(room, troll)

				return room
			},
			monsterIDs:    []string{}, // Empty means remove all
			expectedXP:    600,        // 50 + 100 + 450
			expectedCount: 0,
			notRemovedIDs: []string{}, // All should be removed
		},
		{
			name: "Remove specific monsters",
			setupRoom: func() *entities.Room {
				room := entities.NewRoom(10, 10, entities.LightLevelBright)
				entities.InitializeGrid(room)

				// Add three different monsters
				goblin := entities.Monster{ID: "1", Key: "monster_Goblin", Name: "Goblin", Position: entities.Position{X: 1, Y: 1}}
				orc := entities.Monster{ID: "2", Key: "monster_Orc", Name: "Orc", Position: entities.Position{X: 3, Y: 3}}
				troll := entities.Monster{ID: "3", Key: "monster_Troll", Name: "Troll", Position: entities.Position{X: 5, Y: 5}}

				entities.AddMonster(room, goblin)
				entities.AddMonster(room, orc)
				entities.AddMonster(room, troll)

				return room
			},
			monsterIDs:    []string{"1", "3"}, // Remove goblin and troll
			expectedXP:    500,                // 50 + 450
			expectedCount: 1,                  // Only orc should remain
			notRemovedIDs: []string{},         // All specified monsters should be removed
		},
		{
			name: "Remove non-existent monsters",
			setupRoom: func() *entities.Room {
				room := entities.NewRoom(10, 10, entities.LightLevelBright)
				entities.InitializeGrid(room)

				// Add one monster
				goblin := entities.Monster{ID: "1", Key: "monster_Goblin", Name: "Goblin", Position: entities.Position{X: 1, Y: 1}}
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

			xp, notRemoved, err := service.CleanupRoom(room, tc.monsterIDs)

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
