package signaling

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	conn   *websocket.Conn
	roomID string
	userID string
	server *Server
}

type Room struct {
	clients map[*Client]bool
	mu      sync.Mutex
}

type Server struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		rooms: make(map[string]*Room),
	}
}

func (s *Server) addClient(roomID string, client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, exists := s.rooms[roomID]
	if !exists {
		room = &Room{
			clients: make(map[*Client]bool),
		}
		s.rooms[roomID] = room
	}

	room.mu.Lock()
	room.clients[client] = true
	room.mu.Unlock()
}

func (s *Server) removeClient(roomID string, client *Client) {
	s.mu.RLock()
	room, exists := s.rooms[roomID]
	s.mu.RUnlock()

	if !exists {
		return
	}

	room.mu.Lock()
	delete(room.clients, client)
	room.mu.Unlock()

	// Cleanup empty room
	if len(room.clients) == 0 {
		s.mu.Lock()
		delete(s.rooms, roomID)
		s.mu.Unlock()
	}
}

func (s *Server) broadcast(roomID string, message []byte, sender *Client) {
	s.mu.RLock()
	room, exists := s.rooms[roomID]
	s.mu.RUnlock()

	if !exists {
		return
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	for client := range room.clients {
		if client != sender {
			if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.server.removeClient(c.roomID, c)
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		c.server.broadcast(c.roomID, message, c)
	}
}

func NewClient(conn *websocket.Conn, roomID, userID string, server *Server) *Client {
	client := &Client{
		conn:   conn,
		roomID: roomID,
		userID: userID,
		server: server,
	}

	server.addClient(roomID, client)
	go client.readPump()

	return client
}