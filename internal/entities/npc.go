package entities

// NPC represents a non-player character placed in the room
type NPC struct {
	ID        string   // UUID for this NPC instance
	Key       string   // Reference key from the API (if applicable)
	Name      string   // Name of the NPC
	Inventory []Item   // Items in the NPC's inventory
	Position  Position // Position of the NPC in the room (if grid is used)
}

// GetID returns the unique identifier for this NPC
func (n *NPC) GetID() string {
	return n.ID
}

// GetPosition returns the current position of this NPC in the room
func (n *NPC) GetPosition() Position {
	return n.Position
}

// SetPosition updates the position of this NPC
func (n *NPC) SetPosition(pos Position) {
	n.Position = pos
}

// GetCellType returns the type of cell this NPC occupies
func (n *NPC) GetCellType() CellType {
	return CellNPC
}

// AddItemToInventory adds an item to the NPC's inventory
func (n *NPC) AddItemToInventory(item Item) {
	n.Inventory = append(n.Inventory, item)
}

// RemoveItemFromInventory removes an item from the NPC's inventory by ID
// Returns the removed item and a boolean indicating success
func (n *NPC) RemoveItemFromInventory(itemID string) (Item, bool) {
	for i, item := range n.Inventory {
		if item.ID == itemID {
			// Remove the item from the inventory
			n.Inventory = append(n.Inventory[:i], n.Inventory[i+1:]...)
			return item, true
		}
	}
	return Item{}, false
}

// GetInventory returns all items in the NPC's inventory
func (n *NPC) GetInventory() []Item {
	return n.Inventory
}
