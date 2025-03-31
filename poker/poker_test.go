package poker

import (
	"online-poker/cards"
	"testing"
)

func TestFiveCardsHighCard(t *testing.T) {
	testedHand := cards.DeckOf(
		cards.CardOf(cards.Clubs, cards.Two),
		cards.CardOf(cards.Spades, cards.King),
		cards.CardOf(cards.Spades, cards.Five),
		cards.CardOf(cards.Spades, cards.Seven),
		cards.CardOf(cards.Diamonds, cards.Jack),
	)
	referenceHand := Hand{
		handType:   HighCard,
		comparison: []cards.Rank{cards.King, cards.Jack, cards.Seven, cards.Five, cards.Two},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestFiveCardsOnePair(t *testing.T) {
	testedHand := cards.DeckOf(
		cards.CardOf(cards.Clubs, cards.Two),
		cards.CardOf(cards.Spades, cards.Two),
		cards.CardOf(cards.Clubs, cards.Ace),
		cards.CardOf(cards.Diamonds, cards.Three),
		cards.CardOf(cards.Hearts, cards.Jack),
	)
	referenceHand := Hand{
		handType:   OnePair,
		comparison: []cards.Rank{cards.Two, cards.Ace},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestThreeCardsOnePairWithLameKicker(t *testing.T) {
	testedHand := cards.DeckOf(
		cards.CardOf(cards.Clubs, cards.Two),
		cards.CardOf(cards.Spades, cards.Ace),
		cards.CardOf(cards.Clubs, cards.Ace),
		cards.CardOf(cards.Diamonds, cards.Three),
		cards.CardOf(cards.Hearts, cards.Five),
	)
	referenceHand := Hand{
		handType: OnePair,
		comparison: []cards.Rank{
			cards.Ace,
			cards.Five,
			cards.Three,
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

func TestFiveCardsStraightWithoutAce(t *testing.T) {
	testedHand := cards.DeckOf(
		cards.CardOf(cards.Clubs, cards.Three),
		cards.CardOf(cards.Spades, cards.Four),
		cards.CardOf(cards.Clubs, cards.Five),
		cards.CardOf(cards.Hearts, cards.Six),
		cards.CardOf(cards.Clubs, cards.Seven),
	)
	referenceHand := Hand{
		handType: Straight,
		comparison: []cards.Rank{
			cards.Seven,
		},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)

}

func testPokerHandRecognition(t *testing.T, testedHand cards.Deck, referenceHand Hand) {
	recognisedHand, _ := hand(testedHand)
	if recognisedHand.handType != referenceHand.handType {
		t.Errorf("hand type should be %s (was %s)", referenceHand.handType.String(), recognisedHand.handType.String())
	}
	for i, rank := range referenceHand.comparison {
		if rank != recognisedHand.comparison[i] {
			t.Errorf("hand comparison[%d] should be %s (was %s)", i, rank.String(), referenceHand.comparison[i].String())
		}
	}
}
