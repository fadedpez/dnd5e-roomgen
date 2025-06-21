package repositories

// MonsterRepository defines the interface for accessing monster data
type MonsterRepository interface {
	// GetMonsterXP returns the XP value for a monster based on its key
	GetMonsterXP(monsterKey string) (int, error)
}
