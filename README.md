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

## Installation

```bash
go get github.com/fadedpez/dnd5e-roomgen
```

## Architecture

The library follows a clean layered architecture:

- **Entities Layer**: Core domain objects like Room, Monster, Player, Position, etc.
- **Services Layer**: Business logic for room generation, entity placement, and encounter balancing
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

In both cases, the same entity placement interface is used, making your code consistent regardless of the room type:

```go
// These methods work the same for both grid-based and gridless rooms
err = roomService.AddMonstersToRoom(room, monsterConfigs)
err = roomService.AddPlayersToRoom(room, playerConfigs)
err = roomService.AddItemsToRoom(room, itemConfigs)
```

### Adding Monsters to a Room

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

// Add monsters to the room
err = roomService.AddMonstersToRoom(room, monsterConfigs)
if err != nil {
    // Handle error
}
```

### Adding Players to a Room

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

// Add players to the room
err = roomService.AddPlayersToRoom(room, playerConfigs)
if err != nil {
    // Handle error
}

// Access player properties
for i, player := range room.Players {
    fmt.Printf("Player %d: %s (Level %d) at position (%d,%d)\n",
        i+1, player.Name, player.Level, player.Position.X, player.Position.Y)
}
```

### Generating a Room with Monsters in One Step

```go
// Generate a room and add monsters in one step
room, err := roomService.PopulateRoomWithMonsters(roomConfig, monsterConfigs)
if err != nil {
    // Handle error
}

// Access room properties
fmt.Printf("Room: %dx%d, %s\n", room.Width, room.Height, room.Description)
fmt.Printf("Light level: %s\n", room.LightLevel)

// Access monster properties
for i, monster := range room.Monsters {
    fmt.Printf("Monster %d: %s (CR %.2f) at position (%d,%d)\n",
        i+1, monster.Name, monster.CR, monster.Position.X, monster.Position.Y)
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

### Removing Players from a Room

```go
// Remove a player by ID
removed, err := entities.RemovePlayer(room, playerID)
if err != nil {
    // Handle error
}

if removed {
    fmt.Println("Player successfully removed")
} else {
    fmt.Println("Player not found")
}
```

### Cleaning Up a Room and Calculating XP

```go
// Clean up a room (remove monsters) and calculate XP
totalXP, notRemovedIDs, err := roomService.CleanupRoom(room, []string{})  // Empty slice removes all monsters
if err != nil {
    // Handle error
}

fmt.Printf("Total XP earned: %d\n", totalXP)
```

### API Integration

For detailed information on integrating with the DnD 5e API, including:

- Converting API entities to room service configurations
- Using the encounter balancer with API monsters
- Player and item configuration with API integration
- Generating treasure rooms with API monsters
- Best practices for API integration

Please see the [API Integration Guide](docs/API_INTEGRATION.md).

### Using the Monster Repository

For local testing and non-API usage:

```go
import (
    "github.com/fadedpez/dnd5e-roomgen/internal/repositories"
)

// Create a test monster repository for unit testing
testRepo := &repositories.TestMonsterRepository{
    xpValues: map[string]int{
        "monster_goblin": 50,
        "monster_orc": 100,
    },
}

// Use the test repository in your service
roomService := services.NewRoomService(services.WithMonsterRepository(testRepo))
```

For API integration, see the [API Integration Guide](docs/API_INTEGRATION.md).

### Using the Encounter Balancer

The library includes a powerful encounter balancing system that helps create appropriately challenging encounters based on party composition and desired difficulty level.

For detailed examples of using the encounter balancer with the DnD 5e API, see the [API Integration Guide](docs/API_INTEGRATION.md).

#### Integrated Room Service Balancing

For most use cases, you can use the RoomService's integrated balancing methods:

```go
import (
    "github.com/fadedpez/dnd5e-roomgen/internal/entities"
    "github.com/fadedpez/dnd5e-roomgen/internal/services"
)

// Create a room service (automatically initializes the balancer)
roomService, err := services.NewRoomService()
if err != nil {
    // Handle error
}

// Create a party
party := entities.Party{
    Members: []entities.PartyMember{
        {Name: "Aragorn", Level: 5},
        {Name: "Legolas", Level: 5},
        {Name: "Gimli", Level: 5},
        {Name: "Gandalf", Level: 7},
    },
}

// Define monster configurations
monsterConfigs := []services.MonsterConfig{
    {
        Name:        "Goblin",
        Key:         "monster_goblin",
        CR:          0.25,
        Count:       4,
        RandomPlace: true,
    },
    {
        Name:        "Bugbear",
        Key:         "monster_bugbear",
        CR:          1.0,
        Count:       1,
        RandomPlace: true,
    },
}

// Method 1: Balance monster configurations without creating a room
balancedConfigs, err := roomService.BalanceMonsterConfigs(monsterConfigs, party, entities.EncounterDifficultyHard)
if err != nil {
    // Handle error
}
fmt.Printf("Adjusted monster counts for hard difficulty: %+v\n", balancedConfigs)

// Method 2: Generate a room with automatically balanced monsters in one step
room, err := roomService.PopulateRoomWithBalancedMonsters(roomConfig, monsterConfigs, party, entities.EncounterDifficultyHard)
if err != nil {
    // Handle error
}
```

#### Balancer Use Cases

1. **Creating Balanced Encounters**:
   - Calculate appropriate challenge ratings for your party
   - Automatically adjust monster counts to match desired difficulty

2. **Analyzing Encounter Difficulty**:
   - Determine if an existing encounter is Easy, Medium, Hard, or Deadly
   - Validate encounter designs against party composition

3. **Dynamic Encounter Scaling**:
   - Scale encounters up or down based on party size and level
   - Maintain appropriate challenge as party composition changes

4. **Difficulty Customization**:
   - Choose from standard D&D 5e difficulty levels (Easy, Medium, Hard, Deadly)
   - Apply consistent difficulty calculations across your application

5. **One-Step Room Generation with Balanced Encounters**:
   - Generate complete rooms with monsters balanced for your party
   - Streamline encounter creation while maintaining appropriate challenge

## Advanced Features

### Gridless Rooms

The library supports gridless rooms for applications that don't need spatial positioning or distance calculations:

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

// Add entities to the gridless room
// Entities will still have positions assigned but no grid validation will occur
err = roomService.AddMonstersToRoom(room, monsterConfigs)
if err != nil {
    // Handle error
}

// The placement interface handles gridless rooms automatically
// No need for special configuration in entity addition methods
```

When using gridless rooms:

1. Entities are still added to their respective slices (`Monsters`, `Players`, `Items`)
2. Positions are still assigned but not validated against a grid
3. No grid cell occupancy tracking occurs
4. The `FindEmptyPosition` function returns random positions within room bounds
5. All placement interface methods work seamlessly with both grid and gridless rooms

Gridless rooms are useful for:
- Applications that only need to track entity existence, not positions
- Scenarios where distance calculations are not needed
- Simplified room management without spatial constraints
- Improved performance for large rooms with many entities

## Customization and Advanced Configuration

The DnD 5e Room Generator library is designed to be flexible and extensible. This section covers how to customize the library for your specific needs.

### Custom Configuration Options

#### Complete RoomConfig Options

The `RoomConfig` struct provides several options for customizing room generation:

```go
roomConfig := services.RoomConfig{
    Width:       15,              // Width of the room in grid units
    Height:      10,              // Height of the room in grid units
    LightLevel:  entities.LightLevelDim, // Light level: LightLevelBright, LightLevelDim, or LightLevelDark
    Description: "A dusty chamber with cobwebs in the corners", // Optional room description
    RoomType:    "dungeon",       // Optional room type for categorization
    UseGrid:     true,            // Whether to use a grid for spatial tracking
}
```

#### Complete MonsterConfig Options

The `MonsterConfig` struct allows detailed configuration of monsters:

```go
monsterConfigs := []services.MonsterConfig{
    {
        Name:        "Goblin",    // Display name of the monster
        Key:         "monster_goblin", // API key for the monster (for XP calculation)
        CR:          0.25,        // Challenge Rating
        Count:       3,           // Number of this monster to place
        RandomPlace: true,        // Whether to place randomly or at a specific position
        Position:    nil,         // Position is nil for random placement
    },
    {
        Name:        "Dragon",
        Key:         "monster_adult_red_dragon",
        CR:          17.0,
        Count:       1,
        RandomPlace: false,       // Fixed position placement
        Position:    &entities.Position{X: 7, Y: 5}, // Specific position
    },
}
```

#### Complete PlayerConfig Options

The `PlayerConfig` struct configures player characters:

```go
playerConfigs := []services.PlayerConfig{
    {
        Name:        "Aragorn",   // Player character name
        Level:       5,           // Character level (used for encounter balancing)
        RandomPlace: true,        // Whether to place randomly
        Position:    nil,         // Position is nil for random placement
    },
    {
        Name:        "Gandalf",
        Level:       10,
        RandomPlace: false,
        Position:    &entities.Position{X: 2, Y: 2}, // Specific position
    },
}
```

#### Complete ItemConfig Options

The `ItemConfig` struct configures items and treasures:

```go
itemConfigs := []services.ItemConfig{
    {
        Name:        "Healing Potion", // Item name
        Key:         "item_healing_potion", // API key for the item
        Value:       50,           // Gold piece value
        Count:       2,            // Number of this item to place
        RandomPlace: true,         // Whether to place randomly
        Position:    nil,          // Position is nil for random placement
    },
    {
        Name:        "Magic Sword",
        Key:         "item_magic_sword",
        Value:       500,
        Count:       1,
        RandomPlace: false,
        Position:    &entities.Position{X: 7, Y: 7}, // Specific position
    },
}
```

### Extending the Library

#### Creating Custom Repositories

You can implement custom repositories by implementing the repository interfaces:

```go
// Custom monster repository
type CustomMonsterRepository struct {
    // Your custom fields here
}

// Implement the MonsterRepository interface
func (r *CustomMonsterRepository) GetMonsterXP(key string) (int, error) {
    // Your custom implementation
    return 100, nil
}

// Use your custom repository with the room service
roomService := services.NewRoomService(WithMonsterRepository(&CustomMonsterRepository{}))
```

#### Custom Encounter Balancing

You can implement custom encounter balancing by implementing the `EncounterBalancer` interface:

```go
// Custom encounter balancer
type CustomBalancer struct {
    // Your custom fields here
}

// Implement the EncounterBalancer interface
func (b *CustomBalancer) AdjustMonsterSelection(configs []services.MonsterConfig, party entities.Party, difficulty entities.EncounterDifficulty) ([]services.MonsterConfig, error) {
    // Your custom implementation
    return configs, nil
}

func (b *CustomBalancer) DetermineEncounterDifficulty(monsters []entities.Monster, party entities.Party) (entities.EncounterDifficulty, error) {
    // Your custom implementation
    return entities.EncounterDifficultyMedium, nil
}

func (b *CustomBalancer) CalculateTargetCR(party entities.Party, difficulty entities.EncounterDifficulty) (float64, error) {
    // Your custom implementation
    return 5.0, nil
}

// Use your custom balancer with the room service
roomService := services.NewRoomService(WithEncounterBalancer(&CustomBalancer{}))
```

### Advanced Room Service Configuration

The `NewRoomService` function accepts functional options for advanced configuration:

```go
// Create a room service with custom options
roomService, err := services.NewRoomService(
    // Custom monster repository
    services.WithMonsterRepository(customMonsterRepo),
    
    // Custom item repository
    services.WithItemRepository(customItemRepo),
    
    // Custom encounter balancer
    services.WithEncounterBalancer(customBalancer),
    
    // Custom random source for deterministic testing
    services.WithRandomSource(rand.New(rand.NewSource(42))),
)
```

### Integration with External Systems

#### Using with a Web Framework

```go
// Example using the library with a web framework (e.g., Gin)
func handleGenerateRoom(c *gin.Context) {
    // Parse request
    var req struct {
        Width      int    `json:"width"`
        Height     int    `json:"height"`
        LightLevel string `json:"lightLevel"`
        UseGrid    bool   `json:"useGrid"`
    }
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Create room config
    roomConfig := services.RoomConfig{
        Width:      req.Width,
        Height:     req.Height,
        LightLevel: entities.LightLevel(req.LightLevel),
        UseGrid:    req.UseGrid,
    }
    
    // Generate room
    roomService, _ := services.NewRoomService()
    room, err := roomService.GenerateRoom(roomConfig)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // Return room data
    c.JSON(200, room)
}
```

#### Using with a Game Engine

```go
// Example using the library with a game engine
func createGameRoom(engine GameEngine, width, height int, lightLevel string) {
    // Create room config
    roomConfig := services.RoomConfig{
        Width:      width,
        Height:     height,
        LightLevel: entities.LightLevel(lightLevel),
        UseGrid:    true,
    }
    
    // Generate room with monsters
    roomService, _ := services.NewRoomService()
    monsterConfigs := []services.MonsterConfig{
        {Name: "Goblin", CR: 0.25, Count: 3, RandomPlace: true},
        {Name: "Orc", CR: 0.5, Count: 2, RandomPlace: true},
    }
    
    room, _ := roomService.PopulateRoomWithMonsters(roomConfig, monsterConfigs)
    
    // Convert to game engine entities
    gameRoom := engine.CreateRoom(room.Width, room.Height)
    for _, monster := range room.Monsters {
        gameMonster := engine.CreateMonster(monster.Name)
        engine.PlaceEntity(gameMonster, monster.Position.X, monster.Position.Y)
    }
    
    // Set lighting based on room's light level
    switch room.LightLevel {
    case entities.LightLevelBright:
        engine.SetLighting(1.0)
    case entities.LightLevelDim:
        engine.SetLighting(0.5)
    case entities.LightLevelDark:
        engine.SetLighting(0.1)
    }
}
```

### Performance Optimization

For large rooms or applications with many entities, consider these optimization strategies:

1. **Use gridless rooms** when spatial positioning is not critical
2. **Batch entity additions** rather than adding entities one by one
3. **Implement custom repositories** with caching for frequently accessed data
4. **Use efficient data structures** for custom implementations

```go
// Example of batch entity addition for better performance
func addManyMonstersEfficiently(room *entities.Room, monsterCount int) {
    // Create all configs at once
    configs := make([]services.MonsterConfig, monsterCount)
    for i := 0; i < monsterCount; i++ {
        configs[i] = services.MonsterConfig{
            Name:        fmt.Sprintf("Monster%d", i),
            CR:          0.25,
            Count:       1,
            RandomPlace: true,
        }
    }
    
    // Add all monsters in a single call
    roomService, _ := services.NewRoomService()
    roomService.AddMonstersToRoom(room, configs)
}
```

## Design Decisions

- The library follows a clean layered architecture with clear separation of concerns
- Repository interfaces allow for easy mocking and testing
- The service layer depends on interfaces, not concrete implementations
- External API dependencies are isolated in the repository layer

## Future Improvements

The following improvements are planned for future releases:

1. **API Data Validation** - Update mock repositories in tests to use real API data to validate API contracts and ensure compatibility with the actual D&D 5e API.

2. **Enhanced Error Handling** - Implement more specific error types for different failure scenarios to help consumers of the library better handle errors.

3. **Additional Room Types** - Support for specialized room types like traps, puzzles, and environmental hazards.

4. **Advanced Encounter Balancing** - More sophisticated algorithms for balancing encounters based on party composition and monster synergies.

The current implementation is an MVP focused on room generation and monster placement. Future enhancements planned for the library include:

### Features
- Traps and hazard system
- Advanced room types (circular, irregular, multi-room layouts)
- Enhanced customization parameters
- Support for connected areas/complexes
- Additional environmental effects
- Item integration and placement
- Advanced monster selection algorithms

### Technical Improvements
- Caching system for frequently used monsters
- Enhanced error handling with user notifications
- Fallback mechanisms for API failures
- Performance optimizations for large room complexes

### Architecture Evolution
- Configuration system for library-wide settings
- Event system for room generation lifecycle hooks
- Extensibility points for custom generation algorithms

These enhancements will be prioritized based on user feedback and needs after the MVP has been tested in real applications.

### Using the Encounter Balancer

The library includes a powerful encounter balancing system that helps create appropriately challenging encounters based on party composition and desired difficulty level.

#### Direct Balancer Usage

You can use the balancer directly for fine-grained control:

```go
import (
    "github.com/fadedpez/dnd5e-roomgen/internal/entities"
    "github.com/fadedpez/dnd5e-roomgen/internal/services"
    "github.com/fadedpez/dnd5e-roomgen/internal/repositories"
)

// Create a monster repository
monsterRepo := repositories.NewAPIMonsterRepository()

// Create a balancer
balancer := services.NewBalancer(monsterRepo)

// Create a party
party := entities.Party{
    Members: []entities.PartyMember{
        {Name: "Aragorn", Level: 5},
        {Name: "Legolas", Level: 5},
        {Name: "Gimli", Level: 5},
        {Name: "Gandalf", Level: 7},
    },
}

// Calculate target Challenge Rating for a medium difficulty encounter
targetCR, err := balancer.CalculateTargetCR(party, entities.EncounterDifficultyMedium)
if err != nil {
    // Handle error
}
fmt.Printf("Target CR for medium encounter: %.2f\n", targetCR)

// Determine the difficulty of an existing encounter
monsters := []entities.Monster{
    {Name: "Goblin", CR: 0.25},
    {Name: "Goblin", CR: 0.25},
    {Name: "Orc Chief", CR: 2.0},
}
difficulty, err := balancer.DetermineEncounterDifficulty(monsters, party)
if err != nil {
    // Handle error
}
fmt.Printf("This encounter is %s for the current party\n", difficulty)
```

#### Integrated Room Service Balancing

For most use cases, you can use the RoomService's integrated balancing methods:

```go
import (
    "github.com/fadedpez/dnd5e-roomgen/internal/entities"
    "github.com/fadedpez/dnd5e-roomgen/internal/services"
)

// Create a room service (automatically initializes the balancer)
roomService, err := services.NewRoomService()
if err != nil {
    // Handle error
}

// Create a party
party := entities.Party{
    Members: []entities.PartyMember{
        {Name: "Aragorn", Level: 5},
        {Name: "Legolas", Level: 5},
        {Name: "Gimli", Level: 5},
        {Name: "Gandalf", Level: 7},
    },
}

// Define monster configurations
monsterConfigs := []services.MonsterConfig{
    {
        Name:        "Goblin",
        Key:         "monster_goblin",
        CR:          0.25,
        Count:       4,
        RandomPlace: true,
    },
    {
        Name:        "Bugbear",
        Key:         "monster_bugbear",
        CR:          1.0,
        Count:       1,
        RandomPlace: true,
    },
}

// Method 1: Balance monster configurations without creating a room
balancedConfigs, err := roomService.BalanceMonsterConfigs(monsterConfigs, party, entities.EncounterDifficultyHard)
if err != nil {
    // Handle error
}
fmt.Printf("Adjusted monster counts for hard difficulty: %+v\n", balancedConfigs)

// Method 2: Generate a room with automatically balanced monsters in one step
room, err := roomService.PopulateRoomWithBalancedMonsters(roomConfig, monsterConfigs, party, entities.EncounterDifficultyHard)
if err != nil {
    // Handle error
}

// Method 3: Analyze the difficulty of an existing room
existingRoom := &entities.Room{
    // ... room properties ...
    Monsters: []entities.Monster{
        {Name: "Dragon", CR: 10},
        {Name: "Kobold", CR: 0.25, Count: 8},
    },
}
difficulty, err := roomService.DetermineRoomDifficulty(existingRoom, party)
if err != nil {
    // Handle error
}
fmt.Printf("This room's encounter is %s for the current party\n", difficulty)
```

#### Balancer Use Cases

1. **Creating Balanced Encounters**:
   - Calculate appropriate challenge ratings for your party
   - Automatically adjust monster counts to match desired difficulty

2. **Analyzing Encounter Difficulty**:
   - Determine if an existing encounter is Easy, Medium, Hard, or Deadly
   - Validate encounter designs against party composition

3. **Dynamic Encounter Scaling**:
   - Scale encounters up or down based on party size and level
   - Maintain appropriate challenge as party composition changes

4. **Difficulty Customization**:
   - Choose from standard D&D 5e difficulty levels (Easy, Medium, Hard, Deadly)
   - Apply consistent difficulty calculations across your application

5. **One-Step Room Generation with Balanced Encounters**:
   - Generate complete rooms with monsters balanced for your party
   - Streamline encounter creation while maintaining appropriate challenge

## License

None
