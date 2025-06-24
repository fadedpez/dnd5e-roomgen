# DnD 5e Room Generator

A Go library for generating dynamic rooms for Dungeons & Dragons 5th Edition adventures.

## Features

- Generate room layouts with configurable dimensions and light levels
- Place monsters with random or specific positioning
- Place players with random or specific positioning
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
- **Repository Layer**: Data access for monsters, treasure, and other external resources

## Basic Usage

### Generating a Room

The library supports two types of rooms: grid-based and gridless. The choice depends on your application's needs.

```go
import (
    "github.com/fadedpez/dnd5e-roomgen/internal/entities"
    "github.com/fadedpez/dnd5e-roomgen/internal/services"
)

// Create a room service
roomService := services.NewRoomService()
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

The library provides a unified approach to adding entities to a room through the `AddPlaceablesToRoom` method:

```go
// Create configurations for different entity types
monsterConfigs := []services.MonsterConfig{...}
playerConfigs := []services.PlayerConfig{...}
itemConfigs := []services.ItemConfig{...}

// Convert to generic PlaceableConfig interface
var placeables []services.PlaceableConfig
for _, config := range monsterConfigs {
    placeables = append(placeables, config)
}
for _, config := range playerConfigs {
    placeables = append(placeables, config)
}
for _, config := range itemConfigs {
    placeables = append(placeables, config)
}

// Add all entities to the room with prioritized placement
// Players will be placed first, then monsters, then items
err := roomService.AddPlaceablesToRoom(room, placeables)
if err != nil {
    // Handle error - only occurs if player placement fails
    // Monster and item placement failures are handled gracefully
}
```

You can also use the type-specific convenience methods:

```go
// These methods use AddPlaceablesToRoom internally
err = roomService.AddMonstersToRoom(room, monsterConfigs)
err = roomService.AddPlayersToRoom(room, playerConfigs)
err = roomService.AddItemsToRoom(room, itemConfigs)
```

### Generating and Populating a Room in One Step

The `GenerateAndPopulateRoom` method provides a comprehensive way to create and populate a room with multiple entity types in a single call, with optional encounter balancing:

```go
// Configure room
roomConfig := services.RoomConfig{
    Width:       15,
    Height:      10,
    LightLevel:  entities.LightLevelDim,
    Description: "A dimly lit chamber with ancient runes on the walls",
    UseGrid:     true,
}

// All entity configurations are optional, but at least one type must be provided
// If you don't need a particular entity type, pass an empty slice or nil

// Optional: Configure monsters
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
        Count:       1,
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

// Optional: Configure items
itemConfigs := []services.ItemConfig{
    {
        Name:        "Healing Potion",
        Key:         "item_potion_healing",
        Value:       50,
        Count:       2,
        RandomPlace: true,
    },
    {
        Name:        "Magic Sword",
        Key:         "item_sword_magic",
        Value:       500,
        Count:       1,
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
// IMPORTANT: At least one of monsterConfigs, playerConfigs, or itemConfigs must be non-empty
room, err := roomService.GenerateAndPopulateRoom(
    roomConfig,
    monsterConfigs,  // Pass empty slice ([]services.MonsterConfig{}) if not needed
    playerConfigs,   // Pass empty slice ([]services.PlayerConfig{}) if not needed
    itemConfigs,     // Pass empty slice ([]services.ItemConfig{}) if not needed
    party,           // Pass nil if no balancing is needed
    difficulty,      // Only used if party is provided
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
fmt.Printf("Items: %d\n", len(room.Items))
```

This method handles several tasks in one operation:
1. Creates a new room based on the provided configuration
2. Optionally balances monster selection based on party and difficulty
3. Places entities in the room with proper prioritization
4. Handles placement failures gracefully

### Monster Configuration Examples

```go
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
        Name:        "Orc Chief",
        Key:         "monster_orc_chief",
        CR:          2.0,
        Count:       1,
        RandomPlace: false,
        Position:    &entities.Position{X: 7, Y: 5},
    },
}
```

### Player Configuration Examples

```go
// Configure players
playerConfigs := []services.PlayerConfig{
    {
        Name:        "Aragorn",
        Level:       5,
        RandomPlace: true, // Place player randomly in the room
    },
    {
        Name:        "Gandalf",
        Level:       10,
        RandomPlace: false, // Place player at a specific position
        Position:    &entities.Position{X: 2, Y: 2},
    },
}

// Access player properties
for i, player := range room.Players {
    fmt.Printf("Player %d: %s (Level %d) at position (%d,%d)\n",
        i+1, player.Name, player.Level, player.Position.X, player.Position.Y)
}
```

### Generating a Treasure Room

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

// Configure a room
roomConfig := services.RoomConfig{
    Width:       15,
    Height:      10,
    LightLevel:  entities.LightLevelBright,
    Description: "A glittering treasure chamber",
    UseGrid:     true,
}

// Create a party
party := entities.Party{
    Members: []entities.PartyMember{
        {Name: "Aragorn", Level: 5},
        {Name: "Legolas", Level: 4},
        {Name: "Gimli", Level: 5},
        {Name: "Gandalf", Level: 7},
    },
}

// Generate a treasure room with optional guardian monsters
includeGuardian := true
difficulty := entities.EncounterDifficultyMedium

room, err := roomService.PopulateRandomTreasureRoomWithParty(
    roomConfig, 
    party, 
    includeGuardian, 
    difficulty)
if err != nil {
    // Handle error
}

// Access room properties
fmt.Printf("Room: %dx%d, %s\n", room.Width, room.Height, room.Description)

// Access item properties
fmt.Printf("Found %d items:\n", len(room.Items))
for i, item := range room.Items {
    fmt.Printf("Item %d: %s at position (%d,%d)\n",
        i+1, item.Name, item.Position.X, item.Position.Y)
}

// If guardians were included, access monster properties
if includeGuardian {
    fmt.Printf("Guardian monsters (%d):\n", len(room.Monsters))
    for i, monster := range room.Monsters {
        fmt.Printf("Monster %d: %s (CR %.2f) at position (%d,%d)\n",
            i+1, monster.Name, monster.CR, monster.Position.X, monster.Position.Y)
    }
}
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
```

## Advanced Features

### Encounter Balancing

The library includes built-in encounter balancing based on party composition and desired difficulty:

```go
// Create a party
party := entities.Party{
    Members: []entities.PartyMember{
        {Name: "Aragorn", Level: 5},
        {Name: "Legolas", Level: 4},
        {Name: "Gimli", Level: 5},
    },
}

// Specify difficulty level
difficulty := entities.EncounterDifficultyHard

// Balance monsters for the party
monsterConfigs, err := roomService.BalanceMonsters(monsterConfigs, party, difficulty)
if err != nil {
    // Handle error
}

// Or use GenerateAndPopulateRoom with party and difficulty for automatic balancing
room, err := roomService.GenerateAndPopulateRoom(
    roomConfig,
    monsterConfigs,
    playerConfigs,
    itemConfigs,
    &party,
    difficulty,
)
```

### Placement Prioritization

When adding multiple entity types to a room, the library prioritizes placement in the following order:

1. Players (critical - will return error if placement fails)
2. Monsters (non-critical - will log warning if placement fails)
3. Items (non-critical - will log warning if placement fails)

This ensures that players are always placed in the room, even if it means some monsters or items cannot be placed due to space constraints.

### Moving Entities

The library provides a way to move entities within a room after they've been placed:

```go
// Get a reference to an entity
player := room.Players[0]
monster := room.Monsters[0]
item := room.Items[0]

// Create a new position
newPosition := entities.Position{X: 5, Y: 7}

// Move the entity (works with any entity type that implements Placeable)
err := roomService.MoveEntity(room, &player, newPosition)
if err != nil {
    // Handle error (position out of bounds, cell occupied, etc.)
}

// Move a monster
err = roomService.MoveEntity(room, &monster, newPosition)

// Move an item
err = roomService.MoveEntity(room, &item, newPosition)
```

The `MoveEntity` method:
- Works with any entity type that implements the `Placeable` interface
- Validates that the new position is within room bounds
- Ensures the target cell is unoccupied
- Updates both the entity's position and the room grid
- Works with both grid-based and gridless rooms

## License

[MIT License](LICENSE)
