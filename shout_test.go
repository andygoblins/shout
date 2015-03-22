package shout

import "testing"

func TestBroadcast(*testing.T) {
	//all subscribed channels should receive messages
}

func TestUnSubscribe(*testing.T) {
	//unsubscribed channels should be closed, not receive messages, and not cause a deadlock
}

func TestClose(*testing.T) {
	//Shout.Close() should close all subscribers
}

func TestCloseSend(*testing.T) {
	//should panic
}
