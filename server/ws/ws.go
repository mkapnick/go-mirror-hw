package ws

import (
	"fmt"
	"github.com/mkapnick/go-mirror-hw/server/auth"
	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"strings"
)

// Message is the main interface used for sending data between client
// and server
type Message struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	// websocketId of the connection
	WebsocketId string `json:"websocketId"`
	// userId of the client who established the connection
	UserId  string `json:"userId"`
	Channel string `json:"channel"`
}

// Client associates a web socket connection with a client (UserId)
type Client struct {
	Connection  *websocket.Conn
	WebsocketId string
	UserId      string
}

// keep track of the clients (map connections to a client)
var clients = make(map[*websocket.Conn]*Client)

// map clients to users
var clientUsers = make(map[*Client]auth.UserInfo)

func HandleConnections(ws *websocket.Conn) {
	defer ws.Close()
	fmt.Println("client connected")

	if clients[ws] == nil {
		fmt.Println("adding new client to connection pool")
		var client *Client = &Client{
			Connection:  ws,
			WebsocketId: uuid.Must(uuid.NewV4(), nil).String(),
		}
		clients[ws] = client

		var message Message = Message{
			Code:        "registered.success",
			Message:     "Registered",
			WebsocketId: client.WebsocketId,
		}

		fmt.Println("sending registered ack to client")

		// send the client their assigned userId. This is our own form of `ACK`.
		err := websocket.JSON.Send(clients[ws].Connection, message)

		if err != nil {
			// something weird happened with the connection if the `ack` fails.
			// remove it from out list
			delete(clients, ws)
			fmt.Println("Error sending ack to client")
			fmt.Println(err)
		}
	}

	for {
		var message Message

		if err := websocket.JSON.Receive(ws, &message); err != nil {
			// client most likely closed the connection and an EOF message
			// was sent. Delete the connection here.
			delete(clients, ws)
			fmt.Println("Error in message sent from client")
			fmt.Println(err)
			break
		}

		fmt.Printf("Message received: %s\n", message.Code)

		// update the client to also contain the `UserId`. The `client` is our
		// struct that associates a web socket connection with a client and user
		if message.Code == "register.userId" {
			for index, client := range clients {
				if client.WebsocketId == message.WebsocketId {
					fmt.Println("FOUND: updating client user id to", message.UserId)
					client.UserId = message.UserId
					clients[index] = client
				}
			}
		}

		if message.Code == "message.direct" {
			fmt.Println("direct message channel:", message.Channel)
			for _, client := range clients {
				fmt.Println("Client in question", client.WebsocketId, client.UserId)
				channel := message.Channel
				tokens := strings.Split(message.Channel, ".")
				// the 0th userId is the initiating userId, and they are already
				// subscribed via the front end client. Therefore we only want to notify
				// the 2nd ([1]) user to subscribe to this new channel to receive
				// DMs
				userIdToNotify := tokens[1]
				fmt.Println("comparing client.UserId == userIdToNotify", client.UserId, userIdToNotify)
				if client.UserId == userIdToNotify {
					err := websocket.JSON.Send(client.Connection, &Message{
						Code:    "subscribe.direct",
						Channel: channel,
					})
					if err != nil {
						// most likely the connection is closed, so remove the connection
						// from our list
						delete(clients, ws)
						fmt.Println("Error sending message back to client")
						fmt.Println(err)
						break
					}
				}
			}
		}
	}
}
