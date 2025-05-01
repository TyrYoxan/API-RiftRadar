package memory

import (
	"fmt"
	"github.com/tyryoxan/API-RiftRadar/domain"
	"strings"
	"sync"
)

type PlayerRepo struct {
	players []domain.Player
	mu      sync.RWMutex
}

func NewPlayerRepo() *PlayerRepo {
	return &PlayerRepo{
		players: make([]domain.Player, 0),
	}
}

func (pr *PlayerRepo) Save(player domain.Player) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	// Check if player already exists, update if it does
	for i, p := range pr.players {
		if strings.EqualFold(p.Name, player.Name) && strings.EqualFold(p.Region, player.Region) {
			pr.players[i] = player
			return nil
		}
	}

	// Add new player if not found
	pr.players = append(pr.players, player)
	return nil
}

func (pr *PlayerRepo) FindAll() []domain.Player {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	// Return a copy to avoid race conditions
	result := make([]domain.Player, len(pr.players))
	copy(result, pr.players)
	return result
}

func (pr *PlayerRepo) GetPlayerByNameAndRegion(name, region string) (domain.Player, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	// Case-insensitive search
	for _, player := range pr.players {
		if strings.EqualFold(player.Name, name) && strings.EqualFold(player.Region, region) {
			return player, nil
		}
	}

	return domain.Player{}, fmt.Errorf("player not found with name %s and region %s", name, region)
}

func (pr *PlayerRepo) GetPlayerByNameAndTag(name, tag string) (domain.Player, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	// Case-insensitive search for name
	// For tag, we don't have a direct mapping in our model, so we'll just check the name
	for _, player := range pr.players {
		if strings.EqualFold(player.Name, name) {
			return player, nil
		}
	}

	return domain.Player{}, fmt.Errorf("player not found with name %s and tag %s", name, tag)
}
