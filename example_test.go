package shout_test

import (
	"fmt"
	"github.com/andygoblins/shout"
	"sync"
)

func ExampleShout() {
	var listeners = 3
	var wg sync.WaitGroup

	s := shout.New(1)
	//Start listeners
	for i := 0; i < listeners; i++ {
		l := s.Listen(0)
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("Received %s from s\n", <-l.Rcv())
		}()
	}
	//Send message to all listeners.
	s.Send() <- "broadcast"
	//Wait for listeners to receive message.
	wg.Wait()
	//Output:
	//Received broadcast from s
	//Received broadcast from s
	//Received broadcast from s
}
