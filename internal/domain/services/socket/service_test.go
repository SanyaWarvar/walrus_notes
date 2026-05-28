package socket

import (
	"context"
	"testing"
	"time"
	"wn/internal/domain/dto"
	"wn/pkg/applogger"

	"github.com/google/uuid"
)

type noopLogger struct{}

func (noopLogger) IsDebugLevel() bool                         { return false }
func (noopLogger) IsInfoLevel() bool                          { return false }
func (noopLogger) Debug(string)                               {}
func (noopLogger) Info(string)                                {}
func (noopLogger) Warn(string)                                {}
func (noopLogger) Error(string)                               {}
func (noopLogger) Warnf(string, ...any)                       {}
func (noopLogger) Errorf(string, ...any)                      {}
func (noopLogger) Debugf(string, ...any)                      {}
func (noopLogger) Infof(string, ...any)                       {}
func (l noopLogger) WithCtx(context.Context) applogger.Logger { return l }

type mockConnection struct {
	id       dto.ConnectionID
	userID   uuid.UUID
	sent     []*dto.SocketMessage
	closed   bool
	readMsgs []*dto.SocketMessage
}

func (m *mockConnection) ID() dto.ConnectionID { return m.id }
func (m *mockConnection) UserID() uuid.UUID    { return m.userID }
func (m *mockConnection) Send(msg *dto.SocketMessage) error {
	m.sent = append(m.sent, msg)
	return nil
}
func (m *mockConnection) SendPing() error {
	return m.Send(&dto.SocketMessage{Event: "PING"})
}
func (m *mockConnection) ReadMessage() (*dto.SocketMessage, error) {
	if len(m.readMsgs) == 0 {
		return nil, context.Canceled
	}
	msg := m.readMsgs[0]
	m.readMsgs = m.readMsgs[1:]
	return msg, nil
}
func (m *mockConnection) Close() error {
	m.closed = true
	return nil
}

func TestRegisterHandlerAndHandleMessage(t *testing.T) {
	srv := NewService(noopLogger{})
	userID := uuid.New()
	conn := &mockConnection{id: dto.ConnectionID("c1"), userID: userID}

	srv.RegisterHandler("ECHO", func(msg *dto.SocketMessage, uid uuid.UUID) (*dto.SocketMessage, error) {
		if uid != userID {
			t.Errorf("userID = %v, want %v", uid, userID)
		}
		return msg, nil
	})

	srv.register <- conn
	time.Sleep(20 * time.Millisecond)

	msg := &dto.SocketMessage{Event: "ECHO", Payload: []byte(`{}`)}
	got, err := srv.handleMessage(conn.ID(), msg)
	if err != nil {
		t.Fatalf("handleMessage() error = %v", err)
	}
	if got.Event != "ECHO" {
		t.Errorf("Event = %q, want ECHO", got.Event)
	}
}

func TestHandleMessage_UnknownEvent(t *testing.T) {
	srv := NewService(noopLogger{})
	conn := &mockConnection{id: dto.ConnectionID("c1"), userID: uuid.New()}
	srv.register <- conn
	time.Sleep(20 * time.Millisecond)

	_, err := srv.handleMessage(conn.ID(), &dto.SocketMessage{Event: "UNKNOWN"})
	if err == nil {
		t.Fatal("expected error for unknown event")
	}
}

func TestSendTo(t *testing.T) {
	srv := NewService(noopLogger{})
	conn := &mockConnection{id: dto.ConnectionID("c1"), userID: uuid.New()}
	srv.register <- conn
	time.Sleep(20 * time.Millisecond)

	msg := &dto.SocketMessage{Event: "TEST", Payload: []byte(`{}`)}
	if err := srv.SendTo(conn.ID(), msg); err != nil {
		t.Fatalf("SendTo() error = %v", err)
	}
	if len(conn.sent) != 1 {
		t.Fatalf("sent messages = %d, want 1", len(conn.sent))
	}
}

func TestBroadcast(t *testing.T) {
	srv := NewService(noopLogger{})
	conn1 := &mockConnection{id: dto.ConnectionID("c1"), userID: uuid.New()}
	conn2 := &mockConnection{id: dto.ConnectionID("c2"), userID: uuid.New()}
	srv.register <- conn1
	srv.register <- conn2
	time.Sleep(20 * time.Millisecond)

	srv.Broadcast(&dto.SocketMessage{Event: "BROADCAST", Payload: []byte(`{}`)})
	time.Sleep(20 * time.Millisecond)

	if len(conn1.sent) != 1 || len(conn2.sent) != 1 {
		t.Errorf("broadcast counts: c1=%d c2=%d", len(conn1.sent), len(conn2.sent))
	}
}
