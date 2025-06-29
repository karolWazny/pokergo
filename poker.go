package pokergo

import (
	"errors"
	"fmt"
	"sort"
)

type HandType int

const (
	LowestHandGuardian HandType = -1
	HighCard           HandType = iota
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

type ComparisonResult int

const (
	FirstWins  ComparisonResult = -1
	Tie                         = 0
	SecondWins                  = 1
)

type Hand struct {
	handType   HandType
	comparison []Rank
}

func (h Hand) String() string {
	return h.handType.String() + " " + fmt.Sprintf("%v", h.comparison)
}

func (h HandType) String() string {
	switch h {
	case HighCard:
		return "HighCard"
	case OnePair:
		return "OnePair"
	case TwoPair:
		return "TwoPair"
	case ThreeOfAKind:
		return "ThreeOfAKind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "FullHouse"
	case FourOfAKind:
		return "FourOfAKind"
	case StraightFlush:
		return "StraightFlush"
	case RoyalFlush:
		return "RoyalFlush"
	default:
		return "HighCard"
	}
}

func CompareHands(first Hand, second Hand) ComparisonResult {
	if first.handType > second.handType {
		return FirstWins
	} else if first.handType < second.handType {
		return SecondWins
	}
	for i, card := range first.comparison {
		if card > second.comparison[i] {
			return FirstWins
		} else if card < second.comparison[i] {
			return SecondWins
		} else {
			continue
		}
	}
	return Tie
}

func RecogniseHand(deck Deck) (Hand, error) {
	if len(deck.Cards) != 5 {
		return Hand{}, errors.New("invalid poker RecogniseHand size")
	}
	occurrences := buildOrderedOccurrencesSlice(deck)
	isFlush := isFlush(deck)
	{
		isStraight, occurrences := isStraight(occurrences)
		if isStraight && isFlush && occurrences[0].Rank == Ace {
			return buildStraightHand(RoyalFlush), nil
		}
		if isStraight && isFlush {
			return buildStraightHand(StraightFlush, occurrences[0].Rank), nil
		}
		if occurrences[0].Occurrences == 4 {
			return buildHandWithKickers(occurrences, FourOfAKind), nil
		}
		if occurrences[0].Occurrences == 3 && occurrences[1].Occurrences == 2 {
			return buildHandWithKickers(occurrences, FullHouse), nil
		}
		if isFlush {
			return buildHandWithKickers(occurrences, Flush), nil
		}
		if isStraight {
			return buildStraightHand(Straight, occurrences[0].Rank), nil
		}
		if occurrences[0].Occurrences == 3 {
			return buildHandWithKickers(occurrences, ThreeOfAKind), nil
		}
		if occurrences[0].Occurrences == 2 {
			if occurrences[1].Occurrences == 2 {
				return buildHandWithKickers(occurrences, TwoPair), nil
			} else {
				return buildHandWithKickers(occurrences, OnePair), nil
			}
		}
		return buildHandWithKickers(occurrences, HighCard), nil
	}
}

func CreateLowGuardian() Hand {
	return Hand{
		handType:   LowestHandGuardian,
		comparison: []Rank{},
	}
}

type rankOccurrences struct {
	Rank        Rank
	Occurrences int
}

func buildStraightHand(handType HandType, comparisons ...Rank) Hand {
	return Hand{
		handType:   handType,
		comparison: comparisons,
	}
}

func isFlush(deck Deck) bool {
	suits := map[Suit]bool{}
	for _, card := range deck.Cards {
		suits[card.Suit()] = true
	}
	return len(suits) == 1
}

func isStraight(occurrences []rankOccurrences) (bool, []rankOccurrences) {
	isStraight := true
	for i := range len(occurrences) - 1 {
		if occurrences[i].Rank == Ace && occurrences[i+1].Rank == Five {
			continue
		}
		if int(occurrences[i].Rank-occurrences[i+1].Rank) != 1 {
			isStraight = false
			break
		}
	}
	if isStraight && occurrences[0].Rank == Ace && occurrences[1].Rank == Five {
		return isStraight, append(occurrences[1:], occurrences[0])
	}
	return isStraight, occurrences
}

func buildOrderedOccurrencesSlice(deck Deck) []rankOccurrences {
	unique := map[Rank][]Suit{}
	for _, card := range deck.Cards {
		unique[card.Rank()] = append(unique[card.Rank()], card.Suit())
	}
	occurrences := make([]rankOccurrences, len(unique))
	{
		i := 0
		for rank, suits := range unique {
			occurrences[i] = rankOccurrences{rank, len(suits)}
			i++
		}
	}
	sort.Slice(occurrences, func(i, j int) bool {
		return occurrences[i].Rank > occurrences[j].Rank
	})
	sort.SliceStable(occurrences, func(i, j int) bool {
		return occurrences[i].Occurrences > occurrences[j].Occurrences
	})
	return occurrences
}

func buildHandWithKickers(occurrences []rankOccurrences, handType HandType) Hand {
	comparison := make([]Rank, len(occurrences))

	for i, occurrence := range occurrences {
		comparison[i] = occurrence.Rank
	}
	return Hand{handType: handType, comparison: comparison}
}
