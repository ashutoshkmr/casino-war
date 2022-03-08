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
var colorBlue = "\033[34m"
var colorPurple = "\033[35m"
var colorCyan = "\033[36m"
var colorWhite = "\033[37m"

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
			break
		} else {
			// get bet amount
			chips, err := client.ChipsPrompt(os.Stdin, os.Stdout)

			msg := client.Msg{
				Command: client.DrawCard,
				Content: fmt.Sprint(chips),
			}
			j, _ := json.Marshal(msg)
			c.WriteMessage(websocket.TextMessage, []byte(j))

			if err != nil {
				fmt.Println("Invalid chip amount")
				return
			}
		}

	case client.QuitGame:
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Fatal("write close:", err)
		}
		c.Close()
	case client.DrawResult:
		var parsedResult client.DrawCardResult
		if err := json.Unmarshal(msg, &parsedResult); err != nil {
			panic(err)
		}
		log.Println("DrawResult : ", parsedResult)

	}

	log.Print(parsedMsg)

}

// import (
// 	"fmt"
// 	"main/client"
// 	"main/game"
// 	"math"
// 	"os"
// )

// var colorReset = "\033[0m"

// var colorRed = "\033[31m"
// var colorGreen = "\033[32m"
// var colorYellow = "\033[33m"
// var colorBlue = "\033[34m"
// var colorPurple = "\033[35m"
// var colorCyan = "\033[36m"
// var colorWhite = "\033[37m"

// func main() {

// 	player := game.Player{Chips: 500, Score: 0}
// 	dealer := game.Player{Chips: math.MaxInt32, Score: 0}

// 	table := game.NewTable(player, dealer)

// 	deckLen := len(table.Deck)

// 	fmt.Println(player)

// 	for deckLen > 0 {
// 		result, err := client.PromptToDraw(os.Stdin, os.Stdout)

// 		if err != nil {
// 			fmt.Println(err)
// 			break
// 		}
// 		if err != nil {
// 			fmt.Println(err)
// 			break
// 		}
// 		if result == client.Quit {
// 			printGameStats(table.Player, table.Dealer, 0)
// 			break
// 		} else {
// 			chips, err := client.ChipsPrompt(os.Stdin, os.Stdout)

// 			if err != nil {
// 				fmt.Println("Invalid chip amount")
// 				return
// 			}

// 			table.PlaceBet(chips)

// 			playerCard, dealerCard, result, err := table.DrawCards()

// 			if err != nil {
// 				fmt.Println(err)
// 				break
// 			}

// 			fmt.Println(string(colorCyan), "Cards drawn : ", playerCard, dealerCard)
// 			fmt.Println(string(colorReset))

// 			if result == game.Tie {
// 				// A tie has occured
// 				r, err := client.InitTiePrompt(os.Stdin, os.Stdout)
// 				if err != nil {
// 					return
// 				}

// 				if r == client.Surrender {
// 					fmt.Println("You surrendered")
// 					table.Surrender()
// 					printGameStats(table.Player, table.Dealer, 0)
// 					break
// 				}

// 				if r == client.Gotowar {
// 					fmt.Println("Go to war")
// 					if canGoToWar(table) {
// 						playerCard, dealerCard, result, _ := table.GoToWar()
// 						fmt.Println(string(colorCyan), "Cards drawn : ", playerCard, dealerCard)
// 						fmt.Println(string(colorReset))

// 						if result == game.Tie {
// 							fmt.Println(string(colorYellow), "It's a tie again, you win 10x")
// 							result = game.Won
// 						}

// 						printGameStats(table.Player, table.Dealer, result)

// 					}
// 				}

// 			} else {
// 				printGameStats(table.Player, table.Dealer, result)
// 			}

// 			deckLen = len(table.Deck)
// 			if deckLen < 2 {
// 				fmt.Println("Deck doesn't have enough cards")
// 				break
// 			}
// 		}
// 	}
// }

// func printGameStats(player, dealer game.Player, drawResult game.DrawCardsResult) {
// 	if drawResult == 1 {
// 		fmt.Println(string(colorGreen), "You won!", string(colorReset))
// 	}
// 	if drawResult == -1 {
// 		fmt.Println(string(colorRed), "You lost.", string(colorReset))
// 	}
// 	fmt.Printf("Player Chips : %d  \t\t Score : %d\n", player.Chips, player.Score)
// 	fmt.Printf("Dealer Chips : %d  \t Score : %d\n\n", dealer.Chips, dealer.Score)
// }

// func canGoToWar(table game.CasinoTable) bool {
// 	return table.Player.Chips >= table.Bet
// }
