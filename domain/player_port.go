package domain

type PlayerRepository interface {
	Save(player Player) error
	FindAll() []Player
	GetPlayerByNameAndRegion(name, region string) (Player, error)
	GetPlayerByNameAndTag(name, tag string) (Player, error)
}
