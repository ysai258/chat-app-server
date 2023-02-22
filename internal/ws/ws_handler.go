package ws

import (
	"fmt"
	"net/http"
	"server/internal/constants"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

type CreateRoomRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(c *gin.Context) {

	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}
	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}
	c.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// origin := r.Header.Get("Origin")
		// return origin == constants.BASE_DOMAIN
		return true
	},
}

func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.BadRequest(c, err)
		return
	}

	roomId := c.Param("roomId")
	tkn, ok := c.Get(constants.JWT_TOKEN_CLAIMS_KEY)
	if !ok {
		utils.ServerError(c, fmt.Errorf("invalid token"))
		return
	}
	tokenData := tkn.(*constants.TokenClaims)

	clientID := tokenData.ID
	username := tokenData.Username

	cli := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomId,
		Username: username,
	}
	msg := &Message{
		Content:  fmt.Sprintf("%v has joined the room", username),
		RoomID:   roomId,
		Username: username,
	}

	h.hub.Register <- cli

	h.hub.Broadcast <- msg

	go cli.writeMessage()

	cli.readMessage(h.hub)
}

type RoomResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetRooms(c *gin.Context) {
	rooms := make([]RoomResponse, 0)
	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomResponse{
			ID:   r.ID,
			Name: r.Name,
		})
	}
	c.JSON(http.StatusOK, rooms)
}

type CLientResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(c *gin.Context) {
	clients := make([]CLientResponse, 0)
	roomId := c.Param("roomId")
	if _, ok := h.hub.Rooms[roomId]; !ok {
		c.JSON(http.StatusOK, clients)
		return
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, CLientResponse{
			ID:       c.ID,
			Username: c.Username,
		})
	}
	c.JSON(http.StatusOK, clients)
}
