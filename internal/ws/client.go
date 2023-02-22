package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

func (cli *Client) writeMessage() {
	defer func() {
		cli.Conn.Close()
	}()

	for {
		message, ok := <-cli.Message
		if !ok {
			return
		}
		cli.Conn.WriteJSON(message)
	}
}

func (cli *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- cli
		cli.Conn.Close()
	}()

	for {
		_, msg, err := cli.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err.Error())
			}
			break
		}
		message := &Message{
			Content:  string(msg),
			RoomID:   cli.RoomID,
			Username: cli.Username,
		}
		hub.Broadcast <- message
	}
}
