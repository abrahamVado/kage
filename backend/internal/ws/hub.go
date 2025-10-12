package ws

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Role identifies the type of actor participating in the hub.
type Role string

const (
	RoleRider  Role = "rider"
	RoleDriver Role = "driver"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

// Message envelopes data exchanged through the hub.
type Message struct {
	RoomID  string
	Role    Role
	Type    string
	Payload interface{}
}

// Client maintains websocket state for a single connection.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan interface{}
	room string
	role Role
	mu   sync.Mutex
}

// Hub orchestrates rider and driver communication.
type Hub struct {
	// 1.- Manage client registration and message broadcasting via channels.
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
	shutdown   chan struct{}
	rooms      map[string]map[*Client]struct{}
	logger     *log.Logger
	mu         sync.RWMutex
}

// NewHub constructs a hub with its background goroutine.
func NewHub(logger *log.Logger) *Hub {
	if logger == nil {
		logger = log.Default()
	}
	h := &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
		shutdown:   make(chan struct{}),
		rooms:      make(map[string]map[*Client]struct{}),
		logger:     logger,
	}
	go h.loop()
	return h
}

// RegisterRoutes mounts REST and websocket handlers onto Gin.
func (h *Hub) RegisterRoutes(router *gin.Engine) {
	// 2.- Expose websocket upgrades and health endpoints.
	router.GET("/ws/:role/:room", func(c *gin.Context) {
		role := Role(c.Param("role"))
		room := c.Param("room")
		h.handleUpgrade(c.Writer, c.Request, role, room)
	})

	router.GET("/ws/rooms/:room/occupants", func(c *gin.Context) {
		room := c.Param("room")
		c.JSON(http.StatusOK, gin.H{"room": room, "occupants": h.count(room)})
	})
}

// Broadcast delivers a message to all participants in the specified room.
func (h *Hub) Broadcast(msg Message) {
	h.broadcast <- msg
}

// Shutdown stops the hub loop gracefully.
func (h *Hub) Shutdown(ctx context.Context) {
	close(h.shutdown)
	done := make(chan struct{})
	go func() {
		h.loopDrain()
		close(done)
	}()

	select {
	case <-ctx.Done():
	case <-done:
	}
}

func (h *Hub) handleUpgrade(w http.ResponseWriter, r *http.Request, role Role, room string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Printf("upgrade error: %v", err)
		return
	}
	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan interface{}, 8),
		room: room,
		role: role,
	}
	h.register <- client
	go client.writePump()
	client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(1 << 16)
	_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})
	for {
		var payload interface{}
		if err := c.conn.ReadJSON(&payload); err != nil {
			break
		}
		c.hub.Broadcast(Message{RoomID: c.room, Role: c.role, Type: "update", Payload: payload})
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(50 * time.Second)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case payload, ok := <-c.send:
			if !ok {
				return
			}
			c.mu.Lock()
			if err := c.conn.WriteJSON(payload); err != nil {
				c.mu.Unlock()
				return
			}
			c.mu.Unlock()
		case <-ticker.C:
			c.mu.Lock()
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.mu.Unlock()
				return
			}
			c.mu.Unlock()
		}
	}
}

func (h *Hub) loop() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
		case client := <-h.unregister:
			h.removeClient(client)
		case msg := <-h.broadcast:
			h.push(msg)
		case <-h.shutdown:
			return
		}
	}
}

func (h *Hub) loopDrain() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
		case client := <-h.unregister:
			h.removeClient(client)
		case msg := <-h.broadcast:
			h.push(msg)
		default:
			return
		}
	}
}

func (h *Hub) addClient(client *Client) {
	room := h.roomKey(client.role, client.room)
	h.mu.Lock()
	if _, ok := h.rooms[room]; !ok {
		h.rooms[room] = make(map[*Client]struct{})
	}
	h.rooms[room][client] = struct{}{}
	h.mu.Unlock()
}

func (h *Hub) removeClient(client *Client) {
	room := h.roomKey(client.role, client.room)
	h.mu.Lock()
	if clients, ok := h.rooms[room]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.send)
			if len(clients) == 0 {
				delete(h.rooms, room)
			}
		}
	}
	h.mu.Unlock()
}

func (h *Hub) push(msg Message) {
	room := h.roomKey(msg.Role, msg.RoomID)
	h.mu.RLock()
	clients := h.rooms[room]
	for client := range clients {
		select {
		case client.send <- msg.Payload:
		default:
			go func(c *Client) {
				c.hub.unregister <- c
			}(client)
		}
	}
	h.mu.RUnlock()
}

func (h *Hub) roomKey(role Role, id string) string {
	return string(role) + ":" + id
}

func (h *Hub) count(room string) int {
	sum := 0
	h.mu.RLock()
	for key, clients := range h.rooms {
		if hasRoom(key, room) {
			sum += len(clients)
		}
	}
	h.mu.RUnlock()
	return sum
}

func hasRoom(fullKey, room string) bool {
	parts := strings.SplitN(fullKey, ":", 2)
	if len(parts) != 2 {
		return false
	}
	return parts[1] == room
}
