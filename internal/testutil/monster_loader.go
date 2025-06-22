package testutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
)

// MonsterData represents the structure of a monster JSON file
type MonsterData struct {
	Index           string  `json:"index"`
	Name            string  `json:"name"`
	ChallengeRating float64 `json:"challenge_rating"`
	XP              int     `json:"xp"`
}

// MonstersList represents the structure of the monsters list JSON file
type MonstersList struct {
	Count   int `json:"count"`
	Results []struct {
		Index string `json:"index"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"results"`
}

// getProjectRoot returns the absolute path to the project root directory
func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	// Go up two directories from internal/testutil to reach project root
	return filepath.Dir(filepath.Dir(dir))
}

// LoadMonster loads a single monster from the test data
func LoadMonster(monsterKey string) (*MonsterData, error) {
	// Normalize the key to match file naming conventions
	fileName := strings.ToLower(monsterKey)
	if !strings.HasSuffix(fileName, ".json") {
		fileName += ".json"
	}

	// Get absolute path to the test data file
	projectRoot := getProjectRoot()
	filePath := filepath.Join(projectRoot, "internal", "testdata", "monsters", fileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read monster file %s: %w", filePath, err)
	}

	var monster MonsterData
	if err := json.Unmarshal(data, &monster); err != nil {
		return nil, fmt.Errorf("failed to parse monster data: %w", err)
	}

	return &monster, nil
}

// LoadAllMonsters loads all available monsters from the test data
func LoadAllMonsters() (map[string]*MonsterData, error) {
	// Get absolute path to the monsters list file
	projectRoot := getProjectRoot()
	listPath := filepath.Join(projectRoot, "internal", "testdata", "monsters", "monsterslist.json")

	listData, err := os.ReadFile(listPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read monsters list: %w", err)
	}

	var monstersList MonstersList
	if err := json.Unmarshal(listData, &monstersList); err != nil {
		return nil, fmt.Errorf("failed to parse monsters list: %w", err)
	}

	// Now load each individual monster
	monsters := make(map[string]*MonsterData)
	for _, result := range monstersList.Results {
		// Extract the monster key from the URL
		key := filepath.Base(result.URL)

		// Try to load the monster data
		monster, err := LoadMonster(key)
		if err != nil {
			// Skip monsters that don't have individual files
			continue
		}

		// Store with both the original key and the monster_Key format
		monsters[key] = monster
		monsters["monster_"+key] = monster

		// Also store with capitalized name for backward compatibility
		if len(monster.Name) > 0 {
			capitalizedKey := "monster_" + monster.Name
			monsters[capitalizedKey] = monster
		}
	}

	return monsters, nil
}

// CreateTestMonsterRepository creates a test monster repository with real monster data
func CreateTestMonsterRepository() (map[string]int, error) {
	monsters, err := LoadAllMonsters()
	if err != nil {
		return nil, err
	}

	xpValues := make(map[string]int)
	for key, monster := range monsters {
		xpValues[key] = monster.XP
	}

	return xpValues, nil
}

// CreateEntityMonsters creates a slice of entity.Monster objects from test data
func CreateEntityMonsters() ([]*entities.Monster, error) {
	monsters, err := LoadAllMonsters()
	if err != nil {
		return nil, err
	}

	var entityMonsters []*entities.Monster
	for key, monster := range monsters {
		// Only use the keys that start with "monster_"
		if strings.HasPrefix(key, "monster_") {
			entityMonster := &entities.Monster{
				Key:  key,
				Name: monster.Name,
				CR:   monster.ChallengeRating,
			}
			entityMonsters = append(entityMonsters, entityMonster)
		}
	}

	return entityMonsters, nil
}
