package poker

import (
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
		t.Errorf("When only one player remains, the game is finished (is %s)", visibleGameState.Round)
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
