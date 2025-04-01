package poker

type Player struct {
	name  string
	money int64
}

func NewPlayer(name string, money int64) Player {
	return Player{name: name, money: money}
}
