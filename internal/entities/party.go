package entities

// PartyMember represents a player character in a party
type PartyMember struct {
	Name  string
	Level int
}

// Party represents a group of player characters
type Party struct {
	Members []PartyMember
}

// AverageLevel calculates the average level of the party
func (p *Party) AverageLevel() float64 {
	if len(p.Members) == 0 {
		return 0
	}

	totalLevel := 0
	for _, member := range p.Members {
		totalLevel += member.Level
	}

	return float64(totalLevel) / float64(len(p.Members))
}

// Size returns the number of members in the party
func (p *Party) Size() int {
	return len(p.Members)
}

// PlayersToParty creates a Party from a slice of Player entities
func PlayersToParty(players []Player) Party {
	members := make([]PartyMember, len(players))
	for i, player := range players {
		members[i] = PartyMember{
			Name:  player.Name,
			Level: player.Level,
		}
	}
	return Party{Members: members}
}
