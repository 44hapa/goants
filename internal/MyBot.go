package internal

import (
	"io"
	"log"
	"math/rand"
	"os"
)

type MyBot struct {
}

// NewBot creates a new instance of your bot
func NewBot(s *State) Bot {
	mb := &MyBot{
		//do any necessary initialization here
	}
	return mb
}

// DoTurn is where you should do your bot's actual work.
func (mb *MyBot) DoTurn(s *State) error {
	f, err := os.OpenFile("/tmp/orders.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	//defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	//wrt := io.Writer(f)
	log.SetOutput(wrt)
	log.Println(s.Map.Ants)

	dirs := []Direction{North, East, South, West}
	antNum := 0
	for loc, ant := range s.Map.Ants {
		antNum++
		// Нашли чужого муравья
		if ant != MY_ANT {
			continue
		}

		//try each direction in a random order
		p := []int{}
		if antNum%2 == 0 {
			p = rand.Perm(2)
		} else {
			p = rand.Perm(4)
		}

		for _, i := range p {
			d := dirs[i]

			loc2 := s.Map.Move(loc, d)
			if s.Map.SafeDestination(loc2) {
				s.IssueOrderLoc(loc, d)
				//there's also an s.IssueOrderRowCol if you don't have a Location handy
				break
			}
		}
	}
	//returning an error will halt the whole program!
	return nil
}
