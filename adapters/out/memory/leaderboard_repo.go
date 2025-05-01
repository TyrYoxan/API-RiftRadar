package memory

import (
	"fmt"
	"github.com/tyryoxan/API-RiftRadar/domain"
	"strings"
	"sync"
)

type LeaderboardRepo struct {
	leaderboards []domain.Leaderboard
	mu           sync.RWMutex
}

func NewLeaderboardRepo() *LeaderboardRepo {
	return &LeaderboardRepo{
		leaderboards: make([]domain.Leaderboard, 0),
	}
}

func (lr *LeaderboardRepo) FindAll() []domain.Leaderboard {
	lr.mu.RLock()
	defer lr.mu.RUnlock()

	// Return a copy to avoid race conditions
	result := make([]domain.Leaderboard, len(lr.leaderboards))
	copy(result, lr.leaderboards)
	return result
}

func (lr *LeaderboardRepo) GetLeaderboardByRegion(region string) (domain.Leaderboard, error) {
	lr.mu.RLock()
	defer lr.mu.RUnlock()

	// Case-insensitive search
	for _, leaderboard := range lr.leaderboards {
		if strings.EqualFold(leaderboard.Name, region) {
			return leaderboard, nil
		}
	}

	return domain.Leaderboard{}, fmt.Errorf("leaderboard not found for region %s", region)
}

// AddLeaderboard adds a new leaderboard or updates an existing one
func (lr *LeaderboardRepo) AddLeaderboard(leaderboard domain.Leaderboard) error {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	// Check if leaderboard already exists for this region, update if it does
	for i, lb := range lr.leaderboards {
		if strings.EqualFold(lb.Name, leaderboard.Name) {
			lr.leaderboards[i] = leaderboard
			return nil
		}
	}

	// Add new leaderboard if not found
	lr.leaderboards = append(lr.leaderboards, leaderboard)
	return nil
}
