package chord

import (
	"encoding/json"

	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/craiggwilson/songtool/pkg/theory/note"
)

type Namer interface {
	NameChord(Chord) string
}

type Named struct {
	Parsed

	Name string
}

func (n Named) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name              string              `json:"name"`
		Root              note.Note           `json:"root"`
		Suffix            string              `json:"suffix"`
		BaseNoteDelimiter string              `json:"baseNoteDelimiter,omitempty"`
		Base              *note.Note          `json:"base"`
		Intervals         []interval.Interval `json:"intervals"`
	}{n.Name, n.Chord.root, n.Suffix, n.BaseNoteDelimiter, n.Chord.base, n.Chord.intervals})
}

func (n Named) Transpose(noteNamer note.Namer, by interval.Interval) Named {
	parsed := Parsed{
		Chord:             n.Chord.Transpose(by),
		Suffix:            n.Suffix,
		BaseNoteDelimiter: n.BaseNoteDelimiter,
	}

	return Named{
		Parsed: parsed,
		Name:   parsed.Name(noteNamer),
	}
}
