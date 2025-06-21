# DnD 5e Room Generator

A Go library for generating dynamic rooms for Dungeons & Dragons 5th Edition adventures.

## Features

- Generate room layouts with configurable dimensions and light levels
- Place monsters with random or specific positioning
- Place players with random or specific positioning
- Automatic grid initialization for spatial tracking
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

```go
import (
    "github.com/fadedpez/dnd5e-roomgen/internal/entities"
    "github.com/fadedpez/dnd5e-roomgen/internal/services"
)

// Create a room service
roomService := services.NewRoomService()

// Configure a room
roomConfig := services.RoomConfig{
    Width:       15,
    Height:      10,
    LightLevel:  entities.LightLevelDim,
    Description: "A dimly lit dungeon chamber",
    UseGrid:     true,
}

// Generate the room
room, err := roomService.GenerateRoom(roomConfig)
if err != nil {
    // Handle error
}
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

### Generating a Room with Both Players and Monsters

```go
// Generate a room and add both players and monsters in one step
room, err := roomService.PopulateRoomWithMonstersAndPlayers(roomConfig, monsterConfigs, playerConfigs)
if err != nil {
    // Handle error
}

// Access room properties
fmt.Printf("Room: %dx%d, %s\n", room.Width, room.Height, room.Description)

// Access player properties
for i, player := range room.Players {
    fmt.Printf("Player %d: %s (Level %d) at position (%d,%d)\n",
        i+1, player.Name, player.Level, player.Position.X, player.Position.Y)
}

// Access monster properties
for i, monster := range room.Monsters {
    fmt.Printf("Monster %d: %s (CR %.2f) at position (%d,%d)\n",
        i+1, monster.Name, monster.CR, monster.Position.X, monster.Position.Y)
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

### Using the Monster Repository

```go
import (
    "github.com/fadedpez/dnd5e-roomgen/internal/repositories"
    "net/http"
)

// Create a monster repository that uses the DnD 5e API
httpClient := &http.Client{Timeout: 10 * time.Second}
monsterRepo := repositories.NewAPIMonsterRepository(httpClient)

// Get monster XP
xp, err := monsterRepo.GetMonsterXP("monster_goblin")
if err != nil {
    // Handle error
}
fmt.Printf("Goblin XP: %d\n", xp)
```

### Testing with the Repository Layer

```go
// Create a test monster repository for unit testing
testRepo := &repositories.TestMonsterRepository{
    xpValues: map[string]int{
        "monster_goblin": 50,
        "monster_orc": 100,
    },
}

// Use the test repository in your service
roomService := services.NewRoomService(testRepo)
```

## Design Decisions

- The library follows a clean layered architecture with clear separation of concerns
- Repository interfaces allow for easy mocking and testing
- The service layer depends on interfaces, not concrete implementations
- External API dependencies are isolated in the repository layer

## Future Enhancements

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

## License

None
