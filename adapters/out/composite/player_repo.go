package composite

import (
	"github.com/tyryoxan/API-RiftRadar/adapters/out/memory"
	"github.com/tyryoxan/API-RiftRadar/adapters/out/riot"
	"github.com/tyryoxan/API-RiftRadar/domain"
)

// PlayerRepo is a composite repository that uses both memory and Riot API
type PlayerRepo struct {
	memoryRepo *memory.PlayerRepo
	riotRepo   *riot.RiotPlayerRepo
}

// NewPlayerRepo creates a new composite player repository
func NewPlayerRepo(memoryRepo *memory.PlayerRepo, riotRepo *riot.RiotPlayerRepo) *PlayerRepo {
	return &PlayerRepo{
		memoryRepo: memoryRepo,
		riotRepo:   riotRepo,
	}
}

// Save saves a player to the memory repository
func (c *PlayerRepo) Save(player domain.Player) error {
	return c.memoryRepo.Save(player)
}

// FindAll returns all players from the memory repository
func (c *PlayerRepo) FindAll() []domain.Player {
	return c.memoryRepo.FindAll()
}

// GetPlayerByNameAndRegion gets a player by name and region, first checking memory then Riot API
func (c *PlayerRepo) GetPlayerByNameAndRegion(name, region string) (domain.Player, error) {
	// First try to get from memory
	player, err := c.memoryRepo.GetPlayerByNameAndRegion(name, region)
	if err == nil {
		// Player found in memory
		return player, nil
	}

	// Not found in memory, try Riot API
	player, err = c.riotRepo.GetPlayerByNameAndRegion(name, region)
	if err != nil {
		return domain.Player{}, err
	}

	// Save to memory for future requests
	_ = c.memoryRepo.Save(player)

	return player, nil
}

// GetPlayerByNameAndTag gets a player by name and tag, first checking memory then Riot API
func (c *PlayerRepo) GetPlayerByNameAndTag(name, tag string) (domain.Player, error) {
	// First try to get from memory
	player, err := c.memoryRepo.GetPlayerByNameAndTag(name, tag)
	if err == nil {
		// Player found in memory
		return player, nil
	}

	// Not found in memory, try Riot API
	player, err = c.riotRepo.GetPlayerByNameAndTag(name, tag)
	if err != nil {
		return domain.Player{}, err
	}

	// Save to memory for future requests
	_ = c.memoryRepo.Save(player)

	return player, nil
}
