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
	FINISHED
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
	case FINISHED:
		return "finished"
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
	}
}

type Game struct {
	players           []*TexasPlayer
	winner            *TexasPlayer
	lastBet           int64
	deck              cards.Deck
	activePlayerIndex int
	community         []cards.Card
	round             TexasHoldEmRound
}

func (game *Game) Winner() (*TexasPlayer, error) {
	if game.round != FINISHED {
		return nil, errors.New("there is no winner before game end")
	}
	return game.winner, nil
}

func (game *Game) CurrentPlayer() (*TexasPlayer, error) {
	if game.round != FINISHED {
		return game.unsafeGetCurrentPlayer(), nil
	}
	return nil, errors.New("in finished game there is no current player")
}

func (game *Game) AvailableActions() []TexasHoldEmAction {
	if game.round == FINISHED {
		return []TexasHoldEmAction{}
	}
	currentPlayer := game.unsafeGetCurrentPlayer()
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
	currentPlayer := game.unsafeGetCurrentPlayer()
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
	game.unsafeGetCurrentPlayer().hasFolded = true
	if game.playersInGame() == 1 {
		// finish game
		lastActivePlayerIndex := slices.IndexFunc(game.players, func(player *TexasPlayer) bool {
			return !player.hasFolded
		})
		game.winner = game.players[lastActivePlayerIndex]
		game.activePlayerIndex = -1
		game.round = FINISHED
		game.transferPotToWinner()
	} else {
		game.nextPlayer()
	}
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

func (game *Game) Raise(amount int64) error {
	if amount < game.lastBet {
		return errors.New("amount must be greater than last bet")
	}
	game.lastBet = amount
	realAmount := game.getPreviousPlayerPot() - game.getCurrentPlayerPot() + amount
	currentPlayer, e := game.CurrentPlayer()
	if e == nil {
		currentPlayer.currentPot += realAmount
		currentPlayer.player.money -= realAmount
	}
	game.nextPlayer()
	return nil
}

func (game *Game) CommunityCards() []cards.Card {
	return game.community
}

func (game *Game) GetVisibleGameState() VisibleGameState {
	visibleGameState := VisibleGameState{}
	players := make([]TexasPlayerPublicInfo, len(game.players))
	for i, player := range game.players {
		players[i] = player.GetPublicInfo()
	}
	visibleGameState.Players = players
	activePlayer, e := game.CurrentPlayer()
	if e == nil {
		activePlayerInfo := activePlayer.GetPublicInfo()
		visibleGameState.ActivePlayer = &activePlayerInfo
	}
	visibleGameState.Round = game.round
	visibleGameState.Dealer = game.players[len(game.players)-1].GetPublicInfo()
	visibleGameState.Community = game.CommunityCards()
	return visibleGameState
}

func (game *Game) transferPotToWinner() {
	for _, player := range game.players {
		game.winner.player.money += player.currentPot
	}
}

func (game *Game) unsafeGetCurrentPlayer() *TexasPlayer {
	return game.players[game.activePlayerIndex]
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

func (game *Game) getCurrentPlayerPot() int64 {
	return game.unsafeGetCurrentPlayer().currentPot
}

func (game *Game) nextPlayer() {
	game.unsafeGetCurrentPlayer().hasPlayed = true
	game.changeActivePlayerToFirstNonFolded()
	isRoundFinished := game.isCurrentRoundFinished()
	isGameFinished := (isRoundFinished && game.round == RIVER) || game.playersInGame() == 1
	if isGameFinished {
		game.round = FINISHED
		return
	} else if isRoundFinished {
		game.finishRound()
	}
}

func (game *Game) finishRound() {
	if game.round == RIVER {
		// trigger showdown
	}
	game.activePlayerIndex = len(game.players) - 1
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
	game.incrementActivePlayerIndex()
	for game.unsafeGetCurrentPlayer().hasFolded {
		game.incrementActivePlayerIndex()
	}
}

func (game *Game) incrementActivePlayerIndex() {
	game.activePlayerIndex = (game.activePlayerIndex + 1) % len(game.players)
}

func (game *Game) isCurrentRoundFinished() bool {
	uniquePots := make(map[int64]bool)
	for _, player := range game.players {
		if !player.hasFolded {
			if !player.hasPlayed {
				return false
			}
			uniquePots[player.currentPot] = true
		}
	}
	return len(uniquePots) == 1
}

type VisibleGameState struct {
	Players      []TexasPlayerPublicInfo
	Round        TexasHoldEmRound
	ActivePlayer *TexasPlayerPublicInfo
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
