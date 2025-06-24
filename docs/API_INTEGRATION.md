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
monsterConfig := services.ConvertAPIMonsterToConfig(apiMonster)

// Convert a slice of monsters
monsterConfigs := services.ConvertAPIMonsterSliceToConfigs(apiMonsters)
```

### Item Converters

```go
// Convert a single item
itemConfig := services.ConvertAPIItemToConfig(apiItem)

// Convert a slice of items
itemConfigs := services.ConvertAPIItemSliceToConfigs(apiItems)
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
monsterConfigs := services.ConvertAPIMonsterSliceToConfigs(apiMonsters)

// Create separate configs for each monster instance you want
var finalConfigs []services.MonsterConfig
for _, config := range monsterConfigs {
    // Add two of each monster type
    for i := 0; i < 2; i++ {
        // Create a copy of the config
        newConfig := config
        finalConfigs = append(finalConfigs, newConfig)
    }
}

// Create a room service
roomService, err := services.NewRoomService()
if err != nil {
    // Handle error
}

// Create and populate a room with the monsters
room, err := roomService.GenerateRoom(roomConfig)
if err != nil {
    // Handle error
}

err = roomService.AddMonstersToRoom(room, finalConfigs)
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
        CR:          0.25,
        RandomPlace: true,
    },
    {
        Key:         "goblin",
        Name:        "Goblin",
        CR:          0.25,
        RandomPlace: true,
    },
    {
        Key:         "orc",
        Name:        "Orc",
        CR:          0.5,
        RandomPlace: false,
        Position:    &entities.Position{X: 5, Y: 5},
    },
}

// Create a room service
roomService, err := services.NewRoomService()
if err != nil {
    // Handle error
}

// Create and populate a room with the monsters
room, err := roomService.GenerateRoom(roomConfig)
if err != nil {
    // Handle error
}

err = roomService.AddMonstersToRoom(room, monsterConfigs)
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
apiClient := dnd5eapi.NewClient(httpClient)
apiMonsters, err := apiClient.GetMonsters([]string{"goblin"})
if err != nil {
    // Handle error
}

// Convert API monsters to configs
monsterConfigs := services.ConvertAPIMonsterSliceToConfigs(apiMonsters)

// Create multiple instances of each monster
var finalConfigs []services.MonsterConfig
for _, config := range monsterConfigs {
    // Add two of each monster type
    for i := 0; i < 2; i++ {
        // Create a copy of the config
        newConfig := config
        finalConfigs = append(finalConfigs, newConfig)
    }
}

// Add a custom monster to the configs
customMonster := services.MonsterConfig{
    Key:         "custom_boss",
    Name:        "Dungeon Master's Special",
    CR:          5.0,
    RandomPlace: false,
    Position:    &entities.Position{X: 7, Y: 7},
}
finalConfigs = append(finalConfigs, customMonster)

// Create a room with both API-sourced and custom monsters
roomService, err := services.NewRoomService()
if err != nil {
    // Handle error
}

room, err := roomService.GenerateRoom(roomConfig)
if err != nil {
    // Handle error
}

err = roomService.AddMonstersToRoom(room, finalConfigs)
if err != nil {
    // Handle error
}
```

**Benefits:**
- Flexibility to use official content alongside custom creations
- Gradual migration path from hardcoded to API-sourced entities
- Fallback options when API entities aren't available

## Using the Encounter Balancer with the API

The library includes a powerful encounter balancing system that helps create appropriately challenging encounters based on party composition and desired difficulty level.

### Direct Balancer Usage with API

You can use the balancer directly with API monsters for fine-grained control:

```go
import (
    "fmt"
    "net/http"
    "time"
    
    dnd5eapi "github.com/fadedpez/dnd5e-api/client"
    "github.com/fadedpez/dnd5e-roomgen/internal/entities"
    "github.com/fadedpez/dnd5e-roomgen/internal/services"
)

// Create an HTTP client with timeout
httpClient := &http.Client{Timeout: 10 * time.Second}

// Create a DnD 5e API client
apiClient := dnd5eapi.NewClient(httpClient)

// Create a balancer
balancer := services.NewBalancer()

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

// Fetch some monsters from the API
apiMonsters, err := apiClient.GetMonsters([]string{"goblin", "orc", "bugbear"})
if err != nil {
    // Handle error
}

// Convert API monsters to entities
var monsters []entities.Monster
for _, apiMonster := range apiMonsters {
    monster := entities.Monster{
        Name: apiMonster.Name,
        Key:  apiMonster.Index,
        CR:   apiMonster.ChallengeRating,
    }
    monsters = append(monsters, monster)
}

// Determine the difficulty of an encounter with these monsters
difficulty, err := balancer.DetermineEncounterDifficulty(monsters, party)
if err != nil {
    // Handle error
}
fmt.Printf("This encounter is %s for the current party\n", difficulty)

// Adjust monster selection to match target difficulty
adjustedMonsters, err := balancer.AdjustMonsterSelection(monsters, party, entities.EncounterDifficultyMedium)
if err != nil {
    // Handle error
}
fmt.Printf("Adjusted monster count: %d\n", len(adjustedMonsters))
```

### Integrated Room Service Balancing with API Monsters

For most use cases, you can use the RoomService's integrated balancing methods with API-sourced monsters:

```go
import (
    "fmt"
    "net/http"
    "time"
    
    dnd5eapi "github.com/fadedpez/dnd5e-api/client"
    "github.com/fadedpez/dnd5e-roomgen/internal/entities"
    "github.com/fadedpez/dnd5e-roomgen/internal/services"
)

// Create HTTP client with timeout
httpClient := &http.Client{Timeout: 10 * time.Second}

// Create DnD 5e API client
apiClient := dnd5eapi.NewClient(httpClient)

// Fetch monsters from the API
apiMonsters, err := apiClient.GetMonsters([]string{"goblin", "bugbear"})
if err != nil {
    fmt.Printf("Error fetching monsters: %v\n", err)
    return
}

// Convert API monsters to configs
monsterConfigs := services.ConvertAPIMonsterSliceToConfigs(apiMonsters)

// Create multiple instances of each monster
var finalConfigs []services.MonsterConfig
for _, config := range monsterConfigs {
    // Add two of each monster type
    for i := 0; i < 2; i++ {
        // Create a copy of the config
        newConfig := config
        finalConfigs = append(finalConfigs, newConfig)
    }
}

// Create a room service
roomService, err := services.NewRoomService()
if err != nil {
    // Handle error
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

// Configure room
roomConfig := services.RoomConfig{
    Width:       20,
    Height:      15,
    LightLevel:  entities.LightLevelDim,
    Description: "A musty dungeon chamber with cobwebs in the corners",
    UseGrid:     true,
}

// Generate a room with automatically balanced monsters in one step
room, err := roomService.GenerateAndPopulateRoom(
    roomConfig,
    finalConfigs,
    nil, // No player configs
    nil, // No item configs
    &party,
    entities.EncounterDifficultyHard,
)
if err != nil {
    // Handle error
}

// Analyze the difficulty of the room
difficulty, err := roomService.DetermineRoomDifficulty(room, party)
if err != nil {
    // Handle error
}
fmt.Printf("This room's encounter is %s for the current party\n", difficulty)
```

### Balancer Use Cases with the API

1. **Creating Balanced Encounters with API Monsters**:
   - Fetch monsters from the DnD 5e API
   - Convert to room service configs
   - Calculate appropriate challenge ratings for your party
   - Automatically adjust monster counts to match desired difficulty

2. **Analyzing API-Sourced Encounter Difficulty**:
   - Determine if an existing encounter is Easy, Medium, Hard, or Deadly
   - Validate encounter designs against party composition

3. **Dynamic Encounter Scaling with API Monsters**:
   - Scale encounters up or down based on party size and level
   - Maintain appropriate challenge as party composition changes

## Player Configuration with API Integration

When integrating with the API, you'll often want to add players to your rooms alongside API-sourced monsters and items. Here's how to configure players:

```go
// Configure players
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

// Add players to a room with API-sourced monsters
room, err := roomService.PopulateRoomWithMonstersAndPlayers(roomConfig, monsterConfigs, playerConfigs)
if err != nil {
    // Handle error
}

// Access player properties
for i, player := range room.Players {
    fmt.Printf("Player %d: %s (Level %d) at position (%d,%d)\n",
        i+1, player.Name, player.Level, player.Position.X, player.Position.Y)
}
```

## Item Configuration with API Integration

You can combine API-sourced items with custom items:

```go
// Fetch items from the API
apiItems, err := apiClient.GetEquipment([]string{"potion-of-healing", "longsword"})
if err != nil {
    fmt.Printf("Error fetching items: %v\n", err)
    return
}

// Convert API items to configs
apiItemConfigs := services.ConvertAPIItemSliceToConfigs(apiItems)

// Add custom items
customItemConfigs := []services.ItemConfig{
    {
        Name:        "Magic Amulet",
        Key:         "custom_magic_amulet",
        Value:       500,          // Gold piece value
        RandomPlace: false,
        Position:    &entities.Position{X: 7, Y: 7}, // Specific position
    },
}

// Combine API and custom items
allItemConfigs := append(apiItemConfigs, customItemConfigs...)

// Add items to a room with API-sourced monsters
room, err := roomService.PopulateRoomWithMonstersAndItems(roomConfig, monsterConfigs, allItemConfigs)
if err != nil {
    // Handle error
}

// Access item properties
for i, item := range room.Items {
    fmt.Printf("Item %d: %s at position (%d,%d)\n",
        i+1, item.Name, item.Position.X, item.Position.Y)
}
```

## Generating a Treasure Room with API Integration

You can create treasure rooms with optional guardian monsters sourced from the API:

```go
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
    monsterConfigs := services.ConvertAPIMonsterSliceToConfigs(apiMonsters)
    itemConfigs := services.ConvertAPIItemSliceToConfigs(apiItems)
    
    // Create multiple instances of each monster
    var finalMonsterConfigs []services.MonsterConfig
    for _, config := range monsterConfigs {
        // Add two of each monster type
        for i := 0; i < 2; i++ {
            // Create a copy of the config
            newConfig := config
            finalMonsterConfigs = append(finalMonsterConfigs, newConfig)
        }
    }
    
    // Create a room service
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
    
    // Create a party
    party := entities.Party{
        Members: []entities.PartyMember{
            {Name: "Aragorn", Level: 5},
            {Name: "Legolas", Level: 4},
            {Name: "Gimli", Level: 5},
            {Name: "Gandalf", Level: 7},
        },
    }
    
    // Generate a room with API-sourced monsters and items
    room, err := roomService.GenerateAndPopulateRoom(
        roomConfig,
        finalMonsterConfigs,
        nil, // No player configs
        itemConfigs,
        &party,
        entities.EncounterDifficultyHard,
    )
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

## Best Practices

1. **Error Handling**: Always check for errors when fetching from the API and when converting entities.
2. **Caching**: Consider caching API responses to improve performance for frequently used monsters and items.
3. **Fallbacks**: Implement fallbacks to direct configuration when API calls fail.
4. **Validation**: Validate converted configs before using them, especially for custom properties not derived from the API.
5. **Type Conversion**: Be aware of type conversions (like float32 to float64 for Challenge Rating) when debugging.
6. **API Key Format**: Use the hyphenated keys from the API (e.g., "goblin", "adult-red-dragon", "potion-of-healing").
7. **Monster Selection**: Choose a variety of monsters across different CR values to give the balancer flexibility in adjusting counts.

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

## Important Note on Monster Selection

**The library does not automatically search the DnD 5e API for monsters based on Challenge Rating (CR).** 

The current implementation requires you to explicitly specify which monsters you want to use by providing their keys (e.g., "goblin", "orc", "adult-red-dragon"). The balancer will then:

1. Calculate the appropriate target CR based on your party and desired difficulty
2. Adjust the **counts** of your pre-selected monsters to achieve that target CR

This design gives you full control over which monsters appear in your encounters while still providing automatic balancing.

### Example workflow:

```go
// 1. YOU select which monsters to use (the library doesn't do this automatically)
monsterKeys := []string{"goblin", "orc", "bugbear"}

// 2. Fetch those specific monsters from the API
apiMonsters, err := apiClient.GetMonsters(monsterKeys)
if err != nil {
    // Handle error
}

// 3. Convert to configs with initial counts
monsterConfigs := services.ConvertAPIMonsterSliceToConfigs(apiMonsters)

// 4. Create multiple instances of each monster
var finalConfigs []services.MonsterConfig
for _, config := range monsterConfigs {
    // Add two of each monster type
    for i := 0; i < 2; i++ {
        // Create a copy of the config
        newConfig := config
        finalConfigs = append(finalConfigs, newConfig)
    }
}

// 5. Let the balancer adjust the counts to match your desired difficulty
balancedConfigs, err := roomService.BalanceMonsterConfigs(finalConfigs, party, entities.EncounterDifficultyHard)
if err != nil {
    // Handle error
}

// 6. The balancer has adjusted the counts, but not changed which monsters are used
fmt.Println("Balanced monster counts:")
for _, config := range balancedConfigs {
    fmt.Printf("- %s: %d (CR %.2f)\n", config.Name, config.Count, config.CR)
}
```

### Future Enhancement

A future enhancement to the library may include automatic monster selection from the API based on CR ranges, but currently, you must explicitly choose which monsters to include in your encounters.

## Using the Encounter Balancer with the API

The library includes a powerful encounter balancing system that helps create appropriately challenging encounters based on party composition and desired difficulty level.

### Direct Balancer Usage with API

You can use the balancer directly with API monsters for fine-grained control:

```go
import (
    "fmt"
    "net/http"
    "time"
    
    dnd5eapi "github.com/fadedpez/dnd5e-api/client"
    "github.com/fadedpez/dnd5e-roomgen/internal/entities"
    "github.com/fadedpez/dnd5e-roomgen/internal/services"
)

// Create an HTTP client with timeout
httpClient := &http.Client{Timeout: 10 * time.Second}

// Create a DnD 5e API client
apiClient := dnd5eapi.NewClient(httpClient)

// Create a balancer
balancer := services.NewBalancer()

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

// Fetch some monsters from the API
apiMonsters, err := apiClient.GetMonsters([]string{"goblin", "orc", "bugbear"})
if err != nil {
    // Handle error
}

// Convert API monsters to entities
var monsters []entities.Monster
for _, apiMonster := range apiMonsters {
    monster := entities.Monster{
        Name: apiMonster.Name,
        Key:  apiMonster.Index,
        CR:   apiMonster.ChallengeRating,
    }
    monsters = append(monsters, monster)
}

// Determine the difficulty of an encounter with these monsters
difficulty, err := balancer.DetermineEncounterDifficulty(monsters, party)
if err != nil {
    // Handle error
}
fmt.Printf("This encounter is %s for the current party\n", difficulty)

// Adjust monster selection to match target difficulty
adjustedMonsters, err := balancer.AdjustMonsterSelection(monsters, party, entities.EncounterDifficultyMedium)
if err != nil {
    // Handle error
}
fmt.Printf("Adjusted monster count: %d\n", len(adjustedMonsters))
```

### Integrated Room Service Balancing with API Monsters

For most use cases, you can use the RoomService's integrated balancing methods with API-sourced monsters:

```go
import (
    "fmt"
    "net/http"
    "time"
    
    dnd5eapi "github.com/fadedpez/dnd5e-api/client"
    "github.com/fadedpez/dnd5e-roomgen/internal/entities"
    "github.com/fadedpez/dnd5e-roomgen/internal/services"
)

// Create HTTP client with timeout
httpClient := &http.Client{Timeout: 10 * time.Second}

// Create DnD 5e API client
apiClient := dnd5eapi.NewClient(httpClient)

// Fetch monsters from the API
apiMonsters, err := apiClient.GetMonsters([]string{"goblin", "bugbear"})
if err != nil {
    fmt.Printf("Error fetching monsters: %v\n", err)
    return
}

// Convert API monsters to configs
monsterConfigs := services.ConvertAPIMonsterSliceToConfigs(apiMonsters)

// Create multiple instances of each monster
var finalConfigs []services.MonsterConfig
for _, config := range monsterConfigs {
    // Add two of each monster type
    for i := 0; i < 2; i++ {
        // Create a copy of the config
        newConfig := config
        finalConfigs = append(finalConfigs, newConfig)
    }
}

// Create a room service
roomService, err := services.NewRoomService()
if err != nil {
    // Handle error
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

// Configure room
roomConfig := services.RoomConfig{
    Width:       20,
    Height:      15,
    LightLevel:  entities.LightLevelDim,
    Description: "A musty dungeon chamber with cobwebs in the corners",
    UseGrid:     true,
}

// Generate a room with automatically balanced monsters in one step
room, err := roomService.GenerateAndPopulateRoom(
    roomConfig,
    finalConfigs,
    nil, // No player configs
    nil, // No item configs
    &party,
    entities.EncounterDifficultyHard,
)
if err != nil {
    // Handle error
}

// Analyze the difficulty of the room
difficulty, err := roomService.DetermineRoomDifficulty(room, party)
if err != nil {
    // Handle error
}
fmt.Printf("This room's encounter is %s for the current party\n", difficulty)
```

### Balancer Use Cases with the API

1. **Creating Balanced Encounters with API Monsters**:
   - Fetch monsters from the DnD 5e API
   - Convert to room service configs
   - Calculate appropriate challenge ratings for your party
   - Automatically adjust monster counts to match desired difficulty

2. **Analyzing API-Sourced Encounter Difficulty**:
   - Determine if an existing encounter is Easy, Medium, Hard, or Deadly
   - Validate encounter designs against party composition

3. **Dynamic Encounter Scaling with API Monsters**:
   - Scale encounters up or down based on party size and level
   - Maintain appropriate challenge as party composition changes
