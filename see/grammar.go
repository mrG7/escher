// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly.

//	nand {
//		a and
//		n not
//		a.XandY = n.X
//		n.notX // the empty-string valve of the nand circuit connects to n's notX valve
//		a.X = X
//		a.Y = 1
//	}
//
// Has the following star representation, after seeing (parsing):
//
//	nand {
//		peer {
//			{
//				X Name("X")
//				Name("") //  the empty string valve
//			}
//			a Name("and")
//			n Name("not")
//		}
//		matching {
//			0 {
//				Left {
//					Peer Name("a")
//					Valve Name("XandY")
//				}
//				Right {
//					Peer Name("n")
//					Valve Name("X")
//				}
//			} // end of 0
//			…
//		}
//	}
//
package see

import (
	"fmt"
	"strconv"

	"github.com/gocircuit/escher/star"
)

// Design is one of the built-in designs listed below.
type Design interface{
	String() string
}

type (
	Name string
	RootName string
	String string
	Int int
	Float float64
	Complex complex128
)

func (x Star) String() string {
	return (*star.Star)(x).String()
}

func (x Complex) String() string {
	return fmt.Sprintf("Complex(%g)", x)
}

func (x Float) String() string {
	return fmt.Sprintf("Float(%g)", x)
}

func (x Int) String() string {
	return fmt.Sprintf("Int(%d)", x)
}

func (x Name) String() string {
	return fmt.Sprintf("Name(%d)", x)
}

func (x RootName) String() string {
	return fmt.Sprintf("RootName(%d)", x)
}

func (x String) String() string {
	return fmt.Sprintf("String(%s)", strconv.Quote(chop(string(x))))
}

func chop(x string) string {
	if len(x) < 21 {
		return x
	}
	return x[:20]+"…"
}

// Reflex represents the name of a reflex and its arguments.
type Reflex struct {
	Name  string
	Valve []string
}

type Circuit struct {
	Name string
	Valve []string
	Peer    []*Peer
	Match []*Matching
}

type Peer struct {
	Name   string
	Design Design
}

func (p *Peer) String() string {
	return fmt.Sprintf("Peer{%s, %v}", p.Name, p.Design)
}

type Matching struct {
	Join [2]Join
}

// Join is one of PeerJoin, ValveJoin or DesignJoin.
type Join interface{}

// E.g. “a.X”
type PeerJoin struct {
	Peer string
	Valve string
}

// E.g. “Y”
type ValveJoin struct {
	Valve string
}

// E.g. “12.1e3”
type DesignJoin struct {
	Design Design
}
