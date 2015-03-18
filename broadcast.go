//package broadcast provides a broadcast channel for event loops.
package broadcast

//Broadcast represents a broadcast channel. Broadcast channels have a goroutine that takes messages from Send and transmits them to all subscriber channels.
type Broadcast struct {
	subscribers []chan interface{}
	sub         chan chan interface{}
	unsub       chan int
	send        chan interface{}
}

//Send returns the broadcast channel. All subscribers receive messages
//sent on this channel. Do not close this channel. Instead, call the
//Close() method to close down both the Send channel and all subscriber
//channels.
func (*Broadcast) Send() chan<- interface{} {
	return b.send
}

//run is the broadcast event loop. It dies when all subscribers Close.
func (b *Broadcast) run() {
	for {
		select {
		case s := <-sub:
			subscribers = append(subscribers, s)
		case i := <-unsub:
			subscribers = append(subscribers[:i], subscribers[i+1:]...)
		case m, ok := <-send:
			if !ok {
				panic("Send channel is closed")
			}
			for _, v := range subscribers {
				v <- m
			}
		}
	}
}

//New creates a Broadcast with the given buffer size on the Send channel.
func New(n int) *Broadcast {
	//TODO
}

//Subscribe returns a new subscriber channel with the given buffer size.
func (b *Broadcast) Subscribe(n int) <-chan interface{} {
	c := make(chan interface{}, n)
	b.sub <- c
	return c
}

//Close closes the broadcast channel and all subscriber channels.
func (*Broadcast) Close() {
	//TODO
}

//Chan
type Chan struct {
	b   *Broadcast
	id  int
	rcv chan interface{}
}

//Rcv returns the receiving channel for broadcast messages.
func (c *Chan) Rcv() chan<- interface{} {
	return c.rcv
}

//Close unsubscribes from a Broadcast channel. You should always Close an unused Chan, because eventually (depending on the size of the Chan's buffer) the Broadcast will block trying to send a message to it, and no other subscribed channels will receive messages. Alternatively, closing the Broadcast this Chan is subscribed to will close the Chan.
func (*Chan) Close() {
	//TODO
}
