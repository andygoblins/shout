//package shout provides broadcast channels that work in select
//statements.
package shout

import "sync"

//Shout represents a broadcast channel.
//Each Shout instance has a goroutine that takes messages from the
//Send channel and transmits them to all subscribed Listen channels.
type Shout struct {
	msub        sync.Mutex //subscribers mutex
	subscribers map[chan interface{}]bool
	send        chan interface{}
	done        chan struct{}
}

//Send returns the broadcast channel.
//All Listen channels receive messages sent on this channel. Closing
//this channel causes a panic. Use the Close() method to close down
//both the Send channel and all Listen channels.
func (b *Shout) Send() chan<- interface{} {
	return b.send
}

//run is the Shout event loop.
//It returns when it receives a message on s.done.
func (s *Shout) run() {
	for {
		select {
		case msg, ok := <-s.send:
			if !ok {
				panic("Send channel is closed")
			}
			s.msub.Lock()
			for key := range s.subscribers {
				key <- msg
			}
			s.msub.Unlock()
		case <-s.done:
			return
		}
	}
}

//New creates a Shout with the given buffer size for the Send channel.
//New starts the Shout goroutine that takes messages from Send and
//transmits them to all subscribed Listens.
func New(n int) *Shout {
	s := Shout{}
	s.subscribers = make(map[chan interface{}]bool)
	s.send = make(chan interface{}, n)
	s.done = make(chan struct{})
	go s.run()
	return &s
}

//Listen returns a new Listen channel with the given buffer size.
func (s *Shout) Listen(n int) *Listen {
	s.msub.Lock()
	defer s.msub.Unlock()
	c := make(chan interface{}, n)
	s.subscribers[c] = true
	return &Listen{s, c}
}

//Close closes the Shout broadcast channel and all subscriber channels.
func (s *Shout) Close() {
	s.msub.Lock()
	defer s.msub.Unlock()
	for k := range s.subscribers {
		close(k)
	}
	//Tell run() to return. Can't do close(s.done) because the
	//close(s.send) might be processed by run() first, causing a panic.
	s.done <- struct{}{}
	close(s.send)
}

//Listen is a receiving channel for Shout broadcast messages.
type Listen struct {
	s   *Shout
	rcv chan interface{}
}

//Rcv returns the receiving channel for Shout broadcast messages.
func (c *Listen) Rcv() <-chan interface{} {
	return c.rcv
}

//Close unsubscribes Listen from a Shout channel.
//If an unused Listen is not Closed, the Shout will eventually
//(depends on size of Rcv buffer) block trying to send a message to
//it, and no other subscribed Listen channels will receive messages.
//Closing a Shout will also close all subscribed Listens.
func (c *Listen) Close() {
	c.s.msub.Lock()
	defer c.s.msub.Unlock()
	delete(c.s.subscribers, c.rcv)
	close(c.rcv)
}
