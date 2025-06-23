package services

import (
	"testing"

	apientities "github.com/fadedpez/dnd5e-api/entities"
	"github.com/fadedpez/dnd5e-roomgen/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertAPIMonsterToConfig(t *testing.T) {
	// Load test monster data
	monsterData, err := testutil.LoadAllMonsters()
	require.NoError(t, err, "Failed to load test monster data")

	// Test cases
	testCases := []struct {
		name          string
		monsterKey    string
		count         int
		expectedKey   string
		expectedName  string
		expectedCount int
		expectedCR    float64
	}{
		{
			name:          "Convert goblin with count 5",
			monsterKey:    "goblin",
			count:         5,
			expectedKey:   "goblin",
			expectedName:  "Goblin",
			expectedCount: 5,
			expectedCR:    0.25,
		},
		{
			name:          "Convert dragon with count 1",
			monsterKey:    "adult-blue-dragon",
			count:         1,
			expectedKey:   "adult-blue-dragon",
			expectedName:  "Adult Blue Dragon",
			expectedCount: 1,
			expectedCR:    16.0,
		},
		{
			name:          "Convert with negative count should default to 1",
			monsterKey:    "goblin",
			count:         -3,
			expectedKey:   "goblin",
			expectedName:  "Goblin",
			expectedCount: 1,
			expectedCR:    0.25,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the test monster data
			monster := monsterData[tc.monsterKey]
			require.NotNil(t, monster, "Test monster data not found for key: %s", tc.monsterKey)

			// Debug: Print the monster data to see what we're working with
			t.Logf("Monster from test data: %+v", monster)
			t.Logf("Challenge Rating from test data: %f", monster.ChallengeRating)

			// Create an API monster from the test data
			apiMonster := &apientities.Monster{
				Key:             monster.Index,
				Name:            monster.Name,
				ChallengeRating: float32(monster.ChallengeRating),
				XP:              monster.XP,
			}

			// Debug: Print the API monster
			t.Logf("API Monster: %+v", apiMonster)
			t.Logf("Challenge Rating from API Monster: %f", apiMonster.ChallengeRating)

			// Call the converter function
			config := ConvertAPIMonsterToConfig(apiMonster, tc.count)

			// Debug: Print the config
			t.Logf("Monster Config: %+v", config)
			t.Logf("CR from Monster Config: %f", config.CR)

			// Assert the results
			assert.Equal(t, tc.expectedKey, config.Key, "Monster key should match")
			assert.Equal(t, tc.expectedName, config.Name, "Monster name should match")
			assert.Equal(t, tc.expectedCount, config.Count, "Monster count should match")
			assert.Equal(t, tc.expectedCR, config.CR, "Monster CR should match")
		})
	}
}

func TestDebugMonsterData(t *testing.T) {
	// Load the goblin monster directly to debug
	goblin, err := testutil.LoadMonster("goblin")
	require.NoError(t, err, "Failed to load goblin test data")

	// Print the goblin data
	t.Logf("Goblin data: %+v", goblin)
	t.Logf("Goblin challenge rating: %f", goblin.ChallengeRating)

	// Create API monster
	apiMonster := &apientities.Monster{
		Key:             goblin.Index,
		Name:            goblin.Name,
		ChallengeRating: float32(goblin.ChallengeRating),
		XP:              goblin.XP,
	}

	t.Logf("API Monster: %+v", apiMonster)
	t.Logf("API Monster challenge rating: %f", apiMonster.ChallengeRating)

	// Convert to config
	config := ConvertAPIMonsterToConfig(apiMonster, 1)

	t.Logf("Monster Config: %+v", config)
	t.Logf("Monster Config CR: %f", config.CR)

	// Verify the values directly
	assert.Equal(t, "goblin", config.Key, "Key should be goblin")
	assert.Equal(t, "Goblin", config.Name, "Name should be Goblin")
	assert.Equal(t, 1, config.Count, "Count should be 1")
	assert.Equal(t, 0.25, config.CR, "CR should be 0.25")
}

func TestConvertAPIMonsterSliceToConfigs(t *testing.T) {
	// Create test monsters directly instead of using test data
	apiMonsters := []*apientities.Monster{
		{
			Key:             "goblin",
			Name:            "Goblin",
			ChallengeRating: 0.25,
			XP:              50,
		},
		{
			Key:             "adult-blue-dragon",
			Name:            "Adult Blue Dragon",
			ChallengeRating: 16.0,
			XP:              11500,
		},
	}

	// Test with count 3
	configs := ConvertAPIMonsterSliceToConfigs(apiMonsters, 3)

	// Assert results
	assert.Len(t, configs, 2, "Should have 2 monster configs")

	// Check first monster (goblin)
	assert.Equal(t, "goblin", configs[0].Key, "First monster key should be goblin")
	assert.Equal(t, "Goblin", configs[0].Name, "First monster name should be Goblin")
	assert.Equal(t, 3, configs[0].Count, "First monster count should be 3")
	assert.Equal(t, 0.25, configs[0].CR, "First monster CR should be 0.25")

	// Check second monster (dragon)
	assert.Equal(t, "adult-blue-dragon", configs[1].Key, "Second monster key should be adult-blue-dragon")
	assert.Equal(t, "Adult Blue Dragon", configs[1].Name, "Second monster name should be Adult Blue Dragon")
	assert.Equal(t, 3, configs[1].Count, "Second monster count should be 3")
	assert.Equal(t, 16.0, configs[1].CR, "Second monster CR should be 16.0")
}

func TestConvertAPIItemToConfig(t *testing.T) {
	// Create test equipment
	testEquipment := &apientities.Equipment{
		Key:  "shortsword",
		Name: "Shortsword",
		Cost: &apientities.Cost{
			Quantity: 10,
			Unit:     "gp",
		},
		Weight: 2,
	}

	// Test cases
	testCases := []struct {
		name          string
		count         int
		expectedKey   string
		expectedName  string
		expectedCount int
	}{
		{
			name:          "Convert item with count 2",
			count:         2,
			expectedKey:   "shortsword",
			expectedName:  "Shortsword",
			expectedCount: 2,
		},
		{
			name:          "Convert with negative count should default to 1",
			count:         -1,
			expectedKey:   "shortsword",
			expectedName:  "Shortsword",
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the converter function
			config := ConvertAPIItemToConfig(testEquipment, tc.count)

			// Assert results
			assert.Equal(t, tc.expectedKey, config.Key, "Item key should match")
			assert.Equal(t, tc.expectedName, config.Name, "Item name should match")
			assert.Equal(t, tc.expectedCount, config.Count, "Item count should match")
			assert.True(t, config.RandomPlace, "RandomPlace should default to true")
			assert.Nil(t, config.Position, "Position should be nil")
		})
	}
}

func TestConvertAPIItemSliceToConfigs(t *testing.T) {
	// Create test equipment items
	testEquipment := []*apientities.Equipment{
		{
			Key:  "shortsword",
			Name: "Shortsword",
			Cost: &apientities.Cost{
				Quantity: 10,
				Unit:     "gp",
			},
			Weight: 2,
		},
		{
			Key:  "potion-of-healing",
			Name: "Potion of Healing",
			Cost: &apientities.Cost{
				Quantity: 50,
				Unit:     "gp",
			},
			Weight: 1,
		},
	}

	// Test with count 2
	configs := ConvertAPIItemSliceToConfigs(testEquipment, 2)

	// Assert results
	assert.Len(t, configs, 2, "Should have 2 item configs")

	// Check first item (shortsword)
	assert.Equal(t, "shortsword", configs[0].Key, "First item key should be shortsword")
	assert.Equal(t, "Shortsword", configs[0].Name, "First item name should be Shortsword")
	assert.Equal(t, 2, configs[0].Count, "First item count should be 2")
	assert.True(t, configs[0].RandomPlace, "First item RandomPlace should be true")

	// Check second item (potion)
	assert.Equal(t, "potion-of-healing", configs[1].Key, "Second item key should be potion-of-healing")
	assert.Equal(t, "Potion of Healing", configs[1].Name, "Second item name should be Potion of Healing")
	assert.Equal(t, 2, configs[1].Count, "Second item count should be 2")
	assert.True(t, configs[1].RandomPlace, "Second item RandomPlace should be true")
}
