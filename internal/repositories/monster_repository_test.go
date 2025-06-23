package repositories

import (
	"errors"
	"testing"

	"github.com/fadedpez/dnd5e-roomgen/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMonsterRepository is a simple implementation of MonsterRepository for testing
type TestMonsterRepository struct {
	xpValues map[string]int
}

// GetMonsterXP implements the MonsterRepository interface for testing
func (r *TestMonsterRepository) GetMonsterXP(monsterKey string) (int, error) {
	if xp, ok := r.xpValues[monsterKey]; ok {
		return xp, nil
	}
	return 0, errors.New("monster not found")
}

func TestGetMonsterXP(t *testing.T) {
	// Load monster data from test files
	monsterData, err := testutil.LoadAllMonsters()
	require.NoError(t, err, "Failed to load test monster data")

	// Get specific monsters for testing
	goblin := monsterData["goblin"]
	require.NotNil(t, goblin, "Failed to get goblin test data")

	banditCaptain := monsterData["bandit-captain"]
	require.NotNil(t, banditCaptain, "Failed to get bandit-captain test data")

	// Use a non-existent monster key for error testing
	nonExistentMonsterKey := "monster_non_existent"

	testCases := []struct {
		name        string
		monsterKey  string
		xpValues    map[string]int
		expectedXP  int
		expectError bool
	}{
		{
			name:       "Valid monster key returns XP",
			monsterKey: "goblin",
			xpValues: map[string]int{
				"goblin":         goblin.XP,
				"bandit-captain": banditCaptain.XP,
			},
			expectedXP:  goblin.XP,
			expectError: false,
		},
		{
			name:       "Another valid monster key returns XP",
			monsterKey: "bandit-captain",
			xpValues: map[string]int{
				"goblin":         goblin.XP,
				"bandit-captain": banditCaptain.XP,
			},
			expectedXP:  banditCaptain.XP,
			expectError: false,
		},
		{
			name:       "Non-existent monster key returns error",
			monsterKey: nonExistentMonsterKey,
			xpValues: map[string]int{
				"goblin":         goblin.XP,
				"bandit-captain": banditCaptain.XP,
			},
			expectedXP:  0,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test repository with real XP values
			repo := &TestMonsterRepository{
				xpValues: tc.xpValues,
			}

			// Call the method
			xp, err := repo.GetMonsterXP(tc.monsterKey)

			// Assert results
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedXP, xp)
			}
		})
	}
}
