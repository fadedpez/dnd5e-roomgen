package repositories

import (
	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
)

// ItemRepository defines the interface for retrieving item data
type ItemRepository interface {
	// GetItemByKey fetches an item by its key
	GetItemByKey(key string) (*entities.Item, error)

	// GetRandomItems fetches a specified number of random items
	GetRandomItems(count int) ([]*entities.Item, error)

	// GetRandomItemsByCategory fetches random items of a specific category
	GetRandomItemsByCategory(category string, count int) ([]*entities.Item, error)
}
