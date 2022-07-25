package theory2

import (
	"github.com/craiggwilson/songtool/pkg/theory2/chord"
	"github.com/craiggwilson/songtool/pkg/theory2/note"
)

type Chord struct {
	root   note.Note
	suffix string

	base note.Note

	chord chord.Chord
}
