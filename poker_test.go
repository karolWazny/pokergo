package pokergo

import (
	"testing"
)

func TestFiveCardsHighCard(t *testing.T) {
	testedHand := DeckOf(
		CardOf(Clubs, Two),
		CardOf(Spades, King),
		CardOf(Spades, Five),
		CardOf(Spades, Seven),
		CardOf(Diamonds, Jack),
	)
	referenceHand := Hand{
		handType:   HighCard,
		comparison: []Rank{King, Jack, Seven, Five, Two},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestFiveCardsOnePair(t *testing.T) {
	testedHand := DeckOf(
		CardOf(Clubs, Two),
		CardOf(Spades, Two),
		CardOf(Clubs, Ace),
		CardOf(Diamonds, Three),
		CardOf(Hearts, Jack),
	)
	referenceHand := Hand{
		handType:   OnePair,
		comparison: []Rank{Two, Ace},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestThreeCardsOnePairWithLameKicker(t *testing.T) {
	testedHand := DeckOf(
		CardOf(Clubs, Two),
		CardOf(Spades, Ace),
		CardOf(Clubs, Ace),
		CardOf(Diamonds, Three),
		CardOf(Hearts, Five),
	)
	referenceHand := Hand{
		handType: OnePair,
		comparison: []Rank{
			Ace,
			Five,
			Three,
			Two,
		},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestFiveCardsTwoPair(t *testing.T) {
	testedHand := DeckOf(
		CardOf(Clubs, Two),
		CardOf(Spades, Two),
		CardOf(Clubs, Ace),
		CardOf(Hearts, Ace),
		CardOf(Clubs, Jack),
	)
	referenceHand := Hand{
		handType: TwoPair,
		comparison: []Rank{
			Ace,
			Two,
			Jack,
		},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestFiveCardsThreeOfAKind(t *testing.T) {
	testedHand := DeckOf(
		CardOf(Clubs, Two),
		CardOf(Spades, Two),
		CardOf(Clubs, Two),
		CardOf(Hearts, Jack),
		CardOf(Clubs, Ace),
	)
	referenceHand := Hand{
		handType: ThreeOfAKind,
		comparison: []Rank{
			Two,
			Ace,
			Jack,
		},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)

}

func TestFiveCardsThreeOfAKind2(t *testing.T) {
	testedHand := DeckOf(
		CardOf(Spades, Six),
		CardOf(Diamonds, Six),
		CardOf(Clubs, Six),
		CardOf(Clubs, King),
		CardOf(Diamonds, Ten),
	)
	referenceHand := Hand{
		handType: ThreeOfAKind,
		comparison: []Rank{
			Six,
			King,
			Ten,
		},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)

}

func TestFiveCardsStraightWithoutAce(t *testing.T) {
	testedHand := DeckOf(
		CardOf(Clubs, Three),
		CardOf(Spades, Four),
		CardOf(Clubs, Five),
		CardOf(Hearts, Six),
		CardOf(Clubs, Seven),
	)
	referenceHand := Hand{
		handType: Straight,
		comparison: []Rank{
			Seven,
		},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)

}

func TestStraightWithLowAce(t *testing.T) {
	testedHand := DeckOf(
		CardOf(Clubs, Ace),
		CardOf(Spades, Two),
		CardOf(Clubs, Three),
		CardOf(Hearts, Four),
		CardOf(Clubs, Five),
	)
	referenceHand := Hand{
		handType: Straight,
		comparison: []Rank{
			Five,
		},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)

}

func TestFlush(t *testing.T) {
	testedHand := DeckOf(
		CardOf(Spades, Jack),
		CardOf(Spades, Two),
		CardOf(Spades, Three),
		CardOf(Spades, Four),
		CardOf(Spades, Five),
	)
	referenceHand := Hand{
		handType: Flush,
		comparison: []Rank{
			Jack,
			Five,
			Four,
			Three,
			Two,
		},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestFullHouse(t *testing.T) {

	testedHand := DeckOf(
		CardOf(Spades, Two),
		CardOf(Hearts, Two),
		CardOf(Diamonds, Two),
		CardOf(Clubs, Five),
		CardOf(Spades, Five),
	)
	referenceHand := Hand{
		handType: FullHouse,
		comparison: []Rank{
			Two,
			Five,
		},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestFourOfAKind(t *testing.T) {

	testedHand := DeckOf(
		CardOf(Spades, Two),
		CardOf(Hearts, Two),
		CardOf(Diamonds, Two),
		CardOf(Clubs, Two),
		CardOf(Spades, Five),
	)
	referenceHand := Hand{
		handType: FourOfAKind,
		comparison: []Rank{
			Two,
			Five,
		},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestStraightFlush(t *testing.T) {

	testedHand := DeckOf(
		CardOf(Spades, Ace),
		CardOf(Spades, Two),
		CardOf(Spades, Three),
		CardOf(Spades, Four),
		CardOf(Spades, Five),
	)
	referenceHand := Hand{
		handType: StraightFlush,
		comparison: []Rank{
			Five,
		},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestRoyalFlush(t *testing.T) {

	testedHand := DeckOf(
		CardOf(Spades, Ace),
		CardOf(Spades, King),
		CardOf(Spades, Jack),
		CardOf(Spades, Queen),
		CardOf(Spades, Ten),
	)
	referenceHand := Hand{
		handType:   RoyalFlush,
		comparison: []Rank{},
	}
	testPokerHandRecognition(t, testedHand, referenceHand)
}

func TestPairTrumpsHighCard(t *testing.T) {
	highCardHand := DeckOf(
		CardOf(Clubs, Two),
		CardOf(Spades, King),
		CardOf(Spades, Five),
		CardOf(Spades, Seven),
		CardOf(Diamonds, Jack),
	)
	pairOfTwosHand := DeckOf(
		CardOf(Clubs, Two),
		CardOf(Spades, Two),
		CardOf(Spades, Five),
		CardOf(Spades, Seven),
		CardOf(Diamonds, Jack),
	)
	highCard, _ := RecogniseHand(highCardHand)
	pairOfTwos, _ := RecogniseHand(pairOfTwosHand)
	result := CompareHands(highCard, pairOfTwos)
	if result != SecondWins {
		t.Errorf("pair of twos should trump high card")
	}
}

func TestPairOfKingsTrumpsPairOfTwos(t *testing.T) {
	pairOfKingsHand := DeckOf(
		CardOf(Clubs, King),
		CardOf(Spades, King),
		CardOf(Spades, Five),
		CardOf(Spades, Seven),
		CardOf(Diamonds, Jack),
	)
	pairOfTwosHand := DeckOf(
		CardOf(Clubs, Two),
		CardOf(Spades, Two),
		CardOf(Spades, Five),
		CardOf(Spades, Seven),
		CardOf(Diamonds, Jack),
	)
	pairOfKings, _ := RecogniseHand(pairOfKingsHand)
	pairOfTwos, _ := RecogniseHand(pairOfTwosHand)
	result := CompareHands(pairOfKings, pairOfTwos)
	if result != FirstWins {
		t.Errorf("pair of kings should trump pair of twos")
	}
}

func TestTie(t *testing.T) {
	firstHand := DeckOf(
		CardOf(Clubs, King),
		CardOf(Spades, King),
		CardOf(Spades, Five),
		CardOf(Spades, Seven),
		CardOf(Diamonds, Jack),
	)
	secondHand := DeckOf(
		CardOf(Hearts, King),
		CardOf(Diamonds, King),
		CardOf(Spades, Five),
		CardOf(Spades, Seven),
		CardOf(Diamonds, Jack),
	)
	firstHandRecognised, _ := RecogniseHand(firstHand)
	secondHandRecognised, _ := RecogniseHand(secondHand)
	result := CompareHands(firstHandRecognised, secondHandRecognised)
	if result != Tie {
		t.Errorf("pair of kings is equal to pair of kings with the same kickers")
	}
}

func testPokerHandRecognition(t *testing.T, testedHand Deck, referenceHand Hand) {
	recognisedHand, _ := RecogniseHand(testedHand)
	if recognisedHand.handType != referenceHand.handType {
		t.Errorf("RecogniseHand type should be %s (was %s)", referenceHand.handType.String(), recognisedHand.handType.String())
	}
	for i, rank := range referenceHand.comparison {
		if rank != recognisedHand.comparison[i] {
			t.Errorf("RecogniseHand comparison[%d] should be %s (was %s)", i, rank.String(), recognisedHand.comparison[i].String())
		}
	}
}
