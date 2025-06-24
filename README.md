# DnD 5e Room Generator

A Go library for generating dynamic rooms for Dungeons & Dragons 5th Edition adventures.

## Features

- Generate room layouts with configurable dimensions and light levels
- Place monsters with random or specific positioning
- Place players with random or specific positioning
- Place NPCs with random or specific positioning
- Place obstacles with random or specific positioning
- Place items with random or specific positioning
- Generate treasure rooms with scaled loot based on party size and difficulty
- Automatic grid initialization for spatial tracking
- Support for gridless rooms when spatial positioning is not needed
- Flexible service layer for easy integration with applications
- Support for encounter balancing based on party composition
- Prioritized entity placement (players first, then monsters, then items)
- Graceful handling of placement failures

## Installation

```bash
go get github.com/fadedpez/dnd5e-roomgen
```

## Architecture

The library follows a clean layered architecture:

- **Entities Layer**: Core domain objects like Room, Monster, Player, Position, etc.
- **Services Layer**: Business logic for room generation, entity placement, encounter balancing, and room utilities

## Basic Usage

### Generating a Room

The library supports two types of rooms: grid-based and gridless. The choice depends on your application's needs.

```go
import (
    "github.com/fadedpez/dnd5e-roomgen/internal/entities"
    "github.com/fadedpez/dnd5e-roomgen/internal/services"
)

// Create a room service
roomService, err := services.NewRoomService()
if err != nil {
    // Handle error
}
```

#### Grid-Based Rooms

Grid-based rooms are ideal when you need:
- Spatial positioning and distance calculations
- Cell-based occupancy tracking
- Position validation to prevent entity overlap

```go
// Configure a grid-based room
roomConfig := services.RoomConfig{
    Width:       15,
    Height:      10,
    LightLevel:  entities.LightLevelDim,
    Description: "A dimly lit dungeon chamber",
    UseGrid:     true,  // Enable grid for this room
}

// Generate the room
room, err := roomService.GenerateRoom(roomConfig)
if err != nil {
    // Handle error
}

// With a grid-based room:
// - Entities can only be placed on empty cells
// - Position validation prevents overlap
// - Grid cells track occupancy
```

#### Gridless Rooms

Gridless rooms are useful when:
- Spatial positioning is not important
- You only need to track entity existence
- You want simplified room management
- You need better performance for rooms with many entities

```go
// Configure a gridless room
roomConfig := services.RoomConfig{
    Width:       15,
    Height:      10,
    LightLevel:  entities.LightLevelDim,
    Description: "A dimly lit dungeon chamber",
    UseGrid:     false,  // Disable grid for this room
}

// Generate the gridless room
room, err := roomService.GenerateRoom(roomConfig)
if err != nil {
    // Handle error
}

// With a gridless room:
// - Entities still have positions assigned (within room bounds)
// - No position validation or overlap prevention
// - No grid cell occupancy tracking
// - Positions can be used for visual representation if needed
```

### Adding Entities to a Room

The library supports adding various entity types to a room:

```go
// Create configurations for different entity types
monsterConfigs := []services.MonsterConfig{...}
playerConfigs := []services.PlayerConfig{...}
npcConfigs := []services.NPCConfig{...}
obstacleConfigs := []services.ObstacleConfig{...}
itemConfigs := []services.ItemConfig{...}

// Convert to generic PlaceableConfig interface
var placeables []services.PlaceableConfig
for _, config := range monsterConfigs {
    placeables = append(placeables, config)
}
for _, config := range playerConfigs {
    placeables = append(placeables, config)
}
for _, config := range npcConfigs {
    placeables = append(placeables, config)
}
for _, config := range obstacleConfigs {
    placeables = append(placeables, config)
}
for _, config := range itemConfigs {
    placeables = append(placeables, config)
}

// Add all entities to the room
err := roomService.AddPlaceablesToRoom(room, placeables)
if err != nil {
    // Handle error
}
```

#### Random vs. Fixed Placement

Entities can be placed randomly or at specific positions:

```go
// Random placement (system finds an empty position)
monsterConfig := services.MonsterConfig{
    Name:        "Goblin",
    Key:         "monster_goblin",
    CR:          0.25,
    RandomPlace: true,  // Let the system find a position
}

// Fixed placement (position specified by you)
monsterConfig := services.MonsterConfig{
    Name:        "Goblin Chief",
    Key:         "monster_goblin_chief",
    CR:          1.0,
    RandomPlace: false,  // Use the specified position
    Position:    &entities.Position{X: 7, Y: 5},
}

// NPC with random placement
npcConfig := services.NPCConfig{
    Name:        "Merchant",
    Level:       3,
    Count:       1,
    RandomPlace: true,
}

// Obstacle with fixed placement
obstacleConfig := services.ObstacleConfig{
    Name:        "Stone Wall",
    Key:         "wall_stone",
    Blocking:    true,
    Count:       1,
    RandomPlace: false,
    Position:    &entities.Position{X: 5, Y: 3},
}
```

### Generating a Complete Room with Entities

For convenience, the library provides a method to generate and populate a room in one step:

```go
// Configure the room
roomConfig := services.RoomConfig{
    Width:       20,
    Height:      15,
    LightLevel:  entities.LightLevelBright,
    Description: "A well-lit chamber with stone walls",
    UseGrid:     true,  // Enable grid for this room
}

// Configure monsters
monsterConfigs := []services.MonsterConfig{
    {
        Name:        "Goblin",
        Key:         "monster_goblin",
        CR:          0.25,
        Count:       3,
        RandomPlace: true,
    },
    {
        Name:        "Bugbear",
        Key:         "monster_bugbear",
        CR:          1.0,
        RandomPlace: false,
        Position:    &entities.Position{X: 7, Y: 5},
    },
}

// Optional: Configure players
playerConfigs := []services.PlayerConfig{
    {
        Name:        "Aragorn",
        Level:       5,
        RandomPlace: true,
    },
    {
        Name:        "Gandalf",
        Level:       10,
        RandomPlace: false,
        Position:    &entities.Position{X: 2, Y: 2},
    },
}

// Configure NPCs
npcConfigs := []services.NPCConfig{
    {
        Name:        "Merchant",
        Level:       3,
        Count:       1,
        RandomPlace: true,
    },
}

// Configure obstacles
obstacleConfigs := []services.ObstacleConfig{
    {
        Name:        "Stone Wall",
        Key:         "wall_stone",
        Blocking:    true,
        Count:       2,
        RandomPlace: true,
    },
    {
        Name:        "Barricade",
        Key:         "barricade_wooden",
        Blocking:    true,
        Count:       1,
        RandomPlace: false,
        Position:    &entities.Position{X: 10, Y: 7},
    },
}

// Configure items
itemConfigs := []services.ItemConfig{
    {
        Name:        "Health Potion",
        Key:         "item_potion_health",
        Value:       50,
        Count:       2,
        RandomPlace: true,
    },
    {
        Name:        "Magic Sword",
        Key:         "item_sword_magic",
        Value:       500,
        RandomPlace: false,
        Position:    &entities.Position{X: 12, Y: 8},
    },
}

// Optional: Provide party for automatic encounter balancing
party := &entities.Party{
    Members: []entities.PartyMember{
        {Name: "Aragorn", Level: 5},
        {Name: "Legolas", Level: 4},
        {Name: "Gimli", Level: 5},
    },
}

// Optional: Specify encounter difficulty (defaults to Medium if not provided)
difficulty := entities.EncounterDifficultyHard

// Generate and populate room in one step with automatic balancing
// IMPORTANT: At least one of monsterConfigs, playerConfigs, npcConfigs, obstacleConfigs, or itemConfigs must be non-empty
room, err := roomService.GenerateAndPopulateRoom(
    roomConfig,
    monsterConfigs,   // Pass empty slice ([]services.MonsterConfig{}) if not needed
    playerConfigs,    // Pass empty slice ([]services.PlayerConfig{}) if not needed
    npcConfigs,       // Pass empty slice ([]services.NPCConfig{}) if not needed
    obstacleConfigs,  // Pass empty slice ([]services.ObstacleConfig{}) if not needed
    itemConfigs,      // Pass empty slice ([]services.ItemConfig{}) if not needed
    party,            // Pass nil if no balancing is needed
    difficulty,       // Only used if party is provided
)
if err != nil {
    // Handle error
}

// Access room properties
fmt.Printf("Room: %dx%d, %s\n", room.Width, room.Height, room.Description)
fmt.Printf("Light level: %s\n", room.LightLevel)

// Access entity properties
fmt.Printf("Players: %d\n", len(room.Players))
fmt.Printf("Monsters: %d\n", len(room.Monsters))
fmt.Printf("NPCs: %d\n", len(room.NPCs))
fmt.Printf("Obstacles: %d\n", len(room.Obstacles))
fmt.Printf("Items: %d\n", len(room.Items))
```

### Managing Entities in a Room

```go
// Remove specific monsters from a room (e.g., after defeating them)
monsterIDs := []string{"monster-uuid-1", "monster-uuid-2"}
xpGained, notRemoved, err := roomService.CleanupRoom(room, entities.CellMonster, monsterIDs)
if err != nil {
    // Handle error
}
fmt.Printf("XP gained: %d\n", xpGained)
if len(notRemoved) > 0 {
    fmt.Printf("Some monsters could not be removed: %v\n", notRemoved)
}

// Remove all obstacles from a room (e.g., after breaking them)
_, notRemoved, err = roomService.CleanupRoom(room, entities.CellObstacle, []string{})
if err != nil {
    // Handle error
}
if len(notRemoved) > 0 {
    fmt.Printf("Some obstacles could not be removed: %v\n", notRemoved)
}

// Remove all items from a room (e.g., after collecting all treasure)
_, notRemoved, err = roomService.CleanupRoom(room, entities.CellItem, []string{})
if err != nil {
    // Handle error
}
if len(notRemoved) > 0 {
    fmt.Printf("Some items could not be removed: %v\n", notRemoved)
}

// Remove specific players from a room
playerIDs := []string{"player-uuid-1"}
_, notRemoved, err = roomService.CleanupRoom(room, entities.CellPlayer, playerIDs)
if err != nil {
    // Handle error
}

// Full room cleanup - remove all entities of all types
// This is useful when transitioning to a new room or resetting a room
entityTypes := []entities.CellType{
    entities.CellMonster,
    entities.CellNPC,
    entities.CellObstacle,
    entities.CellItem,
    // Note: You might want to handle players separately
}

totalXP := 0
for _, entityType := range entityTypes {
    xp, _, err := roomService.CleanupRoom(room, entityType, []string{})
    if err != nil {
        // Handle error
    }
    totalXP += xp
}
fmt.Printf("Total XP gained from cleanup: %d\n", totalXP)
```

### Working with NPC Inventories

NPCs in the library can have inventories, making them useful as merchants, traders, or quest-givers:

```go
// Create an NPC with inventory items
npcConfig := services.NPCConfig{
    Name:        "Merchant",
    Level:       3,
    RandomPlace: true,
}

// Add the NPC to the room
placeables := []services.PlaceableConfig{npcConfig}
err := roomService.AddPlaceablesToRoom(room, placeables)
if err != nil {
    // Handle error
}

// Get the created NPC from the room
merchant := room.NPCs[0]

// Create items to add to the merchant's inventory
potion := entities.Item{
    ID:    "item-uuid-1",
    Name:  "Health Potion",
    Key:   "item_potion_health",
    Value: 50,
}

sword := entities.Item{
    ID:    "item-uuid-2",
    Name:  "Magic Sword",
    Key:   "item_sword_magic",
    Value: 500,
}

// Add items to the merchant's inventory
merchant.AddItemToInventory(potion)
merchant.AddItemToInventory(sword)

// View the merchant's inventory
inventory := merchant.GetInventory()
fmt.Printf("Merchant has %d items for sale\n", len(inventory))
for _, item := range inventory {
    fmt.Printf("- %s: %d gold\n", item.Name, item.Value)
}

// Player buys an item from the merchant
boughtItem, success := merchant.RemoveItemFromInventory("item-uuid-1")
if success {
    fmt.Printf("Bought %s for %d gold\n", boughtItem.Name, boughtItem.Value)
    
    // Add the item to the room (e.g., player picks it up)
    room.Items = append(room.Items, boughtItem)
    
    // Or add it directly to a player's inventory if you have a player struct with inventory
    // player.AddItemToInventory(boughtItem)
} else {
    fmt.Println("Item not found in merchant's inventory")
}

// Check updated inventory
inventory = merchant.GetInventory()
fmt.Printf("Merchant now has %d items for sale\n", len(inventory))
```

This inventory system allows for implementing:
- Shops and merchants with items for sale
- NPCs that can give or receive items as part of quests
- Loot that can be transferred between entities
- Trading systems between players and NPCs

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

None
