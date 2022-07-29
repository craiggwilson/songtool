package songio

import "github.com/craiggwilson/songtool/pkg/theory/chord"

type Line interface {
	line()
}

type EmptyLine struct{}

func (EmptyLine) line() {}

type ChordLine struct {
	Chords []*ChordOffset `json:"chords"`
}

func (l *ChordLine) line() {}

type ChordOffset struct {
	Chord  chord.Named `json:"chord"`
	Offset int         `json:"offset"`
}

type TextLine struct {
	Text string `json:"text"`
}

func (l *TextLine) line() {}
