package sockets

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
 ReadBufferSize:  1024,
 WriteBufferSize: 1024,
 CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	feed chan string
	actions chan *Action
	conn *websocket.Conn
	hub  *Hub
}


func (client *Client) sendToEveryone(message string)  {
	for c := range client.hub.clients {
		c.feed <- message
	}
}

// listens to any changes in the feed and action channel
func (client *Client) feedHandler() {
	for {
		select {
		case action, channelOpened := <- client.actions:
			err := client.conn.WriteJSON(action)
			if !channelOpened {
				return
			}

			if err != nil {
				log.Printf("err: %s", err)
				return
			}

		case message, channelOpened := <- client.feed:
			if !channelOpened {
				return
			}
			res := &Action {
				Event: "message",
				Content: message,
			}
			err := client.conn.WriteJSON(res)
			if err != nil {
				log.Printf("err: %s", err)
				return
			}
		}
	}
}

// reads if there was a message sent in the connection
func (client *Client) reader() {
	defer func() {	
		close(client.feed)	
		close(client.actions)
		client.hub.unregister <- client
		client.conn.Close()	
	}()

	for {
		var action Action

		err := client.conn.ReadJSON(&action)
		if err != nil {
			log.Printf("err: %s", err)
			break
		}

		handleActions(client, &action)
	}
}

func HandleConnections(w http.ResponseWriter, r *http.Request, hub *Hub) {
	// upgrades the HTTP server connection to the WebSocket protocol
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("err: %s", err)
		return
	}

	client := &Client{feed: make(chan string), actions: make(chan *Action), conn: ws, hub: hub}
		
	go client.reader()
	go client.feedHandler()
}