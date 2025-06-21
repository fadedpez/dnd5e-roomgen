package services

import (
	"fmt"
	"math"

	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
	"github.com/fadedpez/dnd5e-roomgen/internal/repositories"
)

// Balancer defines the interface for encounter balancing functionality
type Balancer interface {
	// DetermineEncounterDifficulty determines the difficulty of an encounter based on monsters and party
	DetermineEncounterDifficulty(monsters []entities.Monster, party entities.Party) (entities.EncounterDifficulty, error)

	// AdjustMonsterSelection adjusts the monster selection based on the party and desired difficulty
	AdjustMonsterSelection(monsterConfigs []MonsterConfig, party entities.Party, difficulty entities.EncounterDifficulty) ([]MonsterConfig, error)

	// CalculateTargetCR calculates the target CR for a party based on difficulty
	CalculateTargetCR(party entities.Party, difficulty entities.EncounterDifficulty) (float64, error)
}

// StandardBalancer implements the Balancer interface using D&D 5e rules
type StandardBalancer struct {
	monsterRepo repositories.MonsterRepository
}

// NewBalancer creates a new StandardBalancer
func NewBalancer(monsterRepo repositories.MonsterRepository) *StandardBalancer {
	return &StandardBalancer{
		monsterRepo: monsterRepo,
	}
}

// difficultyMultipliers maps difficulty levels to CR multipliers
var difficultyMultipliers = map[entities.EncounterDifficulty]float64{
	entities.EncounterDifficultyEasy:   0.5,  // Easy encounter: CR = 0.5 * party level
	entities.EncounterDifficultyMedium: 0.75, // Medium encounter: CR = 0.75 * party level
	entities.EncounterDifficultyHard:   1.0,  // Hard encounter: CR = 1.0 * party level
	entities.EncounterDifficultyDeadly: 1.5,  // Deadly encounter: CR = 1.5 * party level
}

// partySizeAdjustments maps party size to CR adjustments
var partySizeAdjustments = map[int]float64{
	1: 0.5,  // Solo player: reduce CR
	2: 0.75, // Two players: slightly reduce CR
	3: 1.0,  // Three players: standard CR
	4: 1.0,  // Four players: standard CR (baseline)
	5: 1.25, // Five players: increase CR
	6: 1.5,  // Six+ players: significantly increase CR
}

// CalculateTargetCR calculates the target CR for a party based on difficulty
func (b *StandardBalancer) CalculateTargetCR(party entities.Party, difficulty entities.EncounterDifficulty) (float64, error) {
	if party.Size() == 0 {
		return 0, fmt.Errorf("party cannot be empty")
	}

	// Get the difficulty multiplier
	difficultyMultiplier, ok := difficultyMultipliers[difficulty]
	if !ok {
		return 0, fmt.Errorf("invalid difficulty: %s", difficulty)
	}

	// Calculate the average party level
	avgLevel := party.AverageLevel()

	// Apply party size adjustment
	sizeAdjustment := 1.0 // Default for 3-4 players
	if adj, ok := partySizeAdjustments[party.Size()]; ok {
		sizeAdjustment = adj
	} else if party.Size() > 6 {
		sizeAdjustment = partySizeAdjustments[6] // Use the 6+ player adjustment
	}

	// Calculate the target CR
	targetCR := avgLevel * difficultyMultiplier * sizeAdjustment

	// Round to the nearest 0.25 for standard CR increments
	return math.Round(targetCR*4) / 4, nil
}

// calculateTotalCR calculates the total CR of a set of monsters
func calculateTotalCR(monsters []entities.Monster) float64 {
	totalCR := 0.0
	for _, monster := range monsters {
		totalCR += monster.CR
	}
	return totalCR
}

// DetermineEncounterDifficulty determines the difficulty of an encounter based on monsters and party
func (b *StandardBalancer) DetermineEncounterDifficulty(monsters []entities.Monster, party entities.Party) (entities.EncounterDifficulty, error) {
	if party.Size() == 0 {
		return "", fmt.Errorf("party cannot be empty")
	}

	// Calculate the total CR of the monsters
	totalCR := calculateTotalCR(monsters)

	// Calculate the average party level
	avgLevel := party.AverageLevel()

	// Apply party size adjustment
	sizeAdjustment := 1.0 // Default for 3-4 players
	if adj, ok := partySizeAdjustments[party.Size()]; ok {
		sizeAdjustment = adj
	} else if party.Size() > 6 {
		sizeAdjustment = partySizeAdjustments[6] // Use the 6+ player adjustment
	}

	// Calculate the adjusted CR ratio (totalCR / (avgLevel * sizeAdjustment))
	crRatio := totalCR / (avgLevel * sizeAdjustment)

	// Determine the difficulty based on the CR ratio
	if crRatio >= difficultyMultipliers[entities.EncounterDifficultyDeadly] {
		return entities.EncounterDifficultyDeadly, nil
	} else if crRatio >= difficultyMultipliers[entities.EncounterDifficultyHard] {
		return entities.EncounterDifficultyHard, nil
	} else if crRatio >= difficultyMultipliers[entities.EncounterDifficultyMedium] {
		return entities.EncounterDifficultyMedium, nil
	} else if crRatio >= difficultyMultipliers[entities.EncounterDifficultyEasy] {
		return entities.EncounterDifficultyEasy, nil
	}

	// If the CR ratio is below the easy threshold, it's trivial (we'll return easy)
	return entities.EncounterDifficultyEasy, nil
}

// AdjustMonsterSelection adjusts the monster selection based on the party and desired difficulty
func (b *StandardBalancer) AdjustMonsterSelection(monsterConfigs []MonsterConfig, party entities.Party, difficulty entities.EncounterDifficulty) ([]MonsterConfig, error) {
	if party.Size() == 0 {
		return nil, fmt.Errorf("party cannot be empty")
	}

	// Calculate the target CR for the encounter
	targetCR, err := b.CalculateTargetCR(party, difficulty)
	if err != nil {
		return nil, err
	}

	// Calculate the current total CR of the monster configs
	currentTotalCR := 0.0
	for _, config := range monsterConfigs {
		currentTotalCR += config.CR * float64(config.Count)
	}

	// If we're already close to the target CR (within 10%), return the original configs
	if math.Abs(currentTotalCR-targetCR)/targetCR < 0.1 {
		return monsterConfigs, nil
	}

	// Adjust the monster counts to get closer to the target CR
	adjustedConfigs := make([]MonsterConfig, len(monsterConfigs))
	copy(adjustedConfigs, monsterConfigs)

	// Calculate the scaling factor
	scalingFactor := targetCR / currentTotalCR

	// Apply the scaling factor to each monster count
	for i := range adjustedConfigs {
		// Calculate the new count based on the scaling factor
		newCount := int(math.Round(float64(adjustedConfigs[i].Count) * scalingFactor))

		// Ensure we have at least one monster if the original count was non-zero
		if adjustedConfigs[i].Count > 0 && newCount < 1 {
			newCount = 1
		}

		adjustedConfigs[i].Count = newCount
	}

	return adjustedConfigs, nil
}
