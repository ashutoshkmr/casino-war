package client

import "main/game"

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

type DrawCardResult struct {
	Content string        				`json:"content,omitempty"`
	Command SocketCommand 				`json:"command,omitempty"`
	Err     string        				`json:"err,omitempty"`
	DrawResult 	string					`json:"drawResult,omitempty"`
	Result 		game.DrawCardsResult	`json:"result,omitempty"`
}