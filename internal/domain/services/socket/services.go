package socket

import (
	"context"
	"fmt"
	"sync"
	"time"
	"wn/internal/domain/dto"
	"wn/pkg/applogger"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ConnectionID string

type Connection interface {
	ID() ConnectionID
	UserID() uuid.UUID
	Send(msg *dto.SocketMessage) error
	SendPing() error
	Close() error
	ReadMessage() (*dto.SocketMessage, error)
}

type MessageHandler func(msg *dto.SocketMessage, userId uuid.UUID) (*dto.SocketMessage, error)

type Service struct {
	lgr         applogger.Logger
	connections sync.Map // map[ConnectionID]Connection
	handlers    map[string]MessageHandler
	broadcast   chan *dto.SocketMessage
	register    chan Connection
	unregister  chan Connection
	mu          sync.RWMutex
}

func NewService(lgr applogger.Logger) *Service {
	s := &Service{
		lgr:        lgr,
		handlers:   map[string]MessageHandler{},
		broadcast:  make(chan *dto.SocketMessage, 100),
		register:   make(chan Connection, 10),
		unregister: make(chan Connection, 10),
	}

	go s.run()
	return s
}

// Запуск основного цикла обработки
func (s *Service) run() {
	for {
		select {
		case conn := <-s.register:
			s.connections.Store(conn.ID(), conn)
			s.lgr.Infof("connection registered: %s", conn.ID())

		case conn := <-s.unregister:
			s.connections.Delete(conn.ID())
			conn.Close()
			s.lgr.Infof("connection unregistered: %s", conn.ID())

		case msg := <-s.broadcast:
			s.broadcastMessage(msg)
		}
	}
}

// HTTP хендлер для апгрейда соединения
func (s *Service) HandleConnection(ctx context.Context, conn Connection) {
	// Регистрируем соединение
	s.register <- conn
	defer func() {
		s.unregister <- conn
	}()

	// Создаем тикер для пинга каждые 10 секунд
	pingTicker := time.NewTicker(10 * time.Second)
	defer pingTicker.Stop()

	// Канал для обработки сообщений
	messageChan := make(chan *dto.SocketMessage)
	errorChan := make(chan error)

	// Запускаем горутину для чтения сообщений
	go func() {
		for {
			msg, err := conn.ReadMessage()
			if err != nil {
				errorChan <- err
				return
			}
			messageChan <- msg
		}
	}()

	// Основной цикл обработки
	for {
		select {
		case <-ctx.Done():
			// Контекст отменен
			return

		case <-pingTicker.C:
			// Отправляем пинг каждые 10 секунд
			if err := conn.SendPing(); err != nil {
				s.lgr.Errorf("send ping error: %s", err.Error())
				return
			}
			s.lgr.Debugf("ping sent to connection %s", conn.ID())

		case msg := <-messageChan:
			// Обрабатываем входящее сообщение
			processedMsg, err := s.handleMessage(conn.ID(), msg)
			if err != nil {
				s.lgr.Errorf("handle message error: %s", err)
			}
			if processedMsg != nil {
				if err := conn.Send(processedMsg); err != nil {
					s.lgr.Errorf("send message error: %s", err.Error())
					return
				}
			}

		case err := <-errorChan:
			// Ошибка чтения сообщения
			s.lgr.Errorf("read message error: %s", err.Error())
			return
		}
	}
}

// Обработка входящих сообщений
func (s *Service) handleMessage(connID ConnectionID, msg *dto.SocketMessage) (*dto.SocketMessage, error) {
	s.mu.RLock()
	handler, exists := s.handlers[msg.Event]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no handler for event: %s", msg.Event)
	}
	connAny, _ := s.connections.Load(connID)
	conn := connAny.(Connection)
	return handler(msg, conn.UserID())
}

// Регистрация обработчиков сообщений
func (s *Service) RegisterHandler(event string, handler MessageHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[event] = handler
}

// Бродкаст всем соединениям
func (s *Service) Broadcast(msg *dto.SocketMessage) {
	s.broadcast <- msg
}

// Отправка конкретному соединению
func (s *Service) SendTo(connID ConnectionID, msg *dto.SocketMessage) error {
	if conn, ok := s.connections.Load(connID); ok {
		return conn.(Connection).Send(msg)
	}
	return fmt.Errorf("connection not found: %s", connID)
}

// Бродкаст сообщения
func (s *Service) broadcastMessage(msg *dto.SocketMessage) {
	s.connections.Range(func(key, value interface{}) bool {
		if conn, ok := value.(Connection); ok {
			if err := conn.Send(msg); err != nil {
				s.lgr.Errorf("broadcast error: connId: %s error: %s", key.(string), err.Error())
			}
		}
		return true
	})
}

type WSConnection struct {
	conn   *websocket.Conn
	id     ConnectionID
	userID uuid.UUID
	mu     sync.Mutex
}

func NewWSConnection(conn *websocket.Conn, userID uuid.UUID) *WSConnection {
	return &WSConnection{
		conn:   conn,
		id:     ConnectionID(uuid.New().String()),
		userID: userID,
	}
}

func (w *WSConnection) ID() ConnectionID {
	return w.id
}

func (w *WSConnection) UserID() uuid.UUID {
	return w.UserID()
}

func (w *WSConnection) Send(msg *dto.SocketMessage) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.WriteJSON(msg)
}

func (w *WSConnection) SendPing() error {
	return w.Send(&dto.SocketMessage{
		Event:   "PING",
		Payload: []byte{},
	})
}

func (w *WSConnection) ReadMessage() (*dto.SocketMessage, error) {
	var msg dto.SocketMessage
	err := w.conn.ReadJSON(&msg)
	return &msg, err
}

func (w *WSConnection) Close() error {
	return w.conn.Close()
}
