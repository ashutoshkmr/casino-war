package client

type SocketCommand int

type Msg struct {
   Content string `json:"content,omitempty"`
   Command SocketCommand    `json:"command,omitempty"`
   Err     string `json:"err,omitempty"`
}

const (
	StartGame SocketCommand = iota + 1
	PlaceBet
	DrawCard
	DrawResult
	GoToWar
	SurrenderGame
	QuitGame
	GameStarted
	Error
)

type DrawCardsResult struct {
	PlayerCard	string			`json:"playerCard,omitempty"`
	DealerCard	string			`json:"dealerCard,omitempty"`
	PlayerChips	int				`json:"playerChips,omitempty"`
	Status		DrawCardStatus	`json:"status,omitempty"`
}


type DrawCardStatus int

const (
	Won  DrawCardStatus = 1
	Loss DrawCardStatus = -1
	Tie  DrawCardStatus = 0
)