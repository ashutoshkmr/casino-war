package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"main/client"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var colorReset = "\033[0m"

var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var colorYellow = "\033[33m"
var colorCyan = "\033[36m"
// var colorBlue = "\033[34m"
// var colorPurple = "\033[35m"
// var colorWhite = "\033[37m"

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		chipsAmount, err := client.ChipsPrompt(os.Stdin, os.Stdout)

		if err != nil {
			log.Print(err)
		}

		startGameMsg := client.Msg{
			Content: fmt.Sprint(chipsAmount),
			Command: client.StartGame,
			Err:     "",
		}
		j, _ := json.Marshal(startGameMsg)
		wserr := c.WriteMessage(websocket.TextMessage, []byte(j))

		if wserr != nil {
			log.Println("write:", err)
			return
		}
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			handleWsMessages(c, message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func handleWsMessages(c *websocket.Conn, msg []byte) {
	var parsedMsg client.Msg

	if err := json.Unmarshal(msg, &parsedMsg); err != nil {
		panic(err)
	}

	switch parsedMsg.Command {
	case client.GameStarted:
		promptDraw(c)
	case client.QuitGame:
		fmt.Println(string(colorCyan), "Chips : ", parsedMsg.Content, string(colorReset))
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Fatal("write close:", err)
		}
		// print player stats
		c.Close()
		os.Exit(1)
	case client.DrawResult:
		var drawcardRes client.DrawCardsResult
		if err := json.Unmarshal([]byte(parsedMsg.Content), &drawcardRes); err != nil {
			panic(err)
		}

		fmt.Println(string(colorCyan), "Cards drawn : ", drawcardRes.PlayerCard, drawcardRes.DealerCard, string(colorReset))
		printDrawResult(drawcardRes.Status)
		fmt.Println(string(colorCyan), "Chips : ", drawcardRes.PlayerChips, string(colorReset))


		if drawcardRes.Status == client.Tie {
			promtTie(c)
		} else {
			promptDraw(c)
		}
	}
}

func printDrawResult(result client.DrawCardStatus) {
	if result == client.Won {
		fmt.Println(string(colorGreen), "You won!", string(colorReset))
	} else if result == client.Loss {

		fmt.Println(string(colorRed), "You Lost", string(colorReset))
	} else if result == client.Tie {
		fmt.Println(string(colorYellow), "It's a Tie", string(colorReset))
	} else {
		fmt.Println(string(colorReset))
	}
}

func promptDraw(c *websocket.Conn) {
	result, err := client.PromptToDraw(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	if result == client.Quit {
		msg := client.Msg{
			Command: client.QuitGame,
			Err:     "",
		}
		j, _ := json.Marshal(msg)
		c.WriteMessage(websocket.TextMessage, []byte(j))
	} else {
		// get bet amount
		chips, err := client.BetPrompt(os.Stdin, os.Stdout)

		if err != nil {
			fmt.Println("Invalid chip amount")
			return
		}

		msg := client.Msg{
			Command: client.DrawCard,
			Content: fmt.Sprint(chips),
		}
		j, _ := json.Marshal(msg)
		c.WriteMessage(websocket.TextMessage, []byte(j))
	}
}

func promtTie(c *websocket.Conn) {
	r, err := client.InitTiePrompt(os.Stdin, os.Stdout)
	if err != nil {
		return
	}

	msg := client.Msg{}

	if r == client.Surrender {
		fmt.Println("You surrendered")
		msg = client.Msg{
			Command: client.SurrenderGame,
		}
	} else if r == client.Gotowar {
		fmt.Println("Go to war")
		msg = client.Msg{
			Command: client.GoToWar,
		}
	}

	j, _ := json.Marshal(msg)
	c.WriteMessage(websocket.TextMessage, []byte(j))
}

