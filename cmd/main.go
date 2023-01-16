package main

import (
	"ants/internal"
	"io"
	"log"
)

// main initializes the state and starts the processing loop
func main() {
	var s internal.State
	errStart := s.Start()
	if errStart != nil {
		log.Panicf("Start() failed (%s)", errStart)
	}
	mb := internal.NewBot(&s)
	errStart = s.Loop(mb, func() {
		//if you want to do other between-turn debugging things, you can do them here
	})
	if errStart != nil && errStart != io.EOF {
		log.Panicf("Loop() failed (%s)", errStart)
	}
}
