package poker

import "online-poker/cards"

type Table struct {
	players     []Player
	smallBlind  int64
	bigBlind    int64
	dealerIndex int
}

func NewTable(smallBlind int64, bigBlind int64) Table {
	return Table{
		players:     make([]Player, 0),
		smallBlind:  smallBlind,
		bigBlind:    bigBlind,
		dealerIndex: -1,
	}
}

func (table *Table) AddPlayer(player Player) {
	table.players = append(table.players, player)
}

func (table *Table) StartGame() Game {
	table.dealerIndex = (table.dealerIndex + 1) % len(table.players)
	orderedPlayers := append(table.players[table.dealerIndex+1:], table.players[:table.dealerIndex+1]...)
	texasPlayers := make([]TexasPlayer, len(orderedPlayers))
	deck := cards.CreateDeck().Shuffled()
	for i, player := range orderedPlayers {
		hand, smallerDeck := deck.Deal(2)
		deck = smallerDeck
		texasPlayers[i] = TexasPlayer{
			player:     &player,
			hand:       hand,
			hasFolded:  false,
			currentPot: 0,
		}
	}
	texasPlayers[0].currentPot = table.smallBlind
	texasPlayers[1].currentPot = table.bigBlind
	return Game{
		players:               texasPlayers,
		lastBet:               table.bigBlind,
		deck:                  deck,
		actionsInCurrentRound: 0,
		activePlayerIndex:     3,
	}
}

type Game struct {
	players               []TexasPlayer
	lastBet               int64
	deck                  cards.Deck
	actionsInCurrentRound int
	activePlayerIndex     int
}

type TexasPlayer struct {
	player     *Player
	hand       cards.Deck
	hasFolded  bool
	currentPot int64
}
