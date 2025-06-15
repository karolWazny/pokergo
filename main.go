package main

import (
	"github.com/karolWazny/poker-go/poker"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons._.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	table := poker.NewTable(20, 50)
	master := poker.NewPlayer("MasterOfDisaster", 1500)
	table.AddPlayer(&master)
	badman := poker.NewPlayer("BadMannTM", 1500)
	table.AddPlayer(&badman)
	hanku := poker.NewPlayer("hank.prostokat", 1500)
	table.AddPlayer(&hanku)
	game := table.StartGame()
	// preflop
	game.Fold()
	game.Call()
	game.Check()
	// flop
	game.Check()
	game.Check()
	// turn
	game.Raise(100)
	game.Raise(100)
	game.Call()
	// river
	game.Raise(150)
	game.Call()
	// showdown
	game.GetVisibleGameState().Print()
}
