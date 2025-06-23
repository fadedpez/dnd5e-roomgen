package examples

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fadedpez/dnd5e-api/clients/dnd5e"
	"github.com/fadedpez/dnd5e-roomgen/internal/services"
)

// This file contains examples of how to integrate the DnD 5e API with the room generation library.
// These examples are for documentation purposes and are not meant to be executed directly.

// ExampleAddMonsterFromAPI demonstrates how to retrieve a monster from the DnD 5e API
// and add it to a room using the room generation service.
func ExampleAddMonsterFromAPI() {
	// Create the HTTP client
	httpClient := &http.Client{}

	// Create the API client config
	config := &dnd5e.DND5eAPIConfig{
		Client: httpClient,
	}

	// Initialize the API client
	apiClient, err := dnd5e.NewDND5eAPI(config)
	if err != nil {
		log.Fatalf("Failed to create DND5e API client: %v", err)
	}

	// Retrieve a monster from the API
	monster, err := apiClient.GetMonster("goblin")
	if err != nil {
		log.Fatalf("Failed to get monster: %v", err)
	}

	// Convert the API monster to a MonsterConfig with a count of 5
	monsterConfig := services.ConvertAPIMonsterToConfig(monster, 5)

	// Create a room service
	roomService, err := services.NewRoomService()
	if err != nil {
		log.Fatalf("Failed to create room service: %v", err)
	}

	// Generate a room using the correct method
	room, err := roomService.GenerateRoom(services.RoomConfig{
		Description: "A small cave where goblins have set up camp",
		Width:       30,
		Height:      30,
	})
	if err != nil {
		log.Fatalf("Failed to generate room: %v", err)
	}

	// Add the monster to the room using the appropriate method
	err = roomService.AddMonstersToRoom(room, []services.MonsterConfig{*monsterConfig})
	if err != nil {
		log.Fatalf("Failed to add monster to room: %v", err)
	}

	fmt.Printf("Added %d %s to room %s\n", monsterConfig.Count, monsterConfig.Name, room.Description)
}
