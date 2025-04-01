package main

import (
	"fmt"
	"online-poker/poker"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	table := poker.NewTable(20, 50)
	table.AddPlayer(poker.NewPlayer("MasterOfDistaster", 1500))
	table.AddPlayer(poker.NewPlayer("BadMannTM", 1500))
	table.AddPlayer(poker.NewPlayer("hank.prostokat", 1500))
	game := table.StartGame()
	fmt.Println(game)
}
