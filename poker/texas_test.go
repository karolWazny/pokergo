package poker

import (
	"slices"
	"testing"
)

func TestThreePlayersCanStartAGame(t *testing.T) {
	table := prepareThreePlayerTable()
	game := table.StartGame()
	visibleGameState := game.GetVisibleGameState()
	if len(visibleGameState.Players) != 3 {
		t.Errorf("There should be 3 players")
	}
}

func TestGameIsFinishedWhenAllButOneFold(t *testing.T) {
	table := prepareThreePlayerTable()
	game := table.StartGame()
	game.Fold()
	game.Fold()
	visibleGameState := game.GetVisibleGameState()
	if visibleGameState.Round != FINISHED {
		t.Errorf("When only one player remains, the game should be finished (is %s)", visibleGameState.Round)
	}
}

func TestPlayerCannotCallIfThereWasNoRaise(t *testing.T) {
	table := prepareThreePlayerTable()
	game := table.StartGame()
	game.Call()
	game.Call()
	availableActions := game.AvailableActions()
	if slices.Contains(availableActions, call) {
		t.Errorf("Player cannot call if there was no raise")
	}
}

func TestPlayerCannotCheckIfThereWasRaise(t *testing.T) {
	table := prepareThreePlayerTable()
	game := table.StartGame()
	game.Call()
	game.Call()
	game.Check()
	// flop
	game.Raise(50)
	availableActions := game.AvailableActions()
	if slices.Contains(availableActions, check) {
		t.Errorf("Player cannot check if there was raise")
	}
}

func TestRoundIsNotFinishedIfThereWasRaise(t *testing.T) {
	table := prepareThreePlayerTable()
	game := table.StartGame()
	game.Call()
	game.Call()
	game.Raise(50)
	round := game.GetVisibleGameState().Round
	if round != PREFLOP {
		t.Errorf("Round should be PREFLOP (is %s)", round)
	}
}

func TestSecondRaiseCausesAReRaise(t *testing.T) {
	table := prepareThreePlayerTable()
	game := table.StartGame()
	game.Call()
	game.Call()
	game.Check()
	game.Raise(50)
	player := game.CurrentPlayer()
	currentMoney := player.player.money
	game.Raise(50)
	moneyAfterRaise := player.player.money
	difference := currentMoney - moneyAfterRaise
	if difference != 100 {
		t.Errorf("Raising 50 after raise of 50 should cause re-raise (100$ total) (was %d)", difference)
	}
}

func prepareThreePlayerTable() Table {
	table := NewTable(20, 50)
	master := NewPlayer("MasterOfDisaster", 1500)
	table.AddPlayer(&master)
	badman := NewPlayer("BadMannTM", 1500)
	table.AddPlayer(&badman)
	hanku := NewPlayer("hank.prostokat", 1500)
	table.AddPlayer(&hanku)
	return table
}
