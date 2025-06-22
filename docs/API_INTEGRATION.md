# DnD 5e API Integration Guide

This guide explains how to integrate the DnD 5e API with the room generation library to create rooms with real DnD 5e monsters and items.

## Overview

The room generation library provides converter functions in the `services` package that make it easy to transform data from the DnD 5e API into the format expected by the room generation service.

## Prerequisites

- The DnD 5e API client: `github.com/fadedpez/dnd5e-api/client`
- The room generation library: `github.com/fadedpez/dnd5e-roomgen`

## Converter Functions

### Monster Converters

```go
// Convert a single monster with a specified count
monsterConfig := services.ConvertAPIMonsterToConfig(apiMonster, count)

// Convert a slice of monsters, each with the same count
monsterConfigs := services.ConvertAPIMonsterSliceToConfigs(apiMonsters, count)
```

### Item Converters

```go
// Convert a single item with a specified count
itemConfig := services.ConvertAPIItemToConfig(apiItem, count)

// Convert a slice of items, each with the same count
itemConfigs := services.ConvertAPIItemSliceToConfigs(apiItems, count)
```

## Integration Approaches

The room generation library offers multiple ways to integrate with the DnD 5e API, giving you flexibility based on your specific needs:

### Approach 1: Direct API Entity Conversion

This approach is ideal when you're already working with the DnD 5e API and want to use those entities directly with the room generator.

```go
// Fetch monsters from the DnD 5e API
apiClient := dnd5eapi.NewClient(httpClient)
apiMonsters, err := apiClient.GetMonsters([]string{"goblin", "orc"})
if err != nil {
    // Handle error
}

// Convert API monsters to room service configs
monsterConfigs := services.ConvertAPIMonsterSliceToConfigs(apiMonsters, 2) // 2 of each monster

// Create a room and add the monsters
roomService := services.NewRoomService()
room, err := roomService.PopulateRoomWithMonsters(roomConfig, monsterConfigs)
if err != nil {
    // Handle error
}
```

**Benefits:**
- Seamless integration with the DnD 5e API
- Access to all official DnD monsters and items
- Automatic mapping of challenge ratings and other properties
- Minimal code required to convert between API and service formats

### Approach 2: Direct Configuration Creation

This approach gives you more control by directly creating configuration objects without using the DnD 5e API.

```go
// Create monster configs directly
monsterConfigs := []services.MonsterConfig{
    {
        Key:         "goblin",
        Name:        "Goblin",
        Count:       3,
        CR:          0.25,
        RandomPlace: true,
    },
    {
        Key:         "orc",
        Name:        "Orc",
        Count:       1,
        CR:          0.5,
        RandomPlace: false,
        Position:    &entities.Position{X: 5, Y: 5},
    },
}

// Create a room and add the monsters
roomService := services.NewRoomService()
room, err := roomService.PopulateRoomWithMonsters(roomConfig, monsterConfigs)
if err != nil {
    // Handle error
}
```

**Benefits:**
- No dependency on external API calls
- Complete control over entity properties
- Ability to create custom monsters not in the official DnD database
- Faster execution without network requests

### Approach 3: Hybrid Approach

You can also mix both approaches, using API entities when available and direct configuration for custom entities.

```go
// Start with some API monsters
apiMonsters, err := apiClient.GetMonsters([]string{"goblin"})
if err != nil {
    // Handle error
}

// Convert API monsters to configs
monsterConfigs := services.ConvertAPIMonsterSliceToConfigs(apiMonsters, 2)

// Add a custom monster to the configs
customMonster := services.MonsterConfig{
    Key:         "custom_boss",
    Name:        "Dungeon Master's Special",
    Count:       1,
    CR:          5.0,
    RandomPlace: false,
    Position:    &entities.Position{X: 7, Y: 7},
}
monsterConfigs = append(monsterConfigs, customMonster)

// Create a room with both API-sourced and custom monsters
room, err := roomService.PopulateRoomWithMonsters(roomConfig, monsterConfigs)
if err != nil {
    // Handle error
}
```

**Benefits:**
- Flexibility to use official content alongside custom creations
- Gradual migration path from hardcoded to API-sourced entities
- Fallback options when API entities aren't available

## Field Mapping Details

When converting from DnD 5e API entities to room service configurations, the following field mappings are applied:

### Monster Mapping

| API Monster Field | MonsterConfig Field | Notes |
|-------------------|---------------------|-------|
| `Key`             | `Key`               | Used for identification |
| `Name`            | `Name`              | Display name of the monster |
| `ChallengeRating` | `CR`                | Converted from float32 to float64 |
| -                 | `Count`             | Set by the count parameter |
| -                 | `RandomPlace`       | Defaults to true unless Position is provided |
| -                 | `Position`          | Optional, can be set after conversion |

### Item Mapping

| API Item Field | ItemConfig Field | Notes |
|----------------|------------------|-------|
| `Key`          | `Key`            | Used for identification |
| `Name`         | `Name`           | Display name of the item |
| -              | `Count`          | Set by the count parameter |
| -              | `RandomPlace`    | Defaults to true |
| -              | `Position`       | Optional, can be set after conversion |

## Complete Integration Example

Here's a complete example showing how to fetch monsters and items from the DnD 5e API and use them to create a room:

```go
package main

import (
    "fmt"
    "net/http"
    "time"

    dnd5eapi "github.com/fadedpez/dnd5e-api/client"
    "github.com/fadedpez/dnd5e-roomgen/internal/entities"
    "github.com/fadedpez/dnd5e-roomgen/internal/services"
)

func main() {
    // Create HTTP client with timeout
    httpClient := &http.Client{Timeout: 10 * time.Second}
    
    // Create DnD 5e API client
    apiClient := dnd5eapi.NewClient(httpClient)
    
    // Fetch monsters from the API
    apiMonsters, err := apiClient.GetMonsters([]string{"goblin", "orc"})
    if err != nil {
        fmt.Printf("Error fetching monsters: %v\n", err)
        return
    }
    
    // Fetch items from the API
    apiItems, err := apiClient.GetEquipment([]string{"potion-of-healing", "longsword"})
    if err != nil {
        fmt.Printf("Error fetching items: %v\n", err)
        return
    }
    
    // Convert API entities to room service configs
    monsterConfigs := services.ConvertAPIMonsterSliceToConfigs(apiMonsters, 2)
    itemConfigs := services.ConvertAPIItemSliceToConfigs(apiItems, 1)
    
    // Create room service
    roomService, err := services.NewRoomService()
    if err != nil {
        fmt.Printf("Error creating room service: %v\n", err)
        return
    }
    
    // Configure room
    roomConfig := services.RoomConfig{
        Width:       20,
        Height:      15,
        LightLevel:  entities.LightLevelDim,
        Description: "A musty dungeon chamber with cobwebs in the corners",
        UseGrid:     true,
    }
    
    // Create room with monsters and items
    room, err := roomService.PopulateRoomWithMonstersAndItems(roomConfig, monsterConfigs, itemConfigs)
    if err != nil {
        fmt.Printf("Error populating room: %v\n", err)
        return
    }
    
    // Print room details
    fmt.Printf("Room created: %dx%d, %s\n", room.Width, room.Height, room.Description)
    fmt.Printf("Monsters in room: %d\n", len(room.Monsters))
    for i, monster := range room.Monsters {
        fmt.Printf("  Monster %d: %s (CR %.2f) at position (%d,%d)\n",
            i+1, monster.Name, monster.CR, monster.Position.X, monster.Position.Y)
    }
    
    fmt.Printf("Items in room: %d\n", len(room.Items))
    for i, item := range room.Items {
        fmt.Printf("  Item %d: %s at position (%d,%d)\n",
            i+1, item.Name, item.Position.X, item.Position.Y)
    }
}
```

## Using the Monster Repository with the API

```go
import (
    "github.com/fadedpez/dnd5e-roomgen/internal/repositories"
    "net/http"
    "time"
)

// Create a monster repository that uses the DnD 5e API
httpClient := &http.Client{Timeout: 10 * time.Second}
monsterRepo := repositories.NewAPIMonsterRepository(httpClient)

// Get monster XP
xp, err := monsterRepo.GetMonsterXP("goblin")
if err != nil {
    // Handle error
}
fmt.Printf("Goblin XP: %d\n", xp)
```

## Best Practices

1. **Error Handling**: Always check for errors when fetching from the API and when converting entities.
2. **Caching**: Consider caching API responses to improve performance for frequently used monsters and items.
3. **Fallbacks**: Implement fallbacks to direct configuration when API calls fail.
4. **Validation**: Validate converted configs before using them, especially for custom properties not derived from the API.
5. **Type Conversion**: Be aware of type conversions (like float32 to float64 for Challenge Rating) when debugging.
6. **API Key Format**: Use the hyphenated keys from the API (e.g., "goblin", "adult-red-dragon", "potion-of-healing").

## API Resources

The DnD 5e API provides access to various resources:

### Monster Resources
- Basic monsters: "goblin", "orc", "kobold"
- Dragons: "adult-blue-dragon", "young-red-dragon"
- Humanoids: "bandit-captain", "cultist"

### Item Resources
- Weapons: "longsword", "shortbow", "dagger"
- Armor: "chain-mail", "leather-armor", "shield"
- Potions: "potion-of-healing", "potion-of-invisibility"
- Misc: "rope", "lantern", "spellbook"

## Further Resources

- [DnD 5e API Documentation](https://www.dnd5eapi.co/docs/)
- [Room Generation Library Documentation](https://github.com/fadedpez/dnd5e-roomgen)
- Complete examples in the `/examples` directory
