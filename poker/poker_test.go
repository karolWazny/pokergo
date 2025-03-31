package poker

import (
	"online-poker/cards"
	"testing"
)

func TestOneCard(t *testing.T) {
	testedHand := cards.DeckOf(
		cards.CardOf(cards.Clubs, cards.Ace),
	)
	referenceHand := Hand{
		handType:   HighCard,
		comparison: []cards.Rank{cards.Ace},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestTwoCardsHighCard(t *testing.T) {
	testedHand := cards.DeckOf(
		cards.CardOf(cards.Clubs, cards.Two), cards.CardOf(cards.Spades, cards.King),
	)
	referenceHand := Hand{
		handType:   HighCard,
		comparison: []cards.Rank{cards.King, cards.Two},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestThreeCardsOnePair(t *testing.T) {
	testedHand := cards.DeckOf(
		cards.CardOf(cards.Clubs, cards.Two),
		cards.CardOf(cards.Spades, cards.Two),
		cards.CardOf(cards.Clubs, cards.Ace),
	)
	referenceHand := Hand{
		handType:   OnePair,
		comparison: []cards.Rank{cards.Two, cards.Ace},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestThreeCardsOnePairWithLameKicker(t *testing.T) {
	testedHand := cards.DeckOf(
		cards.CardOf(cards.Clubs, cards.Two), cards.CardOf(cards.Spades, cards.Ace), cards.CardOf(cards.Clubs, cards.Ace),
	)
	referenceHand := Hand{
		handType: OnePair,
		comparison: []cards.Rank{
			cards.Ace,
			cards.Two,
		},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestFiveCardsTwoPair(t *testing.T) {
	testedHand := cards.DeckOf(
		cards.CardOf(cards.Clubs, cards.Two),
		cards.CardOf(cards.Spades, cards.Two),
		cards.CardOf(cards.Clubs, cards.Ace),
		cards.CardOf(cards.Hearts, cards.Ace),
		cards.CardOf(cards.Clubs, cards.Jack),
	)
	referenceHand := Hand{
		handType: TwoPair,
		comparison: []cards.Rank{
			cards.Ace,
			cards.Two,
			cards.Jack,
		},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestFiveCardsThreeOfAKind(t *testing.T) {
	testedHand := cards.DeckOf(
		cards.CardOf(cards.Clubs, cards.Two),
		cards.CardOf(cards.Spades, cards.Two),
		cards.CardOf(cards.Clubs, cards.Two),
		cards.CardOf(cards.Hearts, cards.Jack),
		cards.CardOf(cards.Clubs, cards.Ace),
	)
	referenceHand := Hand{
		handType: ThreeOfAKind,
		comparison: []cards.Rank{
			cards.Two,
			cards.Ace,
			cards.Jack,
		},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)

}

func testPokerHandRecognition(t *testing.T, testedHand cards.Deck, referenceHand Hand) {
	recognisedHand := hand(testedHand)
	if recognisedHand.handType != referenceHand.handType {
		t.Errorf("hand type should be %s (was %s)", referenceHand.handType.String(), recognisedHand.handType.String())
	}
	for i, rank := range recognisedHand.comparison {
		if rank != referenceHand.comparison[i] {
			t.Errorf("hand comparison[%d] should be %s (was %s)", i, rank.String(), referenceHand.comparison[i].String())
		}
	}
}
