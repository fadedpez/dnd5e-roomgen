// Package repositories contains data access interfaces and implementations
package repositories

import (
	"fmt"
	"net/http"

	"github.com/fadedpez/dnd5e-api/clients/dnd5e"
)

// APIMonsterRepository implements MonsterRepository using the dnd5e-api
type APIMonsterRepository struct {
	apiClient dnd5e.Interface
}

// NewAPIMonsterRepository creates a new APIMonsterRepository
func NewAPIMonsterRepository() (*APIMonsterRepository, error) {
	// Create the HTTP client
	httpClient := &http.Client{}

	// Create the API client config
	config := &dnd5e.DND5eAPIConfig{
		Client: httpClient,
	}

	// Initialize the API client
	apiClient, err := dnd5e.NewDND5eAPI(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create DND5e API client: %w", err)
	}

	return &APIMonsterRepository{
		apiClient: apiClient,
	}, nil
}

// GetMonsterXP returns the XP value for a monster based on its key
func (r *APIMonsterRepository) GetMonsterXP(monsterKey string) (int, error) {
	// Use the API client to fetch monster data
	monster, err := r.apiClient.GetMonster(monsterKey)
	if err != nil {
		return 0, fmt.Errorf("failed to get monster data: %w", err)
	}

	// Return the XP value from the monster data
	return monster.XP, nil
}
