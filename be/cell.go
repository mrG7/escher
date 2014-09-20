// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly.

package be

import (
	// "log"
)

type Cell struct {
	show map[string]*ReCognizer
	see map[string]chan interface{}
	ping chan string
}

func NewCell(r Reflex) *Cell {
	x := &Cell{
		show: make(map[string]*ReCognizer),
		see: make(map[string]chan interface{}),
		ping: make(chan string),
	}
	for vlv, syn := range r {
		v := vlv.(string)
		x.show[v] = syn.Focus(
			func(w interface{}) {
				x.cognize(v, w)
			},
		)
		x.see[v] = make(chan interface{})
	}
	return x
}

// ReCognize
func (x *Cell) ReCognize(valve string, value interface{}) {
	x.show[valve].ReCognize(value)
}

func (x *Cell) cognize(valve string, value interface{}) {
	x.ping <- valve
	x.see[valve] <- value
}

func (x *Cell) Cognize() (valve string, value interface{}) {
	valve = <- x.ping
	return valve, <-x.see[valve]
}