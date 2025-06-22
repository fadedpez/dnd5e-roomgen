package repositories

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/fadedpez/dnd5e-api/clients/dnd5e"
	apientities "github.com/fadedpez/dnd5e-api/entities"
	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
	"github.com/google/uuid"
)

// APIItemRepository implements the ItemRepository interface using the DND5e API
type APIItemRepository struct {
	apiClient dnd5e.Interface
}

// NewAPIItemRepository creates a new APIItemRepository
func NewAPIItemRepository() (*APIItemRepository, error) {
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

	return &APIItemRepository{
		apiClient: apiClient,
	}, nil
}

// GetItemByKey fetches an item by its key
func (r *APIItemRepository) GetItemByKey(key string) (*entities.Item, error) {
	// Call the API to get the equipment data
	apiEquipment, err := r.apiClient.GetEquipment(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from API: %w", err)
	}

	// Create base item with common fields
	item := &entities.Item{
		ID:   uuid.NewString(),
		Key:  key,
		Type: apiEquipment.GetType(),
	}

	// Handle different equipment types
	switch equip := apiEquipment.(type) {
	case *apientities.Equipment:
		// Basic equipment
		item.Name = equip.Name
		if equip.EquipmentCategory != nil {
			item.Category = equip.EquipmentCategory.Key
		}
		if equip.Cost != nil {
			item.Value = equip.Cost.Quantity
			item.ValueUnit = equip.Cost.Unit
		}
		item.Weight = int(equip.Weight)

	case *apientities.Weapon:
		// Weapon
		item.Name = equip.Name
		if equip.EquipmentCategory != nil {
			item.Category = equip.EquipmentCategory.Key
		}
		if equip.Cost != nil {
			item.Value = equip.Cost.Quantity
			item.ValueUnit = equip.Cost.Unit
		}
		item.Weight = int(equip.Weight)

		// Weapon-specific properties
		if equip.Properties != nil {
			properties := make([]string, 0, len(equip.Properties))
			for _, prop := range equip.Properties {
				if prop != nil {
					properties = append(properties, prop.Key)
				}
			}
			item.Properties = properties
		}

		// Damage information
		if equip.Damage != nil {
			item.DamageDice = equip.Damage.DamageDice
			if equip.Damage.DamageType != nil {
				item.DamageType = equip.Damage.DamageType.Key
			}
		}

	case *apientities.Armor:
		// Armor
		item.Name = equip.Name
		if equip.EquipmentCategory != nil {
			item.Category = equip.EquipmentCategory.Key
		}
		if equip.Cost != nil {
			item.Value = equip.Cost.Quantity
			item.ValueUnit = equip.Cost.Unit
		}
		item.Weight = int(equip.Weight)

		// Armor-specific properties
		if equip.ArmorClass != nil {
			item.ArmorClass = equip.ArmorClass.Base
		}
		item.StealthDisadvantage = equip.StealthDisadvantage
	}

	return item, nil
}

// GetRandomItems fetches a specified number of random items
func (r *APIItemRepository) GetRandomItems(count int) ([]*entities.Item, error) {
	// Call the API to get all equipment
	allEquipment, err := r.apiClient.ListEquipment()
	if err != nil {
		return nil, fmt.Errorf("failed to list equipment from API: %w", err)
	}

	// Shuffle and select random equipment
	rand.Shuffle(len(allEquipment), func(i, j int) {
		allEquipment[i], allEquipment[j] = allEquipment[j], allEquipment[i]
	})

	// Limit to requested count
	selectedCount := count
	if selectedCount > len(allEquipment) {
		selectedCount = len(allEquipment)
	}
	selectedEquipment := allEquipment[:selectedCount]

	// Convert to our entities
	items := make([]*entities.Item, 0, selectedCount)
	for _, equipRef := range selectedEquipment {
		item, err := r.GetItemByKey(equipRef.Key)
		if err != nil {
			// Skip items that fail to load
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

// GetRandomItemsByCategory fetches random items of a specific category
func (r *APIItemRepository) GetRandomItemsByCategory(category string, count int) ([]*entities.Item, error) {
	// Call the API to get all equipment
	allEquipment, err := r.apiClient.ListEquipment()
	if err != nil {
		return nil, fmt.Errorf("failed to list equipment from API: %w", err)
	}

	// Filter by category - we'll need to fetch each item to check its category
	categoryEquipment := make([]*apientities.ReferenceItem, 0)
	for _, equipRef := range allEquipment {
		// Get the full equipment details to check category
		fullEquip, err := r.apiClient.GetEquipment(equipRef.Key)
		if err != nil {
			continue
		}

		// Check if the equipment category matches the requested category
		switch equip := fullEquip.(type) {
		case *apientities.Equipment:
			if equip.EquipmentCategory != nil &&
				strings.EqualFold(equip.EquipmentCategory.Key, category) {
				categoryEquipment = append(categoryEquipment, equipRef)
			}
		case *apientities.Weapon:
			if equip.EquipmentCategory != nil &&
				strings.EqualFold(equip.EquipmentCategory.Key, category) {
				categoryEquipment = append(categoryEquipment, equipRef)
			}
		case *apientities.Armor:
			if equip.EquipmentCategory != nil &&
				strings.EqualFold(equip.EquipmentCategory.Key, category) {
				categoryEquipment = append(categoryEquipment, equipRef)
			}
		}
	}

	// Shuffle and select random equipment
	rand.Shuffle(len(categoryEquipment), func(i, j int) {
		categoryEquipment[i], categoryEquipment[j] = categoryEquipment[j], categoryEquipment[i]
	})

	// Limit to requested count
	selectedCount := count
	if selectedCount > len(categoryEquipment) {
		selectedCount = len(categoryEquipment)
	}
	selectedEquipment := categoryEquipment[:selectedCount]

	// Convert to our entities
	items := make([]*entities.Item, 0, selectedCount)
	for _, equipRef := range selectedEquipment {
		item, err := r.GetItemByKey(equipRef.Key)
		if err != nil {
			// Skip items that fail to load
			continue
		}
		items = append(items, item)
	}

	return items, nil
}
