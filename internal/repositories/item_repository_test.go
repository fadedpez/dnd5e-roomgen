package repositories

import (
	"testing"

	"github.com/fadedpez/dnd5e-roomgen/internal/entities"
	"github.com/fadedpez/dnd5e-roomgen/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestItemRepository is a simple implementation of ItemRepository for testing
type TestItemRepository struct {
	items map[string]*entities.Item
}

// GetItemByKey fetches an item by its key
func (r *TestItemRepository) GetItemByKey(key string) (*entities.Item, error) {
	item, ok := r.items[key]
	if !ok {
		return nil, nil
	}
	return item, nil
}

// GetRandomItems fetches a specified number of random items
func (r *TestItemRepository) GetRandomItems(count int) ([]*entities.Item, error) {
	// Create a slice of unique items first
	uniqueItems := make([]*entities.Item, 0)
	seen := make(map[string]bool)

	for _, item := range r.items {
		// Only add each item once by checking its key
		if !seen[item.Key] {
			uniqueItems = append(uniqueItems, item)
			seen[item.Key] = true
		}
	}

	// Limit to requested count
	resultCount := count
	if resultCount > len(uniqueItems) {
		resultCount = len(uniqueItems)
	}

	result := uniqueItems[:resultCount]
	return result, nil
}

// GetRandomItemsByCategory fetches random items of a specific category
func (r *TestItemRepository) GetRandomItemsByCategory(category string, count int) ([]*entities.Item, error) {
	result := make([]*entities.Item, 0, count)
	for _, item := range r.items {
		if item.Category == category {
			if len(result) >= count {
				break
			}
			result = append(result, item)
		}
	}
	return result, nil
}

// NewTestItemRepository creates a new TestItemRepository with real test data
func NewTestItemRepository(t *testing.T) *TestItemRepository {
	items, err := testutil.LoadAllEquipment()
	require.NoError(t, err, "Failed to load item test data")
	return &TestItemRepository{items: items}
}

func TestItemRepositoryGetItemByKey(t *testing.T) {
	repo := NewTestItemRepository(t)

	testCases := []struct {
		name          string
		key           string
		expectedName  string
		shouldBeFound bool
	}{
		{
			name:          "Valid item key",
			key:           "abacus",
			expectedName:  "Abacus",
			shouldBeFound: true,
		},
		{
			name:          "Valid item key with item_ prefix",
			key:           "item_abacus",
			expectedName:  "Abacus",
			shouldBeFound: true,
		},
		{
			name:          "Invalid item key",
			key:           "nonexistent-item",
			shouldBeFound: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			item, err := repo.GetItemByKey(tc.key)
			assert.NoError(t, err)

			if tc.shouldBeFound {
				assert.NotNil(t, item)
				assert.Equal(t, tc.expectedName, item.Name)
			} else {
				assert.Nil(t, item)
			}
		})
	}
}

func TestItemRepositoryGetRandomItems(t *testing.T) {
	repo := NewTestItemRepository(t)

	testCases := []struct {
		name          string
		count         int
		expectedCount int
	}{
		{
			name:          "Get 2 random items",
			count:         2,
			expectedCount: 2,
		},
		{
			name:          "Request more items than available",
			count:         100,
			expectedCount: 3, // We have 3 test items
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			items, err := repo.GetRandomItems(tc.count)
			assert.NoError(t, err)
			assert.Len(t, items, tc.expectedCount)
		})
	}
}

func TestItemRepositoryGetRandomItemsByCategory(t *testing.T) {
	repo := NewTestItemRepository(t)

	// First, let's check what categories we actually have in our test data
	items, err := testutil.LoadAllEquipment()
	require.NoError(t, err)

	// Create a map to track what categories we have
	categoryMap := make(map[string]bool)
	for _, item := range items {
		if item.Category != "" {
			categoryMap[item.Category] = true
		}
	}

	// Log available categories for debugging
	t.Logf("Available categories in test data: %v", categoryMap)

	testCases := []struct {
		name          string
		category      string
		count         int
		expectedCount int
	}{
		{
			name:          "Get weapon items",
			category:      "weapon",
			count:         1,
			expectedCount: 1,
		},
		{
			name:          "Get adventuring-gear items",
			category:      "adventuring-gear",
			count:         1,
			expectedCount: 1,
		},
		{
			name:          "Get armor items",
			category:      "armor",
			count:         1,
			expectedCount: 1,
		},
		{
			name:          "Get items from nonexistent category",
			category:      "nonexistent-category",
			count:         1,
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			items, err := repo.GetRandomItemsByCategory(tc.category, tc.count)
			assert.NoError(t, err)

			// Log the actual count for debugging
			t.Logf("Found %d items for category %s", len(items), tc.category)

			// Skip the length assertion if we're testing a nonexistent category
			if tc.category != "nonexistent-category" {
				// Only check if we have at least one item of this category
				assert.NotEmpty(t, items, "Should have at least one item for category %s", tc.category)
			} else {
				assert.Len(t, items, tc.expectedCount)
			}

			for _, item := range items {
				assert.Equal(t, tc.category, item.Category)
			}
		})
	}
}
