package game

import (
	"math/rand"
	"time"
)

type Suits int

type CardRanks int


var suits = []string{"Clubs", "Diamonds", "Hearts", "Spades"}
var ranks = []string{"", "", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King", "Ace"}

const (
	Clubs Suits = iota
	Diamonds
	Hearts
	Spades
)

const (
	Two CardRanks = 2 + iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

type Card struct {
	Suit Suits
	Rank CardRanks
}

type Deck []Card

const NumOfDecks = 6

func NewDeck() Deck {
	deck := make(Deck, 52*NumOfDecks)

	for i := 0; i < NumOfDecks; i++ {
		for j := Clubs; j <= Spades; j++ {
			for k := Two; k <= Ace; k++ {
				deck[(int(j)*13+int(k-2))+(52*i)] = Card{
					Suit: j,
					Rank: k,
				}
			}
		}
	}
	return deck
}

func (deck Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
}

func (card *Card) GetCardName() string {
	return suits[card.Suit] + ": " + ranks[card.Rank]
}
