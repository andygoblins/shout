//package broadcast provides a broadcast channel for event loops.
package broadcast

type broadcast struct {
	subscribers []*bchan
	sub         chan *bchan
	unsub       chan int
	bufsize     int
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

type bchan struct {
	send chan interface{}
	recv chan interface{}
}

func (b *bchan) close() {
	close(b.send)
	close(b.recv)
}

type Chan struct {
	Send chan<- interface{}
	Recv <-chan interface{}
	id   int
	bc   *broadcast
}

func New(n int) *Chan {
	//TODO init broadcast
	go func() {
		defer close(bc.sub)
		defer close(bc.unsub)
		bc.run()
	}()
}

func (c *Chan) Close() {
	//TODO should I call close here, or in run() ?
	c.bc.subscribers[c.id].close()
	c.bc.unsub <- c.id
}

//Is having a Chan beget more Chans make sense? Is this better than
//having a public Broadcast type whose only purpose is to beget Chans?
func (c *Chan) Subscribe() *Chan {
	s := make(chan<- interface{}, c.bf.bufsize)
	r := make(chan<- interface{}, c.bf.bufsize)
	c.bc.sub <- bchan{s, r}
	return &Chan{s, r, len(c.bc.subscribers) - 1, bc}
}
