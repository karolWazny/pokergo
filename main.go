package main

import (
	"fmt"
	"online-poker/poker"
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
	player, _ := game.CurrentPlayer()
	fmt.Println(player)
	fmt.Println(game.CommunityCards())
	game.Call()
	player, _ = game.CurrentPlayer()
	fmt.Println(player)
	fmt.Println(game.CommunityCards())
	game.Call()
	player, _ = game.CurrentPlayer()
	fmt.Println(player)
	fmt.Println(game.CommunityCards())
	game.Call()
	player, _ = game.CurrentPlayer()
	fmt.Println(player)
	// here it should already be flop
	fmt.Println(game.CommunityCards())
}
