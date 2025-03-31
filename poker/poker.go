package poker

import (
	"online-poker/cards"
	"slices"
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
	unique := map[cards.Rank][]cards.Suit{}
	for _, card := range deck.Cards {
		unique[card.Rank()] = append(unique[card.Rank()], card.Suit())
	}
	comparison := make([]cards.Rank, len(unique))
	occurrences := make([]rankOccurrences, len(unique))
	i := 0
	for rank, suits := range unique {
		occurrences[i] = rankOccurrences{rank, len(suits)}
		comparison[i] = rank
		i++
	}
	sort.Slice(occurrences, func(i, j int) bool {
		return occurrences[i].Rank > occurrences[j].Rank
	})
	sort.SliceStable(occurrences, func(i, j int) bool {
		return occurrences[i].Occurrences > occurrences[j].Occurrences
	})
	for i, occurrence := range occurrences {
		// some pairs
		if occurrence.Occurrences == 2 {
			// two pair
			if occurrences[i+1].Occurrences == 2 {
				kickers := occurrences[i+2:]
				sort.Slice(kickers, func(i, j int) bool {
					return kickers[i].Rank > kickers[j].Rank
				})
				comparison = make([]cards.Rank, 0)
				comparison = append(comparison, occurrence.Rank)
				comparison = append(comparison, occurrences[i+1].Rank)
				for _, kicker := range kickers {
					comparison = append(comparison, kicker.Rank)
				}
				return Hand{handType: TwoPair, comparison: comparison}
			} else {
				// one pair
				kickers := occurrences[i+1:]
				sort.Slice(kickers, func(i, j int) bool {
					return kickers[i].Rank > kickers[j].Rank
				})
				comparison = make([]cards.Rank, 0)
				comparison = append(comparison, occurrence.Rank)
				for _, kicker := range kickers {
					comparison = append(comparison, kicker.Rank)
				}
				return Hand{handType: OnePair, comparison: comparison}
			}
		}
	}

	slices.Sort(comparison)
	slices.Reverse(comparison)
	if len(deck.Cards) > 1 && deck.Cards[0].Rank() == deck.Cards[1].Rank() {
		slices.Reverse(comparison)
		return Hand{handType: OnePair, comparison: comparison}
	}
	return Hand{handType: HighCard, comparison: comparison}
}
