package usecases

import (
	"github.com/tyryoxan/API-RiftRadar/domain"
)

type LeaderboardService struct {
	Repo       domain.LeaderboardRepository
	PlayerRepo domain.PlayerRepository
}

func (ls *LeaderboardService) GetLeaderboardByRegion(region string) (domain.Leaderboard, error) {
	return ls.Repo.GetLeaderboardByRegion(region)
}

func (ls *LeaderboardService) GetLeaderboard() []domain.Leaderboard {
	return ls.Repo.FindAll()
}

func (ls *LeaderboardService) CreateLeaderboard(name string) domain.Leaderboard {
	leaderboard := domain.Leaderboard{
		Name: name,
	}

	ls.Repo.AddLeaderboard(leaderboard)

	return leaderboard
}

func (ls *LeaderboardService) AddPlayer(name string, tag string) error {
	// Get the player by name and tag
	player, err := ls.PlayerRepo.GetPlayerByNameAndTag(name, tag)
	if err != nil {
		return err
	}

	// Get the leaderboard for the player's region
	leaderboard, err := ls.Repo.GetLeaderboardByRegion(player.Region)
	if err != nil {
		// Create a new leaderboard if it doesn't exist
		leaderboard = domain.Leaderboard{
			Name:    player.Region,
			Players: []domain.Player{},
		}
	}

	// Add the player to the leaderboard
	leaderboard.Players = append(leaderboard.Players, player)

	// Save the updated leaderboard
	ls.Repo.AddLeaderboard(leaderboard)

	return nil
}
