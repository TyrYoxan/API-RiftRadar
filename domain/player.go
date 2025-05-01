package domain

type Player struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	Elo    string `json:"elo"`
	Lp     int    `json:"lp"`
}
