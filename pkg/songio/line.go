package songio

import "github.com/craiggwilson/songtool/pkg/theory"

type Line interface {
	line()
}

type EmptyLine struct{}

func (EmptyLine) line() {}

type ChordLine struct {
	Chords []*ChordOffset
}

func (l *ChordLine) line() {}

type ChordOffset struct {
	theory.Chord
	Offset int
}

type TextLine struct {
	Text string
}

func (l *TextLine) line() {}
