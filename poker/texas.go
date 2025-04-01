package poker

import (
	"errors"
	"fmt"
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

func (round TexasHoldEmRound) String() string {
	switch round {
	case PREFLOP:
		return "preflop"
	case FLOP:
		return "flop"
	case TURN:
		return "turn"
	case RIVER:
		return "river"
	default:
		return "unknown"
	}
}

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
	game.incrementActivePlayerIndex()
	game.changeActivePlayerToFirstNonFolded()
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
	game.changeActivePlayerToFirstNonFolded()
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

func (game *Game) changeActivePlayerToFirstNonFolded() {
	for game.players[game.activePlayerIndex].hasFolded {
		game.incrementActivePlayerIndex()
	}
}

func (game *Game) incrementActivePlayerIndex() {
	game.activePlayerIndex = (game.activePlayerIndex + 1) % len(game.players)
}

func (game *Game) isCurrentRoundFinished() bool {
	for _, player := range game.players {
		if !player.hasFolded && !player.hasPlayed {
			return false
		}
	}
	return true
}

func (game *Game) GetVisibleGameState() VisibleGameState {
	players := make([]TexasPlayerPublicInfo, len(game.players))
	for i, player := range game.players {
		players[i] = player.GetPublicInfo()
	}
	activePlayer := game.players[game.activePlayerIndex].GetPublicInfo()
	round := game.round
	dealer := game.players[len(game.players)-1].GetPublicInfo()
	community := game.CommunityCards()
	return VisibleGameState{
		Players:      players,
		Round:        round,
		Dealer:       dealer,
		ActivePlayer: activePlayer,
		Community:    community,
	}
}

type VisibleGameState struct {
	Players      []TexasPlayerPublicInfo
	Round        TexasHoldEmRound
	ActivePlayer TexasPlayerPublicInfo
	Dealer       TexasPlayerPublicInfo
	Community    []cards.Card
}

func (gameState VisibleGameState) Print() {
	fmt.Printf("Little Friendly Game of Poker, stage: %s\n", gameState.Round)
	fmt.Printf("Dealer: %s\n", gameState.Dealer.Name)
	fmt.Printf("Community Cards:\n")
	for _, card := range gameState.Community {
		fmt.Printf("- %s\n", card)
	}
	fmt.Printf("Players:\n")
	for _, player := range gameState.Players {
		fmt.Printf("- %s", player)
	}
	fmt.Printf("Now playing: %s\n", gameState.ActivePlayer.Name)
}

type TexasPlayerPublicInfo struct {
	Name       string
	Money      int64
	HasFolded  bool
	CurrentPot int64
	Cards      []cards.Card
}

func (playerPublicInfo TexasPlayerPublicInfo) String() string {
	foldedString := "in game"
	if playerPublicInfo.HasFolded {
		foldedString = "has folded"
	}
	return fmt.Sprintf("%s, pot: %d$, %s, total: %d$, Cards: %s\n",
		playerPublicInfo.Name,
		playerPublicInfo.CurrentPot,
		foldedString,
		playerPublicInfo.Money,
		playerPublicInfo.Cards)
}

type TexasPlayer struct {
	player     *Player
	hand       cards.Deck
	hasFolded  bool
	hasPlayed  bool
	currentPot int64
}

func (texasPlayer TexasPlayer) GetPublicInfo() TexasPlayerPublicInfo {
	return TexasPlayerPublicInfo{
		Name:       texasPlayer.player.name,
		Money:      texasPlayer.player.money,
		HasFolded:  texasPlayer.hasFolded,
		CurrentPot: texasPlayer.currentPot,
		Cards:      []cards.Card{},
	}
}

func (texasPlayer TexasPlayer) String() string {
	return texasPlayer.player.String() + " " + texasPlayer.hand.String() + " " + strconv.FormatInt(texasPlayer.currentPot, 10)
}
