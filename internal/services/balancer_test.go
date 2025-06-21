package services

import (
	"testing"

	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
	"github.com/stretchr/testify/assert"
)

// createTestParty creates a party with the given number of members at the specified level
func createTestParty(memberCount int, level int) entities.Party {
	members := make([]entities.PartyMember, memberCount)
	for i := 0; i < memberCount; i++ {
		members[i] = entities.PartyMember{
			Name:  "Player",
			Level: level,
		}
	}
	return entities.Party{Members: members}
}

// createTestMonsters creates a slice of monsters with the given CRs
func createTestMonsters(crs ...float64) []entities.Monster {
	monsters := make([]entities.Monster, len(crs))
	for i, cr := range crs {
		monsters[i] = entities.Monster{
			ID:   "monster_" + string(rune(i+65)), // Rune is an alias for int32 and typically used to differentiate character values
			Name: "Monster " + string(rune(i+65)),
			CR:   cr,
		}
	}
	return monsters
}

// createTestBalancer creates a balancer with a mock repository for testing
func createTestBalancer() *StandardBalancer {
	mockRepo := &MockMonsterRepository{
		xpValues: map[string]int{
			"monster_goblin": 50,
			"monster_orc":    100,
			"monster_troll":  450,
		},
	}
	return NewBalancer(mockRepo)
}

func TestCalculateTargetCR(t *testing.T) {
	balancer := createTestBalancer()

	testCases := []struct {
		name           string
		party          entities.Party
		difficulty     entities.EncounterDifficulty
		expectedCR     float64
		expectError    bool
		errorSubstring string
	}{
		{
			name:        "Solo player easy encounter",
			party:       createTestParty(1, 5),
			difficulty:  entities.EncounterDifficultyEasy,
			expectedCR:  1.25, // 5 * 0.5 * 0.5 = 1.25
			expectError: false,
		},
		{
			name:        "Four players medium encounter",
			party:       createTestParty(4, 3),
			difficulty:  entities.EncounterDifficultyMedium,
			expectedCR:  2.25, // 3 * 0.75 * 1.0 = 2.25
			expectError: false,
		},
		{
			name:        "Six players hard encounter",
			party:       createTestParty(6, 10),
			difficulty:  entities.EncounterDifficultyHard,
			expectedCR:  15, // 10 * 1.0 * 1.5 = 15
			expectError: false,
		},
		{
			name:        "Five players deadly encounter",
			party:       createTestParty(5, 8),
			difficulty:  entities.EncounterDifficultyDeadly,
			expectedCR:  15, // 8 * 1.5 * 1.25 = 15
			expectError: false,
		},
		{
			name:           "Empty party",
			party:          entities.Party{Members: []entities.PartyMember{}},
			difficulty:     entities.EncounterDifficultyEasy,
			expectError:    true,
			errorSubstring: "party cannot be empty",
		},
		{
			name:           "Invalid difficulty",
			party:          createTestParty(4, 5),
			difficulty:     "impossible",
			expectError:    true,
			errorSubstring: "invalid difficulty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cr, err := balancer.CalculateTargetCR(tc.party, tc.difficulty)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorSubstring != "" {
					assert.Contains(t, err.Error(), tc.errorSubstring)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCR, cr)
			}
		})
	}
}

func TestDetermineEncounterDifficulty(t *testing.T) {
	balancer := createTestBalancer()

	testCases := []struct {
		name           string
		monsters       []entities.Monster
		party          entities.Party
		expectedDiff   entities.EncounterDifficulty
		expectError    bool
		errorSubstring string
	}{
		{
			name:         "Empty monsters, returns easy",
			monsters:     []entities.Monster{},
			party:        createTestParty(4, 5),
			expectedDiff: entities.EncounterDifficultyEasy,
			expectError:  false,
		},
		{
			name:         "Low CR monsters for party level, returns easy",
			monsters:     createTestMonsters(0.25, 0.25),
			party:        createTestParty(4, 5),
			expectedDiff: entities.EncounterDifficultyEasy,
			expectError:  false,
		},
		{
			name:         "Medium difficulty monsters",
			monsters:     createTestMonsters(1, 1, 1),
			party:        createTestParty(4, 5),
			expectedDiff: entities.EncounterDifficultyMedium,
			expectError:  false,
		},
		{
			name:         "Hard difficulty monsters",
			monsters:     createTestMonsters(2, 2, 1),
			party:        createTestParty(4, 5),
			expectedDiff: entities.EncounterDifficultyHard,
			expectError:  false,
		},
		{
			name:         "Deadly difficulty monsters",
			monsters:     createTestMonsters(5, 5, 5),
			party:        createTestParty(4, 5),
			expectedDiff: entities.EncounterDifficultyDeadly,
			expectError:  false,
		},
		{
			name:         "Solo player adjustment",
			monsters:     createTestMonsters(1),
			party:        createTestParty(1, 5),
			expectedDiff: entities.EncounterDifficultyDeadly, // 1 CR vs level 5 solo player (with 0.5 adjustment) is deadly
			expectError:  false,
		},
		{
			name:           "Empty party",
			monsters:       createTestMonsters(1),
			party:          entities.Party{Members: []entities.PartyMember{}},
			expectError:    true,
			errorSubstring: "party cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			difficulty, err := balancer.DetermineEncounterDifficulty(tc.monsters, tc.party)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorSubstring != "" {
					assert.Contains(t, err.Error(), tc.errorSubstring)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedDiff, difficulty)
			}
		})
	}
}

func TestAdjustMonsterSelection(t *testing.T) {
	balancer := createTestBalancer()

	testCases := []struct {
		name           string
		monsterConfigs []MonsterConfig
		party          entities.Party
		difficulty     entities.EncounterDifficulty
		checkFunc      func(*testing.T, []MonsterConfig)
		expectError    bool
		errorSubstring string
	}{
		{
			name: "Already balanced encounter",
			monsterConfigs: []MonsterConfig{
				{Name: "Goblin", Key: "monster_goblin", CR: 0.25, Count: 4, RandomPlace: true},
			},
			party:      createTestParty(4, 1), // Level 1 party
			difficulty: entities.EncounterDifficultyEasy,
			checkFunc: func(t *testing.T, configs []MonsterConfig) {
				// Should remain unchanged since it's already balanced
				assert.Len(t, configs, 1)
				assert.Equal(t, "Goblin", configs[0].Name)
				assert.Equal(t, 4, configs[0].Count)
			},
			expectError: false,
		},
		{
			name: "Scale up monsters for harder difficulty",
			monsterConfigs: []MonsterConfig{
				{Name: "Goblin", Key: "monster_goblin", CR: 0.25, Count: 2, RandomPlace: true},
				{Name: "Orc", Key: "monster_orc", CR: 0.5, Count: 1, RandomPlace: true},
			},
			party:      createTestParty(4, 5),              // Level 5 party
			difficulty: entities.EncounterDifficultyDeadly, // Should scale up significantly
			checkFunc: func(t *testing.T, configs []MonsterConfig) {
				assert.Len(t, configs, 2)

				// Total CR before: 2*0.25 + 1*0.5 = 1
				// Target CR for deadly: 5*1.5*1.0 = 7.5
				// Scaling factor: 7.5/1 = 7.5

				// Check that counts are scaled up
				totalCountBefore := 3 // 2 goblins + 1 orc
				totalCountAfter := 0
				for _, config := range configs {
					totalCountAfter += config.Count
				}

				assert.Greater(t, totalCountAfter, totalCountBefore)
			},
			expectError: false,
		},
		{
			name: "Scale down monsters for easier difficulty",
			monsterConfigs: []MonsterConfig{
				{Name: "Troll", Key: "monster_troll", CR: 5, Count: 3, RandomPlace: true},
			},
			party:      createTestParty(4, 5),            // Level 5 party
			difficulty: entities.EncounterDifficultyEasy, // Should scale down significantly
			checkFunc: func(t *testing.T, configs []MonsterConfig) {
				assert.Len(t, configs, 1)

				// Total CR before: 3*5 = 15
				// Target CR for easy: 5*0.5*1.0 = 2.5
				// Scaling factor: 2.5/15 = 0.167

				// Check that count is scaled down but at least 1
				assert.Less(t, configs[0].Count, 3)
				assert.GreaterOrEqual(t, configs[0].Count, 1)
			},
			expectError: false,
		},
		{
			name:           "Empty party",
			monsterConfigs: []MonsterConfig{{Name: "Goblin", Key: "monster_goblin", CR: 0.25, Count: 2, RandomPlace: true}},
			party:          entities.Party{Members: []entities.PartyMember{}},
			difficulty:     entities.EncounterDifficultyEasy,
			expectError:    true,
			errorSubstring: "party cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			adjustedConfigs, err := balancer.AdjustMonsterSelection(tc.monsterConfigs, tc.party, tc.difficulty)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorSubstring != "" {
					assert.Contains(t, err.Error(), tc.errorSubstring)
				}
			} else {
				assert.NoError(t, err)
				if tc.checkFunc != nil {
					tc.checkFunc(t, adjustedConfigs)
				}
			}
		})
	}
}
