package usecases

import "github.com/tyryoxan/API-RiftRadar/domain"

type PlayerService struct {
	Repo domain.PlayerRepository
}

func (ps *PlayerService) GetPlayers() []domain.Player {
	return ps.Repo.FindAll()
}

func (ps *PlayerService) GetPlayerByNameAndRegion(name, region string) (domain.Player, error) {
	return ps.Repo.GetPlayerByNameAndRegion(name, region)
}

func (ps *PlayerService) GetPlayerByNameAndTag(name, tag string) (domain.Player, error) {
	return ps.Repo.GetPlayerByNameAndTag(name, tag)
}
