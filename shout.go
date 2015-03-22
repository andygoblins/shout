//package shout provides a broadcast channel that works in select
//statements.
package shout

//Shout represents a broadcast channel. Each Shout instance has a
//goroutine that takes messages from the Send channel and transmits
//them to all subscribed Listen channels.
type Shout struct {
	subscribers map[chan interface{}]bool
	sub         chan chan interface{}
	unsub       chan chan interface{}
	send        chan interface{}
}

//Send returns the broadcast channel. All Listen channels receive
//messages sent on this channel. Closing this channel causes a panic.
//Use the Close() method to close down both the Send channel and all
//Listen channels.
func (b *Shout) Send() chan<- interface{} {
	return b.send
}

//run is the Shout event loop. It dies when unsub is closed.
func (s *Shout) run() {
	for {
		select {
		case sub := <-s.sub:
			s.subscribers[sub] = true
		case unsub, ok := <-s.unsub:
			if !ok {
				return
			}
			delete(s.subscribers, unsub)
		case msg, ok := <-s.send:
			if !ok {
				panic("Send channel is closed")
			}
			for key := range s.subscribers {
				key <- msg
			}
		}
	}
}

//New creates a Shout with the given buffer size on the Send channel.
func New(n int) *Shout {
	s := Shout{}
	s.subscribers = make(map[chan interface{}]bool)
	s.sub = make(chan chan interface{})
	s.unsub = make(chan chan interface{})
	s.send = make(chan interface{}, n)
	go s.run()
	return &s
}

//Listen returns a new Listen channel with the given buffer size.
func (s *Shout) Listen(n int) *Listen {
	c := make(chan interface{}, n)
	s.sub <- c
	return &Listen{s, c}
}

//Close closes the Shout broadcast channel and all subscriber channels.
func (s *Shout) Close() {
	for k := range s.subscribers {
		close(k)
	}
	close(s.unsub) //this causes run() to return
	close(s.sub)
	close(s.send)
}

//Listen is a receiving channel for Shout broadcast messages.
type Listen struct {
	s   *Shout
	rcv chan interface{}
}

//Rcv returns the receiving channel for Shout broadcast messages.
func (c *Listen) Rcv() chan<- interface{} {
	return c.rcv
}

//Close unsubscribes Listen from a Shout channel. You should always
//Close an unused Listen, because eventually Shout will block trying to
//send a message to it, and no other subscribed Listen channels will
//receive messages. Alternatively, closing the Shout this Listen is
//subscribed to will close the Listen.
func (c *Listen) Close() {
	c.s.unsub <- c.rcv
	close(c.rcv)
}
