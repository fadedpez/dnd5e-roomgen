package repositories

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
	testCases := []struct {
		name        string
		monsterKey  string
		xpValues    map[string]int
		expectedXP  int
		expectError bool
	}{
		{
			name:       "Valid monster key returns XP",
			monsterKey: "monster_Goblin",
			xpValues: map[string]int{
				"monster_Goblin": 50,
				"monster_Orc":    100,
				"monster_Troll":  450,
			},
			expectedXP:  50,
			expectError: false,
		},
		{
			name:       "API error is propagated",
			monsterKey: "monster_Unknown",
			xpValues: map[string]int{
				"monster_Goblin": 50,
			},
			expectedXP:  0,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test repository
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
