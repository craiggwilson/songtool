package note

import "github.com/craiggwilson/songtool/pkg/theory/interval"

const ()

func New(diatonic, chromatic int) Note {
	return Note{interval.New(diatonic, chromatic)}
}

type Note struct {
	interval interval.Interval
}

func (n Note) Diatonic() int {
	return n.interval.Diatonic()
}

func (n Note) Chromatic() int {
	return n.interval.Chromatic()
}

func (n Note) Interval() interval.Interval {
	return n.interval
}

func (n Note) Transpose(interval interval.Interval) Note {
	return Note{n.interval.Transpose(interval)}
}
