package poker

import (
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

func hand(deck cards.Deck) Hand {
	occurrences := buildOrderedOccurrencesSlice(deck)
	{
		// three of a kind
		if occurrences[0].Occurrences == 3 {
			return buildHandWithKickers(1, occurrences, ThreeOfAKind)
		}
		// some pairs
		if occurrences[0].Occurrences == 2 {
			// two pair
			if occurrences[1].Occurrences == 2 {
				return buildHandWithKickers(2, occurrences, TwoPair)
			} else {
				// one pair
				return buildHandWithKickers(1, occurrences, OnePair)
			}
		}
		return buildHandWithKickers(0, occurrences, HighCard)
	}
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

func buildHandWithKickers(kickerStart int, occurrences []rankOccurrences, handType HandType) Hand {
	comparison := make([]cards.Rank, 0)
	for i := 0; i < kickerStart; i++ {
		comparison = append(comparison, occurrences[i].Rank)
	}
	comparison = createAndAppendKickers(occurrences, kickerStart, comparison)
	return Hand{handType: handType, comparison: comparison}
}

func createAndAppendKickers(occurrences []rankOccurrences, kickersStart int, comparison []cards.Rank) []cards.Rank {
	kickers := occurrences[kickersStart:]
	sort.Slice(kickers, func(i, j int) bool {
		return kickers[i].Rank > kickers[j].Rank
	})
	for _, kicker := range kickers {
		comparison = append(comparison, kicker.Rank)
	}
	return comparison
}
