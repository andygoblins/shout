//package shout provides a broadcast channel that works in select 
//statements.
package shout

//Shout represents a broadcast channel. Each Shout instance has a 
//goroutine that takes messages from the Send channel and transmits 
//them to all subscribed Listen channels.
type Shout struct {
	subscribers map[chan chan interface{}] bool
	sub         chan chan interface{}
	unsub       chan *Listen
	send        chan interface{}
}

//Send returns the broadcast channel. All Listen channels receive 
//messages sent on this channel. Closing this channel causes a panic. 
//Use the Close() method to close down both the Send channel and all 
//Listen channels.
func (*Shout) Send() chan<- interface{} {
	return b.send
}

//run is the Shout event loop. It dies when all Listen channels Close.
func (b *Shout) run() {
	for {
		select {
		case s := <-sub:
			subscribers[s] = true
		case s := <-unsub:
			delete(subscribers, s)
		case m, ok := <-send:
			if !ok {
				panic("Send channel is closed")
			}
			for k := range subscribers {
				k <- m
			}
		}
	}
}

//New creates a Shout with the given buffer size on the Send channel.
func New(n int) *Shout {
	s := Shout{}
	//TODO
}

//Listen returns a new Listen channel with the given buffer size.
func (b *Shout) Listen(n int) *Listen {
	c := make(chan interface{}, n)
	b.sub <- c
	return &Listen{b, c}
}

//Close closes the Shout broadcast channel and all subscriber channels.
func (*Shout) Close() {
	//TODO
}

//Listen is a receiving channel for Shout broadcast messages. 
type Listen struct {
	b   *Shout
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
func (*Listen) Close() {
	//TODO
}
