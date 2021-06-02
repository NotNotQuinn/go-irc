package errorStream

import (
	"fmt"

	"github.com/NotNotQuinn/go-irc/channels"
)

func Listen() {
	for {
		// although it doesnt seem like much, it allows for good error loggin later on.
		// Errors should only be passed to this stream if there is no other place, and
		// a panic is not sutible
		err := <-channels.Errors
		fmt.Printf("Error: %+v\n", err)
	}
}
