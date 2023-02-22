package ws

import "fmt"

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cli := <-h.Register:
			if val, ok := h.Rooms[cli.RoomID]; ok {
				if _, ok := val.Clients[cli.ID]; !ok {
					val.Clients[cli.ID] = cli
				}
			}
		case cli := <-h.Unregister:
			if val, ok := h.Rooms[cli.RoomID]; ok {
				if _, ok := val.Clients[cli.ID]; ok {
					if len(val.Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  fmt.Sprintf("%v left the chat", cli.Username),
							RoomID:   cli.RoomID,
							Username: cli.Username,
						}
					}

					delete(val.Clients, cli.ID)
					close(cli.Message)
				}
			}

		case msg := <-h.Broadcast:
			if val, ok := h.Rooms[msg.RoomID]; ok {
				for _, cli := range val.Clients {
					cli.Message <- msg
				}
			}
		}
	}
}
