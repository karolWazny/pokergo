package poker

import (
	"errors"
	"online-poker/cards"
	"sort"
)

type HandType int

const (
	HighCard HandType = iota
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

type Hand struct {
	handType   HandType
	comparison []cards.Rank
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

type rankOccurrences struct {
	Rank        cards.Rank
	Occurrences int
}

func hand(deck cards.Deck) (Hand, error) {
	if len(deck.Cards) != 5 {
		return Hand{}, errors.New("invalid poker hand size")
	}
	occurrences := buildOrderedOccurrencesSlice(deck)
	isFlush := isFlush(deck)
	{
		isStraight, occurrences := isStraight(occurrences)
		if isStraight && isFlush && occurrences[0].Rank == cards.Ace {
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

func buildStraightHand(handType HandType, comparisons ...cards.Rank) Hand {
	return Hand{
		handType:   handType,
		comparison: comparisons,
	}
}

func isFlush(deck cards.Deck) bool {
	suits := map[cards.Suit]bool{}
	for _, card := range deck.Cards {
		suits[card.Suit()] = true
	}
	return len(suits) == 1
}

func isStraight(occurrences []rankOccurrences) (bool, []rankOccurrences) {
	isStraight := true
	for i := range len(occurrences) - 1 {
		if occurrences[i].Rank == cards.Ace && occurrences[i+1].Rank == cards.Five {
			continue
		}
		if int(occurrences[i].Rank-occurrences[i+1].Rank) != 1 {
			isStraight = false
			break
		}
	}
	if isStraight && occurrences[0].Rank == cards.Ace && occurrences[1].Rank == cards.Five {
		return isStraight, append(occurrences[1:], occurrences[0])
	}
	return isStraight, occurrences
}

func buildOrderedOccurrencesSlice(deck cards.Deck) []rankOccurrences {
	unique := map[cards.Rank][]cards.Suit{}
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
	comparison := make([]cards.Rank, len(occurrences))

	for i, occurrence := range occurrences {
		comparison[i] = occurrence.Rank
	}
	return Hand{handType: handType, comparison: comparison}
}
