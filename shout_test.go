package shout_test

import (
	"github.com/andygoblins/shout"
	"testing"
)

var size = 5 //arbitrary

func TestNewShout(t *testing.T) {
	s := shout.New(5)
	defer s.Close()
	if len(s.Send()) != size {
		t.Fatalf("New Shout should have size %d but has size %d\n", size, len(s.Send()))
	}
}

func TestNewListen(t *testing.T) {
	s := shout.New(0)
	defer s.Close()
	l := s.Listen(size)
	if len(l.Rcv()) != size {
		t.Fatalf("New Listen should have size %d but has size %d\n", size, len(l.Rcv()))
	}
}

func TestUnSubscribe(t *testing.T) {
	s := shout.New(1)
	defer s.Close()
	l := make([]*shout.Listen, size)
	for i := range l {
		l[i] = s.Listen(1)
	}

	l[0].Close()
	s.Send() <- "msg"
	if _, ok := <-l[0].Rcv(); ok {
		t.Error("l[0] isn't closed")
	}
	for i, v := range l[1:] {
		if _, ok := <-v.Rcv(); !ok {
			t.Errorf("l[%d] is closed", i)
		}
	}
}

func TestClose(t *testing.T) {
	s := shout.New(1)
	l := make([]*shout.Listen, size)
	for i := range l {
		l[i] = s.Listen(1)
	}

	s.Close()
	for i, v := range l {
		if _, ok := <-v.Rcv(); ok {
			t.Errorf("l[%d] isn't closed", i)
		}
	}
}

func TestCloseSend(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("Shout didn't panic when Send() was closed.")
		}
	}()
	s := shout.New(1)
	close(s.Send())
}
