package multyplayer

import (
	"testing"
)

func TestConnectAndGetRoomParticipants(t *testing.T) {
	srv := NewService()
	conn := &Connection{Send: make(chan []byte, 1)}

	srv.Connect("user-1", "note-1", conn)

	participants := srv.GetRoomParticipants("note-1")
	if len(participants) != 1 || participants[0] != "user-1" {
		t.Errorf("participants = %v, want [user-1]", participants)
	}
	if !srv.RoomExists("note-1") {
		t.Error("room should exist")
	}
}

func TestDisconnect_RemovesUserAndRoom(t *testing.T) {
	srv := NewService()
	conn := &Connection{Send: make(chan []byte, 1)}

	srv.Connect("user-1", "note-1", conn)
	srv.Disconnect("user-1", "note-1")

	if srv.RoomExists("note-1") {
		t.Error("room should be removed after last user disconnects")
	}
	if len(srv.GetRoomParticipants("note-1")) != 0 {
		t.Error("participants should be empty")
	}
}

func TestHandleMessage_BroadcastsToRoom(t *testing.T) {
	srv := NewService()
	conn1 := &Connection{Send: make(chan []byte, 1)}
	conn2 := &Connection{Send: make(chan []byte, 1)}

	srv.Connect("user-1", "note-1", conn1)
	srv.Connect("user-2", "note-1", conn2)

	msg := []byte("hello")
	got := srv.HandleMessage("user-1", "note-1", msg)
	if string(got) != "hello" {
		t.Errorf("returned message = %q, want %q", got, msg)
	}

	select {
	case received := <-conn1.Send:
		if string(received) != "hello" {
			t.Errorf("conn1 received = %q, want hello", received)
		}
	default:
		t.Error("conn1 should receive broadcast message")
	}

	select {
	case received := <-conn2.Send:
		if string(received) != "hello" {
			t.Errorf("conn2 received = %q, want hello", received)
		}
	default:
		t.Error("conn2 should receive broadcast message")
	}
}

func TestHandleMessage_UnknownRoom(t *testing.T) {
	srv := NewService()
	msg := []byte("ping")
	got := srv.HandleMessage("user-1", "missing", msg)
	if string(got) != "ping" {
		t.Errorf("returned message = %q, want ping", got)
	}
}

func TestHandleMessage_RemovesDeadConnection(t *testing.T) {
	srv := NewService()
	closed := false
	deadConn := &Connection{
		Send: make(chan []byte),
		Close: func() error {
			closed = true
			return nil
		},
	}
	aliveConn := &Connection{Send: make(chan []byte, 1)}

	srv.Connect("dead-user", "note-1", deadConn)
	srv.Connect("alive-user", "note-1", aliveConn)

	srv.HandleMessage("alive-user", "note-1", []byte("msg"))

	if !closed {
		t.Error("dead connection Close should be called")
	}
	participants := srv.GetRoomParticipants("note-1")
	if len(participants) != 1 || participants[0] != "alive-user" {
		t.Errorf("participants = %v, want [alive-user]", participants)
	}
}
