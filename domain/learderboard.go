package domain

type Leaderboard struct {
	Name    string   `json:"name"`
	Players []Player `json:"players"`
}

type LeaderboardRepository interface {
	FindAll() []Leaderboard
	GetLeaderboardByRegion(region string) (Leaderboard, error)
	AddLeaderboard(leaderboard Leaderboard) error
}
