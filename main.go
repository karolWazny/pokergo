package main

import (
	"fmt"
	"online-poker/cards"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	deck := cards.CreateDeck()
	deck = deck.Shuffled()
	fmt.Println(deck)
	hand, deck := deck.Deal(5)
	fmt.Println(hand)
	fmt.Println(deck)
}
