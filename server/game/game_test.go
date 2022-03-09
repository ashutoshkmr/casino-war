package game

import (
	"math"
	"testing"

	"github.com/google/uuid"
)

var totalChips = 1000
var betAmount = 10
var clientId = uuid.New()
var Palyer = Player{Chips: totalChips, Score: 0}
var Dealer = Player{Chips: math.MaxInt32, Score: 0}
var table = NewTable(Palyer, Dealer, clientId)

func TestBet(t *testing.T) {
	table.PlaceBet(betAmount)
	expectedBet := betAmount
	actualBet := table.Bet
	if expectedBet != actualBet {
		t.Errorf("got %q, wanted %q", actualBet, expectedBet)
	}
}

func TestChips(t *testing.T) {
	expectedChips := totalChips - betAmount
	actualChips := table.Player.Chips

	if expectedChips != actualChips {
		t.Errorf("got %q, wanted %q", actualChips, expectedChips)
	}
}

func TestDrawCard(t *testing.T) {
	// Arrange
	playerCard := table.Deck[len(table.Deck)-1]
	dealerCard := table.Deck[len(table.Deck)-2]

	len1 := len(table.Deck)

	res, _ := table.DrawCards()

	len2 := len(table.Deck)

	if playerCard.GetCardName() != res.PlayerCard {
		t.Errorf("got %q, wanted %q", res.PlayerCard, playerCard.GetCardName())
	}

	if dealerCard.GetCardName() != res.DealerCard {
		t.Errorf("got %q, wanted %q", res.DealerCard, dealerCard.GetCardName())
	}

	if len1-2 != len2 {
		t.Errorf("got %q, wanted %q", len2, len1-2)
	}

	if playerCard.Rank == dealerCard.Rank && res.Status != Tie {
		t.Errorf("got %q, wanted %q", res.Status, Tie)
	}

	if playerCard.Rank > dealerCard.Rank && res.Status != Won {
		t.Errorf("got %q, wanted %q", res.Status, Won)
	}
	if playerCard.Rank < dealerCard.Rank && res.Status != Loss {
		t.Errorf("got %q, wanted %q", res.Status, Loss)
	}
}

func TestSurrender(t *testing.T) {
	table.PlaceBet(betAmount)
	chipsBeforeSurrender := table.Player.Chips
	table.Surrender()
	chipsAfterSurrender := table.Player.Chips

	if chipsBeforeSurrender+betAmount/2 != chipsAfterSurrender {
		t.Errorf("got %q, wanted %q", chipsAfterSurrender, chipsBeforeSurrender+betAmount/2)
	}
}
