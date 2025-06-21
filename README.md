# DnD 5e Room Generator

A Go library for generating dynamic rooms for Dungeons & Dragons 5th Edition adventures.

## Features

- Generate room layouts with configurable dimensions and light levels
- Place monsters with random or specific positioning
- Automatic grid initialization for spatial tracking
- Flexible service layer for easy integration with applications

## Installation

```bash
go get github.com/fadedpez/dnd5e-roomgen
```

## Architecture

The library follows a clean layered architecture:

- **Entities Layer**: Core domain objects like Room, Monster, Position, etc.
- **Services Layer**: Business logic for room generation and monster placement
- **Repository Layer**: (Future) Data access for monsters, items, etc.

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

## Design Decisions

- The library focuses on spatial layout and monster placement, not combat mechanics
- Monster stats (AC, HP, etc.) are intentionally not included and should be handled by consumer applications
- The service layer provides a clean API that hides implementation details

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
- Repository layer for data access (monsters, items, etc.)
- Configuration system for library-wide settings
- Event system for room generation lifecycle hooks
- Extensibility points for custom generation algorithms

These enhancements will be prioritized based on user feedback and needs after the MVP has been tested in real applications.

## License

None
