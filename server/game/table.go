package game

type DrawCardsResult int

const (
	Won  DrawCardsResult = 1
	Loss DrawCardsResult = -1
	Tie  DrawCardsResult = 0
)

type CasinoTable struct {
	Deck   Deck
	Player Player
	Dealer Player
	Bet    int
	IsActive bool
}

func NewTable(player, dealer Player) CasinoTable {
	deck := NewDeck()
	deck.Shuffle()
	table := CasinoTable{
		deck,
		player,
		dealer,
		0,
		true,
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

func (table *CasinoTable) DrawCards() (string, string, DrawCardsResult, error) {
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

	return playerCard.GetCardName(), dealerCard.GetCardName(), result, nil
}

func (table *CasinoTable) Surrender() {
	table.Player.Chips += table.Bet / 2
	table.Dealer.Chips += table.Bet / 2
	table.Bet = 0
}

func (table *CasinoTable) GoToWar() (string, string, DrawCardsResult, error) {
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
	}
	table.Bet = 0

	return playerCard.GetCardName(), dealerCard.GetCardName(), result, nil

}

func compareCards(first, second Card) DrawCardsResult {
	if first.Rank < second.Rank {
		return Loss
	}
	if first.Rank > second.Rank {
		return Won
	}
	return Tie
}
