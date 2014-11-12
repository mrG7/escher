// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly.

package be

import (
	"fmt"

	. "github.com/gocircuit/escher/circuit"
)

func Materialize(index Index, design Value) (residue Value) {
	var reflex Reflex
	reflex, residue = MaterializeReflex(index, design)
	if len(reflex) > 0 {
		panic("circuit not closed")
	}
	return
}

func MaterializeReflex(index Index, design Value) (reflex Reflex, residue Value) {
	renderer := newRenderer(index)
	matter := &Matter{
		Index: index,
		Design: design,
		View: New(),
		Path: []Name{},
		Super: nil,
	}
	return renderer.Materialize(matter, design, true)
}

type renderer struct {
	index Index
}

func newRenderer(index Index) *renderer {
	return &renderer{index}
}

func (b *renderer) lookup(addr Address) Value {
	return b.index.Recall(addr.Path...)
}

func (b *renderer) expandAddress(matter *Matter, addr Address) (Reflex, Value) {
	addr, monkey := filterMonkey(addr) // parse @-sign out from front of address

	// first, looking up addr within the circuit that encloses this address reference
	var val Value
	if matter != nil && matter.Super != nil {
		enclosing := matter.Super.Address
		if len(enclosing.Path) > 0 {
			abs := Address{enclosing.Path[:len(enclosing.Path)-1]}
			abs = abs.Append(addr)
			val = b.lookup(abs)
			if val != nil {
				addr = abs
			}
		}
	}

	matter.Address = addr

	// if not found locally, find the addr starting from root
	if val == nil {
		val = b.lookup(addr)
	}
	if val == nil {
		panicf("Address %v is dangling in matter:\n%v\n", addr, matter)
	}
	if monkey {
		return MaterializeNoun(matter, val)
	}

	// fill in address and source value and return
	matter.Design = val

	return b.Materialize(matter, val, true)
}

// expand is false, if and only if it is invoked by MaterialiazeCircuit.
func (b *renderer) Materialize(matter *Matter, x Value, expand bool) (Reflex, Value) {

	switch t := x.(type) {
	// Addresses are materialized recursively
	case Address:
		return b.expandAddress(matter, t)
	// Primitive types are materialized as gates that emit their values once (these gates are called nouns)
	case int, float64, complex128, string:
		return MaterializeNoun(matter, t)
	case Materializer:
		return t.Materialize(matter)
	case Circuit:
		if expand {
			return b.materializeCircuit(matter, t)
		}
		return MaterializeNoun(matter, t)
	default:
		panicf("Source address %v points to unknown type %T", matter.Address, x)
	}
	panic(0)
}

var SpiritAddress = NewAddress("escher", "Spirit")

func (b *renderer) materializeCircuit(matter *Matter, u Circuit) (Reflex, Value) {
	residue := New()
	gates := make(map[Name]Reflex)
	spirit := make(map[Name]interface{})

	// iterate and materialize gates
	for g, _ := range u.Gate {
		if g == Super {
			panicf("Circuit design overwrites the “%s” gate. In:\n%v\n", Super, u)
		}
		m := u.At(g)
		var gv Value
		if Same(m, SpiritAddress) {
			gates[g], gv, spirit[g] = MaterializeNativeInstance(
				&Matter{
					Index: b.index,
					Address: Address{},
					Design: m,
					View: u.View(g),
					Path: append(matter.Path, g),
					Super: matter,
				},
				&Future{},
			)
		} else {
			gates[g], gv = b.Materialize(
				&Matter{
					Index: b.index,
					Address: Address{},
					Design: m,
					View: u.View(g),
					Path: append(matter.Path, g),
					Super: matter,
				},
				m,
				false,
			)
		}
		residue.Gate[g] = gv
	}

	// compute the super reflex to be returned by this circuit's materialization
	var super Reflex
	super, gates[Super] = make(Reflex), make(Reflex)
	for v, _ := range u.Valves(Super) {
		super[v], gates[Super][v] = NewSynapse()
	}

	// residue.Gate[Genus] = matter.Circuit()

	// link up all gates
	for _, g_ := range append(u.Names(), Super) {
		g := g_
		for v_, t := range u.Valves(g) {
			v := v_
			checkLink(u, gates, g, v, t.Gate, t.Valve)
			residue.Link(Vector{g, v}, Vector{t.Gate, t.Valve})
			go Link(gates[g][v], gates[t.Gate][t.Valve])
			// go func() {
			//	log.Printf("%v:%v -> %v:%v | %v %v", g, v, t.Gate, t.Valve, gates[g][v], gates[t.Gate][t.Valve])
			// 	Link(gates[g][v], gates[t.Gate][t.Valve])
			// }()
		}
	}

	// send residue of this circuit to all escher.Spirit reflexes
	res := CleanUp(residue)
	go func() {
		for _, f := range spirit {
			f.(*Future).Charge(res)
		}
	}()

	if !res.IsEmpty() {
		res.Grow("?", "Index")
	}
	return super, res
}

func checkLink(u Circuit, gates map[Name]Reflex, sg, sv, tg, tv Name) {
	// log.Printf(" %v:%v <=> %v:%v", sg, sv, tg, tv)
	if _, ok := gates[sg]; !ok {
		panicf("In circuit:\n%v\nHas no gate %v\n",u,  sg)
	}
	if _, ok := gates[tg]; !ok {
		panicf("In circuit:\n%v\nHas no gate %v\n",u,  tg)
	}
	if _, ok := gates[sg][sv]; !ok {
		panicf("In circuit:\n%v\nGate %v has no valve :%v\n",u,  sg, sv)
	}
	if _, ok := gates[tg][tv]; !ok {
		panicf("In circuit:\n%v\nGate %v has no valve :%v\n",u,  tg, tv)
	}
}

func panicf(f string, a ...interface{}) {
	panic(fmt.Sprintf(f, a...))
}