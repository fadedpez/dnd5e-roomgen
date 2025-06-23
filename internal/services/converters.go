package services

import (
	apientities "github.com/fadedpez/dnd5e-api/entities"
)

// ConvertAPIMonsterToConfig converts a monster from the DnD 5e API format to a MonsterConfig
// with the specified count.
func ConvertAPIMonsterToConfig(apiMonster *apientities.Monster, count int) *MonsterConfig {
	if count < 1 {
		count = 1 // Ensure at least one monster
	}

	return &MonsterConfig{
		Key:   apiMonster.Key,
		Name:  apiMonster.Name,
		Count: count,
		CR:    float64(apiMonster.ChallengeRating),
		// Additional fields can be mapped as needed
	}
}

// ConvertAPIMonsterSliceToConfigs converts a slice of monsters from the DnD 5e API format
// to a slice of MonsterConfigs, each with the specified count.
func ConvertAPIMonsterSliceToConfigs(apiMonsters []*apientities.Monster, count int) []*MonsterConfig {
	configs := make([]*MonsterConfig, len(apiMonsters))

	for i, monster := range apiMonsters {
		configs[i] = ConvertAPIMonsterToConfig(monster, count)
	}

	return configs
}

// ConvertAPIItemToConfig converts an item from the DnD 5e API format to an ItemConfig
// with the specified count.
func ConvertAPIItemToConfig(apiItem *apientities.Equipment, count int) *ItemConfig {
	if count < 1 {
		count = 1 // Ensure at least one item
	}

	return &ItemConfig{
		Key:         apiItem.Key,
		Name:        apiItem.Name,
		Count:       count,
		RandomPlace: true, // Default to random placement
		// Additional fields can be mapped as needed
	}
}

// ConvertAPIItemSliceToConfigs converts a slice of items from the DnD 5e API format
// to a slice of ItemConfigs, each with the specified count.
func ConvertAPIItemSliceToConfigs(apiItems []*apientities.Equipment, count int) []*ItemConfig {
	configs := make([]*ItemConfig, len(apiItems))

	for i, item := range apiItems {
		configs[i] = ConvertAPIItemToConfig(item, count)
	}

	return configs
}
