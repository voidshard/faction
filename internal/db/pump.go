package db

import (
	"log"
	"sync"
	"time"

	"github.com/voidshard/faction/pkg/structs"
)

const (
	chunksize = 1000
)

// Pump is a simple struct that implements a buffer + timer to flush to the DB as needed
// as new data comes in.
type Pump struct {
	db   *FactionDB
	size int

	family chan *structs.Family
	people chan *structs.Person
	tuples map[Relation]chan *structs.Tuple

	tuplesShared chan *rTuple

	bufFamily []*structs.Family
	bufPeople []*structs.Person
	bufTuples map[Relation][]*structs.Tuple

	// Routine control
	// order routine to emit rows to db
	//  - emit True; routine should emit, then exit
	//  - emit False; routine should emit, not exit
	emit chan bool
	done sync.WaitGroup
	errs chan error
}

type rTuple struct {
	Relation Relation
	Tuple    *structs.Tuple
}

func newPump(db *FactionDB) *Pump {
	tupShared := make(chan *rTuple)

	tup := map[Relation]chan *structs.Tuple{}
	tupBuf := map[Relation][]*structs.Tuple{}
	for _, r := range allRelations {
		tupBuf[r] = []*structs.Tuple{}
		tchan := make(chan *structs.Tuple)
		tup[r] = tchan

		go func(rel Relation, c <-chan *structs.Tuple) {
			// forward incoming tuples to shared tuple chan
			for t := range c {
				tupShared <- &rTuple{rel, t}
			}
		}(r, tchan)
	}

	p := &Pump{
		db:           db,
		size:         chunksize,
		family:       make(chan *structs.Family),
		people:       make(chan *structs.Person),
		tuples:       tup,
		tuplesShared: tupShared,
		bufFamily:    []*structs.Family{},
		bufPeople:    []*structs.Person{},
		bufTuples:    tupBuf,
		emit:         make(chan bool),
		done:         sync.WaitGroup{},
		errs:         make(chan error),
	}
	p.done.Add(1)
	go p.doPump()

	return p
}

// doPump is the main pump routine.
func (p *Pump) doPump() {
	defer p.done.Done()
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C: // try to flush buffer every few seconds
			p.errs <- p.write()
		case exit := <-p.emit: // we've been explicitly told to flush
			err := p.write()
			if err != nil {
				p.errs <- err
			}
			if exit {
				return
			}
		case i := <-p.family:
			if i == nil {
				continue
			}
			p.bufFamily = append(p.bufFamily, i)
			if len(p.bufFamily) >= p.size {
				p.errs <- p.write()
			}
		case i := <-p.people:
			if i == nil {
				continue
			}
			p.bufPeople = append(p.bufPeople, i)
			if len(p.bufPeople) >= p.size {
				p.errs <- p.write()
			}
		case i := <-p.tuplesShared:
			if i == nil {
				continue
			}
			buf, ok := p.bufTuples[i.Relation]
			if !ok {
				continue
			}
			p.bufTuples[i.Relation] = append(buf, i.Tuple)
			if len(buf)+1 >= p.size {
				p.errs <- p.write()
			}
		}
	}
}

// write will write the current buffer(s) to the database.
// Only used by pump routine
func (p *Pump) write() error {
	// check all of our buffers, and write to db if they're full
	return p.db.InTransaction(func(tx ReaderWriter) error {
		if len(p.bufFamily) >= p.size {
			log.Printf("Flushing %d families", len(p.bufFamily))
			err := tx.SetFamilies(p.bufFamily...)
			if err != nil {
				return err
			}
			p.bufFamily = []*structs.Family{}
		}
		if len(p.bufPeople) >= p.size {
			log.Printf("Flushing %d people", len(p.bufPeople))
			err := tx.SetPeople(p.bufPeople...)
			if err != nil {
				return err
			}
			p.bufPeople = []*structs.Person{}
		}
		for r, buf := range p.bufTuples {
			if len(buf) >= p.size {
				log.Printf("Flushing %d tuples for %s", len(buf), r)
				err := tx.SetTuples(r, buf...)
				if err != nil {
					return err
				}
				p.bufTuples[r] = []*structs.Tuple{}
			}
		}
		return nil
	})
}

// Errors returns a channel of errors that the pump has encountered.
// This must be checked regularly, or the pump will block.
func (p *Pump) Errors() <-chan error {
	return p.errs
}

// Close stops internal channels, flushes data to the database.
func (p *Pump) Close() {
	close(p.family)
	close(p.people)
	for _, c := range p.tuples {
		close(c)
	}
	p.size = 1 // explicitly set size to 1 so we flush everything
	p.emit <- true
	p.done.Wait() // wait for pump routine to exit
	close(p.errs)
}

// SetFamilies will queue a family for writing to the database
func (p *Pump) SetFamilies(in ...*structs.Family) {
	for _, i := range in {
		p.family <- i
	}
}

// SetPeople will queue a person for writing to the database
func (p *Pump) SetPeople(in ...*structs.Person) {
	for _, i := range in {
		p.people <- i
	}
}

// SetTuples will queue a tuple for writing to the database
func (p *Pump) SetTuples(r Relation, in ...*structs.Tuple) {
	chnl, ok := p.tuples[r]
	if !ok {
		return
	}
	for _, i := range in {
		chnl <- i
	}
}
