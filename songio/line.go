package songio

import "github.com/craiggwilson/songtools/theory"

type Line interface {
	line()
}

type EmptyLine struct{}

func (EmptyLine) line() {}

type ChordLine struct {
	Chords []ChordSegment
}

func (l *ChordLine) line() {}

type ChordSegment struct {
	theory.Chord
	Offset int
}

type TextLine struct {
	Text string
}

func (l *TextLine) line() {}
