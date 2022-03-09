package game

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DrawCardsResult struct {
	PlayerCard  string         `json:"playerCard,omitempty"`
	DealerCard  string         `json:"dealerCard,omitempty"`
	PlayerChips int            `json:"playerChips,omitempty"`
	Status      DrawCardStatus `json:"status,omitempty"`
}

type DrawCardStatus int

const (
	Won  DrawCardStatus = 1
	Loss DrawCardStatus = -1
	Tie  DrawCardStatus = 0
)

type CasinoTable struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Deck     Deck			`bson:"deck,omitempty"`
	Player   Player			`bson:"player,omitempty"`
	Dealer   Player			`bson:"dealer,omitempty"`
	Bet      int			`bson:"bet,omitempty"`
	GameId   uuid.UUID		`bson:"gameId,omitempty"`
}

func NewTable(player, dealer Player, gameId uuid.UUID) CasinoTable {
	deck := NewDeck()
	deck.Shuffle()
	table := CasinoTable{
		Deck :deck,
		Player:player,
		Dealer: dealer,
		Bet: 0,
		GameId: gameId,
	}
	return table
}

func (table *CasinoTable) PlaceBet(betAmount int) {
	table.Player.Chips -= betAmount
	table.Bet = betAmount
}

func (table *CasinoTable) CanPlaceBet(betAmount int) bool {
	return table.Player.Chips >= betAmount
}

func (table *CasinoTable) DrawCards() (DrawCardsResult, error) {
	playerCard := table.Deck[len(table.Deck)-1]
	dealerCard := table.Deck[len(table.Deck)-2]
	table.Deck = table.Deck[:len(table.Deck)-2]

	result := compareCards(playerCard, dealerCard)

	if result == Won {
		table.Player.Chips += table.Bet * 2
		table.Dealer.Chips -= table.Bet
		table.Bet = 0
		table.Player.Score++
	}

	if result == Loss {
		table.Dealer.Chips += table.Bet
		table.Bet = 0
		table.Dealer.Score++
	}

	return DrawCardsResult{
		playerCard.GetCardName(),
		dealerCard.GetCardName(),
		table.Player.Chips,
		result,
	}, nil
}

func (table *CasinoTable) Surrender() {
	table.Player.Chips += table.Bet / 2
	table.Dealer.Chips += table.Bet / 2
	table.Bet = 0
}

func (table *CasinoTable) CanGoToWar() bool {
	return table.Player.Chips >= table.Bet
}

func (table *CasinoTable) GoToWar() (DrawCardsResult, error) {
	warBet := table.Bet * 3
	table.Player.Chips -= table.Bet
	table.Dealer.Chips -= table.Bet

	// Burn three cards

	table.Deck = table.Deck[:len(table.Deck)-3]

	playerCard := table.Deck[len(table.Deck)-1]

	// Burn three cards

	table.Deck = table.Deck[:len(table.Deck)-4]

	dealerCard := table.Deck[len(table.Deck)-1]

	table.Deck = table.Deck[:len(table.Deck)-4]

	result := compareCards(playerCard, dealerCard)

	if result == Won {
		table.Player.Chips += warBet
		table.Player.Score++
	}

	if result == Loss {
		table.Dealer.Chips += warBet
		table.Dealer.Score++
	} else {
		table.Player.Chips += table.Bet * 10
		table.Dealer.Chips -= table.Bet * 10
		table.Player.Score++
		// Won by second tie, can be displayed better on client
		result = Won
	}
	table.Bet = 0

	return DrawCardsResult{
		playerCard.GetCardName(),
		dealerCard.GetCardName(),
		table.Player.Chips,
		result,
	}, nil
}

func compareCards(first, second Card) DrawCardStatus {
	if first.Rank < second.Rank {
		return Loss
	}
	if first.Rank > second.Rank {
		return Won
	}
	return Tie
}
