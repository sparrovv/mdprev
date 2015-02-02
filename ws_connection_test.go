package mdprev

import (
	"testing"
	"time"
)

func TestBroadcasting(t *testing.T) {
	exitCh := make(chan bool)
	broadcastCh := make(chan []byte)
	hub := newHub(broadcastCh, exitCh)

	fc_1 := NewFakeConnection()
	fc_2 := NewFakeConnection()
	go fc_1.writer()
	go fc_2.writer()

	go hub.run()
	hub.register <- fc_1
	hub.register <- fc_2

	// test broadcasting
	testMsg := []byte("TEST")
	broadcastCh <- testMsg
	broadcastCh <- testMsg

	time.Sleep(100 * time.Millisecond)
	if len(fc_1.fakeDataBag) != 2 {
		t.Errorf("expected to receive 2 messages, but got %+v", len(fc_1.fakeDataBag))
	}
	if len(fc_2.fakeDataBag) != 2 {
		t.Errorf("expected to receive 2 messages, but got %+v", len(fc_2.fakeDataBag))
	}

	// test that connections are removed
	hub.unregister <- fc_1

	time.Sleep(100 * time.Millisecond)
	if len(hub.connections) != 1 {
		t.Errorf("expected no connections, but got %+v", len(hub.connections))
	}

}

type FakeConnection struct {
	send        chan []byte
	fakeDataBag [][]byte
}

func NewFakeConnection() *FakeConnection {
	return &FakeConnection{
		send: make(chan []byte),
	}
}

func (fc *FakeConnection) writer() {
	for message := range fc.send {
		fc.fakeDataBag = append(fc.fakeDataBag, message)
	}
}

func (fc *FakeConnection) closeCh() {
	close(fc.send)
}

func (fc *FakeConnection) sendMsg(msg []byte) {
	fc.send <- msg
}
