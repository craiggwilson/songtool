package theory2

import (
	"github.com/craiggwilson/songtool/pkg/theory2/note"
	"github.com/craiggwilson/songtool/pkg/theory2/scale"
)

var std = func() *Theory {
	cfg := DefaultConfig()

	return New(cfg, cfg, cfg)
}()

func Default() *Theory {
	return std
}

func New(noteNamer NoteNamer, noteParser NoteParser, scaleParser ScaleParser) *Theory {
	return &Theory{
		noteNamer:   noteNamer,
		noteParser:  noteParser,
		scaleParser: scaleParser,
	}
}

type Theory struct {
	noteNamer   NoteNamer
	noteParser  NoteParser
	scaleParser ScaleParser
}

func (t *Theory) NameNote(n note.Note) string {
	return t.noteNamer.NameNote(n)
}

func (t *Theory) ParseNote(text string) (note.Note, error) {
	return t.noteParser.ParseNote(text)
}

func (t *Theory) ParseScale(text string) (scale.Scale, error) {
	return t.scaleParser.ParseScale(text)
}
