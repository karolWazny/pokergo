package poker

import (
	"fmt"
	"online-poker/cards"
	"testing"
)

func TestOneCard(t *testing.T) {
	testedHand := cards.DeckOf(
		cards.CardOf(cards.Clubs, cards.Ace),
	)
	recognisedHand := hand(testedHand)
	if recognisedHand.handType != HighCard {
		t.Errorf("hand type should be HighCard")
	}
	if recognisedHand.comparison[0] != cards.Ace {
		fmt.Println(recognisedHand.comparison[0])
		t.Errorf("hand comparison should be Ace")
	}
}

func TestTwoCardsHighCard(t *testing.T) {
	testedHand := cards.DeckOf(
		cards.CardOf(cards.Clubs, cards.Two), cards.CardOf(cards.Spades, cards.King),
	)
	recognisedHand := hand(testedHand)
	if recognisedHand.handType != HighCard {
		t.Errorf("hand type should be HighCard")
	}
	if recognisedHand.comparison[0] != cards.King {
		fmt.Println(recognisedHand.comparison[0])
		t.Errorf("hand comparison[0] should be King")
	}
	if recognisedHand.comparison[1] != cards.Two {
		t.Errorf("hand comparison[1] should be Two")
	}
}

func TestThreeCardsOnePair(t *testing.T) {
	testedHand := cards.DeckOf(
		cards.CardOf(cards.Clubs, cards.Two), cards.CardOf(cards.Spades, cards.Two), cards.CardOf(cards.Clubs, cards.Ace),
	)
	recognisedHand := hand(testedHand)
	if recognisedHand.handType != OnePair {
		t.Errorf("hand type should be OnePair")
	}
	if recognisedHand.comparison[0] != cards.Two {
		t.Errorf("hand comparison[0] should be Two")
	}
	if recognisedHand.comparison[1] != cards.Ace {
		t.Errorf("hand comparison[1] should be Ace")
	}
}

func TestThreeCardsOnePairWithLameKicker(t *testing.T) {
	testedHand := cards.DeckOf(
		cards.CardOf(cards.Clubs, cards.Two), cards.CardOf(cards.Spades, cards.Ace), cards.CardOf(cards.Clubs, cards.Ace),
	)
	recognisedHand := hand(testedHand)
	if recognisedHand.handType != OnePair {
		t.Errorf("hand type should be OnePair")
	}
	if recognisedHand.comparison[0] != cards.Ace {
		t.Errorf("hand comparison[0] should be Ace")
	}
	if recognisedHand.comparison[1] != cards.Two {
		t.Errorf("hand comparison[1] should be Two")
	}
}

func TestFiveCardsTwoPair(t *testing.T) {
	testedHand := cards.DeckOf(
		cards.CardOf(cards.Clubs, cards.Two),
		cards.CardOf(cards.Spades, cards.Two),
		cards.CardOf(cards.Clubs, cards.Ace),
		cards.CardOf(cards.Hearts, cards.Ace),
		cards.CardOf(cards.Clubs, cards.Jack),
	)
	recognisedHand := hand(testedHand)
	if recognisedHand.handType != TwoPair {
		t.Errorf("hand type should be TwoPair")
	}
	if recognisedHand.comparison[0] != cards.Ace {
		t.Errorf("hand comparison[0] should be Ace")
	}
	if recognisedHand.comparison[1] != cards.Two {
		t.Errorf("hand comparison[1] should be Two")
	}
	if recognisedHand.comparison[2] != cards.Jack {
		t.Errorf("hand comparison[2] should be Jack")
	}
}
