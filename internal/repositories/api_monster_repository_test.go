package repositories

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Create a direct mock of the APIMonsterRepository for testing
type MockAPIMonsterRepository struct {
	GetMonsterXPFunc func(monsterKey string) (int, error)
}

// GetMonsterXP implements the MonsterRepository interface
func (m *MockAPIMonsterRepository) GetMonsterXP(monsterKey string) (int, error) {
	return m.GetMonsterXPFunc(monsterKey)
}

func TestAPIMonsterRepository_GetMonsterXP(t *testing.T) {
	testCases := []struct {
		name        string
		monsterKey  string
		setupMock   func() *MockAPIMonsterRepository
		expectedXP  int
		expectError bool
	}{
		{
			name:       "Valid monster key returns XP",
			monsterKey: "monster_Goblin",
			setupMock: func() *MockAPIMonsterRepository {
				mock := &MockAPIMonsterRepository{}
				mock.GetMonsterXPFunc = func(monsterKey string) (int, error) {
					if monsterKey == "monster_Goblin" {
						return 50, nil
					}
					return 0, errors.New("monster not found")
				}
				return mock
			},
			expectedXP:  50,
			expectError: false,
		},
		{
			name:       "API error is propagated",
			monsterKey: "monster_Unknown",
			setupMock: func() *MockAPIMonsterRepository {
				mock := &MockAPIMonsterRepository{}
				mock.GetMonsterXPFunc = func(monsterKey string) (int, error) {
					return 0, errors.New("monster not found")
				}
				return mock
			},
			expectedXP:  0,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock
			repo := tc.setupMock()

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
