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

### Entity Placement Details

For each entity type, you can specify whether to place it at a random position or at a specific position:

```go
// Random placement (position will be determined automatically)
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
```

**Important notes about entity placement:**

1. For grid-based rooms:
   - Random placement will find an empty cell in the grid
   - Fixed placement will fail if the specified position is occupied or out of bounds
   - The `PlaceEntity` function ensures entities don't overlap

2. For gridless rooms:
   - Random placement assigns a position within room bounds (no collision detection)
   - Fixed placement only validates that the position is within room bounds
   - Multiple entities can occupy the same position

3. The `Count` field in entity configs:
   - Each config creates exactly one entity
   - To create multiple entities of the same type, create multiple configs

### Moving Entities

The library provides functions to move entities within a room:

```go
// Get a monster from the room
monster := room.Monsters[0]

// Define a new position
newPosition := entities.Position{X: 3, Y: 4}

// Move the monster to the new position
err := services.MovePlaceable(room, &monster, newPosition)
if err != nil {
    // Handle error (position occupied, out of bounds, etc.)
}

// The monster's position is now updated both in the room's slice and in the original variable
fmt.Printf("Monster position: (%d, %d)\n", monster.Position.X, monster.Position.Y)
```

### Removing Entities

To remove entities from a room:

```go
// Remove a monster from the room
monster := room.Monsters[0]
removed, err := services.RemovePlaceable(room, &monster)
if err != nil {
    // Handle error
}
if removed {
    fmt.Println("Monster was successfully removed")
}
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

// Optional: Configure items
itemConfigs := []services.ItemConfig{
    {
        Name:        "Healing Potion",
        Key:         "item_potion_healing",
        Value:       50,
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

### Encounter Balancing

The library includes a balancer that can adjust monster selection based on party composition and desired difficulty:

```go
// Create a balancer
balancer := services.NewBalancer()

// Create a party
party := entities.Party{
    Members: []entities.PartyMember{
        {Name: "Aragorn", Level: 5},
        {Name: "Legolas", Level: 4},
        {Name: "Gimli", Level: 5},
    },
}

// Calculate target CR for this party at Hard difficulty
targetCR, err := balancer.CalculateTargetCR(party, entities.EncounterDifficultyHard)
if err != nil {
    // Handle error
}
fmt.Printf("Target CR: %.2f\n", targetCR)

// Determine the difficulty of an encounter with specific monsters
monsters := []entities.Monster{
    {Name: "Goblin", CR: 0.25},
    {Name: "Goblin", CR: 0.25},
    {Name: "Hobgoblin", CR: 0.5},
    {Name: "Bugbear", CR: 1.0},
}
difficulty, err := balancer.DetermineEncounterDifficulty(monsters, party)
if err != nil {
    // Handle error
}
fmt.Printf("Encounter difficulty: %s\n", difficulty)

// Adjust monster selection to match target difficulty
adjustedMonsters, err := balancer.AdjustMonsterSelection(monsters, party, entities.EncounterDifficultyMedium)
if err != nil {
    // Handle error
}
fmt.Printf("Adjusted monster count: %d\n", len(adjustedMonsters))
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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

None
