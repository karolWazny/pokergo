package cards

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"math/rand/v2"
)

type Suit string

const (
	Clubs    Suit = "clubs"
	Hearts   Suit = "hearts"
	Spades   Suit = "spades"
	Diamonds Suit = "diamonds"
)

func (s Suit) String() string {
	return cases.Title(language.English, cases.Compact).String(string(s))
}

type Rank int64

const (
	AltAce Rank = 1
	Two    Rank = 2
	Three  Rank = 3
	Four   Rank = 4
	Five   Rank = 5
	Six    Rank = 6
	Seven  Rank = 7
	Eight  Rank = 8
	Nine   Rank = 9
	Ten    Rank = 10
	Jack   Rank = 11
	Queen  Rank = 12
	King   Rank = 13
	Ace    Rank = 14
)

func (r Rank) String() string {
	switch r {
	case AltAce:
		return "Ace"
	case Two:
		return "Two"
	case Three:
		return "Three"
	case Four:
		return "Four"
	case Five:
		return "Five"
	case Six:
		return "Six"
	case Seven:
		return "Seven"
	case Eight:
		return "Eight"
	case Nine:
		return "Nine"
	case Ten:
		return "Ten"
	case Jack:
		return "Jack"
	case Queen:
		return "Queen"
	case King:
		return "King"
	case Ace:
		return "Ace"
	default:
		return "Unknown"
	}
}

type Card struct {
	suit Suit
	rank Rank
}

func CardOf(s Suit, r Rank) Card {
	return Card{suit: s, rank: r}
}

func (c Card) Suit() Suit {
	return c.suit
}

func (c Card) Rank() Rank {
	return c.rank
}

func (c Card) String() string {
	return fmt.Sprintf("%v of %s", c.rank, c.suit)
}

type Deck struct {
	Cards []Card
}

func (d Deck) String() string {
	return fmt.Sprint(d.Cards)
}

func CreateDeck() Deck {
	var cards []Card

	ranks := []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
	suits := []Suit{Clubs, Hearts, Spades, Diamonds}
	for _, rank := range ranks {
		for _, suit := range suits {
			cards = append(cards, CardOf(suit, rank))
		}
	}
	return Deck{
		Cards: cards,
	}
}

func DeckOf(cards ...Card) Deck {
	return Deck{cards}
}

func (d Deck) Shuffled() Deck {
	shuffled := make([]Card, len(d.Cards))
	copy(shuffled, d.Cards)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return DeckOf(shuffled...)
}

func (d Deck) Deal(cards int) (hand Deck, deck Deck) {
	hand = DeckOf(d.Cards[:cards]...)
	deck = DeckOf(d.Cards[cards:]...)
	return
}
