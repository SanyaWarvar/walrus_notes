package multyplayer

import (
	"sync"
)

// Service handles WebSocket connections and message proxying
type Service struct {
	// rooms maps noteId to a map of userId to connection
	rooms map[string]map[string]*Connection
	mu    sync.RWMutex
}

// Connection represents a WebSocket connection
type Connection struct {
	// Send is a channel for sending messages to the client
	Send chan []byte

	// Close is a function to close the connection
	Close func() error
}

// NewService creates a new instance of Service
func NewService() *Service {
	return &Service{
		rooms: make(map[string]map[string]*Connection),
	}
}

// Connect handles new WebSocket connection
// userId - identifier of the user
// noteId - room identifier
// conn - WebSocket connection with send channel and close function
func (s *Service) Connect(userId, noteId string, conn *Connection) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create room if not exists
	if _, exists := s.rooms[noteId]; !exists {
		s.rooms[noteId] = make(map[string]*Connection)
	}

	// Store connection
	s.rooms[noteId][userId] = conn
}

// Disconnect removes user from room
func (s *Service) Disconnect(userId, noteId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if room, exists := s.rooms[noteId]; exists {
		delete(room, userId)

		// Clean up empty rooms
		if len(room) == 0 {
			delete(s.rooms, noteId)
		}
	}
}

// HandleMessage processes incoming message and broadcasts to room participants
// userId - sender identifier
// noteId - room identifier
// message - raw bytes received from WebSocket
// Returns the same message that was broadcasted
func (s *Service) HandleMessage(userId, noteId string, message []byte) []byte {
	s.mu.RLock()
	room, exists := s.rooms[noteId]
	s.mu.RUnlock()

	if !exists {
		return message
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Broadcast to all participants in the room
	for id, conn := range room {
		// Skip sender if needed (uncomment line below)
		// if id == userId { continue }

		select {
		case conn.Send <- message:
			// Message sent successfully
		default:
			// Connection is dead, remove it
			if conn.Close != nil {
				conn.Close()
			}
			delete(room, id)
		}
	}

	// Clean up empty rooms
	if len(room) == 0 {
		delete(s.rooms, noteId)
	}

	return message
}

// GetRoomParticipants returns list of user IDs in a room
func (s *Service) GetRoomParticipants(noteId string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, exists := s.rooms[noteId]
	if !exists {
		return []string{}
	}

	participants := make([]string, 0, len(room))
	for id := range room {
		participants = append(participants, id)
	}

	return participants
}

// RoomExists checks if room exists and has participants
func (s *Service) RoomExists(noteId string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, exists := s.rooms[noteId]
	return exists && len(room) > 0
}
