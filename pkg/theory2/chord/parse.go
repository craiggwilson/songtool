package chord

import (
	"encoding/json"

	"github.com/craiggwilson/songtool/pkg/theory2/interval"
	"github.com/craiggwilson/songtool/pkg/theory2/note"
)

type Parser interface {
	ParseChord(string) (Parsed, error)
}

type Parsed struct {
	Chord

	Suffix            string
	BaseNoteDelimiter string
}

func (p Parsed) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Root              note.Note           `json:"root"`
		Suffix            string              `json:"suffix"`
		BaseNoteDelimiter string              `json:"baseNoteDelimiter,omitempty"`
		Base              *note.Note          `json:"base"`
		Intervals         []interval.Interval `json:"intervals"`
	}{p.Chord.root, p.Suffix, p.BaseNoteDelimiter, p.Chord.base, p.Chord.intervals})
}

func (p Parsed) Name(noteNamer note.Namer) string {
	name := noteNamer.NameNote(p.Chord.root)
	name += p.Suffix
	if p.base != nil {
		name += p.BaseNoteDelimiter
		name += noteNamer.NameNote(*p.Chord.base)
	}

	return name
}
