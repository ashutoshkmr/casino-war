package main

import (
	"bytes"
	"encoding/json"
	"log"
	"main/game"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var gameMap = make(map[uuid.UUID]game.CasinoTable)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	clientId uuid.UUID

	hub *Hub

	conn *websocket.Conn

	send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		handleMessage(c, message)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	id := uuid.New()
	client := &Client{clientId: id, hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func handleMessage(c *Client, msg []byte) {
	var parsedMsg SocketMessage
	wsResponse := SocketMessage{}
	if err := json.Unmarshal(msg, &parsedMsg); err != nil {
		panic(err)
	}

	table := gameMap[c.clientId]

	log.Println(table)

	switch parsedMsg.Command {
	case StartGame:
		initialChipsAmount, err := strconv.Atoi(parsedMsg.Content)
		if err != nil {
			log.Println(err)
			wsResponse = SocketMessage{Command: Error, Content: "Invalid chips amount"}
		}
		startGame(initialChipsAmount, c.clientId)

		wsResponse = SocketMessage{Command: GameStarted}

	// case PlaceBet:
	// 	log.Println("PlaceBet : ", parsedMsg)
	case DrawCard:
		if !table.IsActive {
			wsResponse = SocketMessage{Command: Error, Content: "Game doesn't exist or is not available anymore"}
			break
		}

		betAmount, err := strconv.Atoi(parsedMsg.Content)

		if err != nil {
			wsResponse = SocketMessage{Command: Error, Content: "Invalid bet amount"}
			break
		}

		if table.CanPlaceBet(betAmount) {
			wsResponse = SocketMessage{}

			table.PlaceBet(betAmount)
			playerCard, dealerCard, result, _ := table.DrawCards()
			// y := DrawCardResult{
			// 	Command:    DrawResult,
			// 	DrawResult: fmt.Sprint(playerCard + "\t" + dealerCard),
			// 	Result:     result,
			// }

			// res, _ := json.Marshal(y)
			// c.hub.broadcast <- []byte(res)

		}

	case SurrenderGame:
		log.Println("SurrenderGame : ", parsedMsg)
	case GoToWar:
		log.Println("GoToWar : ", parsedMsg)
	case QuitGame:
		log.Println("Quit Game")
		wsResponse = SocketMessage{Command: QuitGame}
	default:
		wsResponse = SocketMessage{Command: Error, Content: "Invalid Message"}

	}

	res, err := json.Marshal(wsResponse)

	if err != nil {
		log.Println(err)
	}

	c.hub.broadcast <- []byte(res)

}

func startGame(initialChipsAmount int, clientId uuid.UUID) {
	player := game.Player{Chips: initialChipsAmount, Score: 0}
	// Assuming dealer has enough chips to handle bet of any size
	dealer := game.Player{Chips: math.MaxInt32, Score: 0}
	table := game.NewTable(player, dealer)
	gameMap[clientId] = table
}
