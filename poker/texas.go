package poker

import (
	"errors"
	"online-poker/cards"
	"slices"
	"strconv"
)

type TexasHoldEmAction string

const (
	check TexasHoldEmAction = "check"
	fold                    = "fold"
	call                    = "call"
	raise                   = "raise"
)

type TexasHoldEmRound int8

const (
	PREFLOP TexasHoldEmRound = iota
	FLOP
	TURN
	RIVER
)

type Table struct {
	players     []*Player
	smallBlind  int64
	bigBlind    int64
	dealerIndex int
}

func NewTable(smallBlind int64, bigBlind int64) Table {
	return Table{
		players:     make([]*Player, 0),
		smallBlind:  smallBlind,
		bigBlind:    bigBlind,
		dealerIndex: -1,
	}
}

func (table *Table) AddPlayer(player *Player) {
	table.players = append(table.players, player)
}

func (table *Table) StartGame() Game {
	table.dealerIndex = (table.dealerIndex + 1) % len(table.players)
	orderedPlayers := append(table.players[table.dealerIndex+1:], table.players[:table.dealerIndex+1]...)
	texasPlayers := make([]*TexasPlayer, len(orderedPlayers))
	deck := cards.CreateDeck().Shuffled()
	for i, player := range orderedPlayers {
		hand, smallerDeck := deck.Deal(2)
		deck = smallerDeck
		texasPlayers[i] = &TexasPlayer{
			player:     player,
			hand:       hand,
			hasFolded:  false,
			currentPot: 0,
		}
	}
	texasPlayers[0].currentPot = table.smallBlind
	texasPlayers[0].player.money -= table.smallBlind
	texasPlayers[1].currentPot = table.bigBlind
	texasPlayers[1].player.money -= table.bigBlind
	return Game{
		players:           texasPlayers,
		lastBet:           table.bigBlind,
		deck:              deck,
		activePlayerIndex: 2,
		community:         make([]cards.Card, 0),
		round:             PREFLOP,
		isFinished:        false,
	}
}

type Game struct {
	players           []*TexasPlayer
	lastBet           int64
	deck              cards.Deck
	activePlayerIndex int
	community         []cards.Card
	round             TexasHoldEmRound
	isFinished        bool
}

func (game *Game) CurrentPlayer() *TexasPlayer {
	return game.players[game.activePlayerIndex]
}

func (game *Game) AvailableActions() []TexasHoldEmAction {
	if game.isFinished {
		return []TexasHoldEmAction{}
	}
	currentPlayer := game.players[game.activePlayerIndex]
	previousPlayerPot := game.getPreviousPlayerPot()
	availableActions := []TexasHoldEmAction{fold, raise}
	if previousPlayerPot == currentPlayer.currentPot {
		availableActions = append(availableActions, check)
	} else {
		availableActions = append(availableActions, call)
	}
	return availableActions
}

func (game *Game) Call() error {
	availableActions := game.AvailableActions()
	if !slices.Contains(availableActions, call) {
		return errors.New("call action not available")
	}
	currentPlayer := game.players[game.activePlayerIndex]
	pot := game.getPreviousPlayerPot()
	difference := pot - currentPlayer.currentPot
	currentPlayer.currentPot = pot
	currentPlayer.player.money -= difference
	game.nextPlayer()
	return nil
}

func (game *Game) Fold() error {
	availableActions := game.AvailableActions()
	if !slices.Contains(availableActions, fold) {
		return errors.New("fold action not available")
	}
	game.players[game.activePlayerIndex].hasFolded = true
	game.nextPlayer()
	return nil
}

func (game *Game) Check() error {
	availableActions := game.AvailableActions()
	if !slices.Contains(availableActions, check) {
		return errors.New("check action not available")
	}
	game.nextPlayer()
	return nil
}

func (game *Game) CommunityCards() []cards.Card {
	return game.community
}

func (game *Game) playersInGame() int {
	playersInGame := 0
	for _, player := range game.players {
		if !player.hasFolded {
			playersInGame++
		}
	}
	return playersInGame
}

func (game *Game) getPreviousPlayerPot() int64 {
	for i := 1; i < len(game.players); i++ {
		previousPlayerIndex := (game.activePlayerIndex - i + len(game.players)) % len(game.players)
		if !game.players[previousPlayerIndex].hasFolded {
			return game.players[previousPlayerIndex].currentPot
		}
	}
	panic("There should be at least two active players!")
}

func (game *Game) nextPlayer() {
	game.players[game.activePlayerIndex].hasPlayed = true
	game.activePlayerIndex = (game.activePlayerIndex + 1) % len(game.players)
	for game.players[game.activePlayerIndex].hasFolded {
		game.activePlayerIndex = (game.activePlayerIndex + 1) % len(game.players)
	}
	isRoundFinished := game.isCurrentRoundFinished()
	isGameFinished := (isRoundFinished && game.round == RIVER) || game.playersInGame() == 1
	if isGameFinished {
		game.isFinished = true
		return
	} else if isRoundFinished {
		game.finishRound()
	}
}

func (game *Game) finishRound() {
	if game.round == RIVER {
		// trigger showdown
	}
	game.activePlayerIndex = 0
	for _, player := range game.players {
		player.hasPlayed = false
	}
	_, game.deck = game.deck.Deal(1)
	cardsToShow := 1
	isFlop := game.round == PREFLOP
	if isFlop {
		cardsToShow = 3
	}
	newCards, deck := game.deck.Deal(cardsToShow)
	game.deck = deck
	game.community = append(game.CommunityCards(), newCards.Cards...)
	game.round++
}

func (game *Game) isCurrentRoundFinished() bool {
	for _, player := range game.players {
		if !player.hasFolded && !player.hasPlayed {
			return false
		}
	}
	return true
}

type TexasPlayer struct {
	player     *Player
	hand       cards.Deck
	hasFolded  bool
	hasPlayed  bool
	currentPot int64
}

func (texasPlayer TexasPlayer) String() string {
	return texasPlayer.player.String() + " " + texasPlayer.hand.String() + " " + strconv.FormatInt(texasPlayer.currentPot, 10)
}
