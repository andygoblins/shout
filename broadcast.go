//package broadcast provides a broadcast channel for event loops.
package broadcast

//Broadcast represents a broadcast channel. Broadcast channels have a goroutine that takes messages from Send and transmits them to all subscriber channels.
type Broadcast struct {
	subscribers []chan interface{}
	sub         chan chan interface
	unsub       chan int
	//Send is the broadcast channel. All subscribers receive
	//interfaces sent on this channel. Do not close this channel. 
	//Instead, call the Close() method to close down all both the 
	//Send channel and all subscriber channels.
	Send        chan<-interface{}
}

//run is the broadcast event loop. It dies when all subscribers Close.
func (b *broadcast) run() {
	for {
		select {
		case s := <-sub:
			subscribers = append(subscribers, s)
		case i := <-unsub:
			subscribers = append(subscribers[:i], subscribers[i+1:]...)
		case m, ok := <-msg:
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

}

//Subscribe returns a new subscriber channel with the given buffer size.
func (b *Broadcast) Subscribe(n int) <-chan interface{} {
	c := make(chan interface{}, n)
	b.sub <- c
	return c
}

//Close closes the broadcast channel and all subscriber channels.
func (*Broadcast) Close() {

}

//Chan
type Chan struct {
	b *Broadcast
	id int
	Rcv <-chan interface{}
}

//Close unsubscribes from a Broadcast channel. You should always Close an unused Chan, because eventually (depending on the size of the Chan's buffer) the Broadcast will block trying to send a message to it, and no other subscribed channels will receive messages. Alternatively, closing the Broadcast this Chan is subscribed to will close the Chan.
func (*Chan) Close() {

}
