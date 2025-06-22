package testutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
)

// Equipment represents the structure of equipment data from the D&D 5e API
type Equipment struct {
	Index             string `json:"index"`
	Name              string `json:"name"`
	EquipmentCategory struct {
		Index string `json:"index"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"equipment_category"`
	Cost struct {
		Quantity int    `json:"quantity"`
		Unit     string `json:"unit"`
	} `json:"cost"`
	Weight int      `json:"weight"`
	Desc   []string `json:"desc"`
}

// EquipmentList represents the structure of the equipment list from the D&D 5e API
type EquipmentList struct {
	Count   int `json:"count"`
	Results []struct {
		Index string `json:"index"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"results"`
}

// GetTestDataDir returns the absolute path to the testdata directory
func GetTestDataDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}
	dir := filepath.Dir(filename)
	// Go up two directories from internal/testutil to reach project root, then to testdata
	return filepath.Join(filepath.Dir(filepath.Dir(dir)), "internal", "testdata"), nil
}

// LoadAllEquipment loads all equipment data from the test files
func LoadAllEquipment() (map[string]*entities.Item, error) {
	// Get the absolute path to the testdata directory
	testdataDir, err := GetTestDataDir()
	if err != nil {
		return nil, err
	}

	// Path to the equipment list file
	equipmentListPath := filepath.Join(testdataDir, "equipment", "equipmentlist.json")

	// Read the equipment list file
	equipmentListData, err := os.ReadFile(equipmentListPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read equipment list file: %w", err)
	}

	// Parse the equipment list
	var equipmentList EquipmentList
	err = json.Unmarshal(equipmentListData, &equipmentList)
	if err != nil {
		return nil, fmt.Errorf("failed to parse equipment list: %w", err)
	}

	// Create a map to store items
	items := make(map[string]*entities.Item)

	// Load each equipment item
	for _, result := range equipmentList.Results {
		key := result.Index
		item, err := LoadEquipment(key)
		if err != nil {
			// Skip items that don't have a corresponding file
			continue
		}

		// Store the item with different key formats for easy lookup
		items[key] = item
		items["item_"+key] = item
		if len(item.Name) > 0 {
			capitalizedKey := "item_" + item.Name
			items[capitalizedKey] = item
		}
	}

	return items, nil
}

// LoadEquipment loads a specific equipment item from the test files
func LoadEquipment(key string) (*entities.Item, error) {
	// Get the absolute path to the testdata directory
	testdataDir, err := GetTestDataDir()
	if err != nil {
		return nil, err
	}

	// Path to the equipment file
	equipmentPath := filepath.Join(testdataDir, "equipment", key+".json")

	// Read the equipment file
	equipmentData, err := os.ReadFile(equipmentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read equipment file: %w", err)
	}

	// Parse the equipment data
	var equipment Equipment
	err = json.Unmarshal(equipmentData, &equipment)
	if err != nil {
		return nil, fmt.Errorf("failed to parse equipment data: %w", err)
	}

	// Convert to an entity.Item
	item := &entities.Item{
		Key:      equipment.Index,
		Name:     equipment.Name,
		Type:     equipment.EquipmentCategory.Index,
		Category: equipment.EquipmentCategory.Index, // Set Category to match Type
		// Add other fields as needed
	}

	return item, nil
}

// CreateTestItemRepository creates a test item repository with real item data
func CreateTestItemRepository() (map[string]*entities.Item, error) {
	return LoadAllEquipment()
}
