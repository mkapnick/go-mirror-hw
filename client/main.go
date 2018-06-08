package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/websocket"
	"os"
	"time"
)

// This file is not used in the app - I originally created it to test
// my websocket server.

var shouldAllowPrompt bool = true

type Message struct {
	Message string `json:"message"`
}

func main() {
	fmt.Println("Starting client")

	address := fmt.Sprintf("ws://%s:%s/ws", "localhost", "9000")
	ws, err := websocket.Dial(address, "", fmt.Sprintf("http://%s/", "localhost:9000"))

	if err != nil {
		fmt.Printf("Dial failed: %s\n", err.Error())
		os.Exit(1)
	}

	msgChan := make(chan Message)
	go readClientMessages(ws, msgChan)

	// block here, wait for messages from the server
	for {
		select {
		case <-time.After(time.Duration(2e9)):
			go func() {
				if shouldAllowPrompt {
					shouldAllowPrompt = false
					reader := bufio.NewReader(os.Stdin)
					fmt.Print("Message: ")
					text, _ := reader.ReadString('\n')

					// set shouldAllowPrompt it to true

					msg := Message{
						Message: text,
					}
					err = websocket.JSON.Send(ws, msg)
					if err != nil {
						fmt.Printf("Error sending data to the server %s\n", err.Error())
						os.Exit(1)
					}
					// set it back to false
					shouldAllowPrompt = true
				}
			}()
		case message := <-msgChan:
			fmt.Println(message.Message)
		}
	}
}

// we wan run as many clients as we want, each client has its own
// ws connection with the server
func readClientMessages(ws *websocket.Conn, msgChan chan Message) {
	for {
		var message Message

		// check to see if there's anything shouldAllowPrompt from the server and if so
		// send to our channel
		err := websocket.JSON.Receive(ws, &message)

		if err != nil {
			fmt.Printf("Error on receiving json on web socket connection %s\n", err.Error())
			return
		}

		msgChan <- message
	}
}
