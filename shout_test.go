package shout_test

import (
	"testing"
	"github.com/andygoblins/shout"
)

var size = 5 //arbitrary

func TestNewShout(t *testing.T) {
	s := shout.New(5)
	if len(s.Send()) != size {
		t.Fatalf("New Shout should have size %d but has size %d\n", size, len(s.Send()))
	}
}

func TestNewListen(t *testing.T) {
	s := shout.New(0)
	l := s.Listen(size)
	if len(l.Rcv()) != size {
		t.Fatalf("New Listen should have size %d but has size %d\n", size, len(l.Rcv()))
	}
}

func TestUnSubscribe(t *testing.T) {
	//TODO unsubscribed channels should be closed, not receive messages, and not cause a deadlock
}

func TestClose(t *testing.T) {
	//TODO Shout.Close() should close all subscribers
}

func TestCloseSend(*testing.T) {
	//TODO should panic
}
