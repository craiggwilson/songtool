package theory2

import (
	"github.com/craiggwilson/songtool/pkg/theory2/interval"
	"github.com/craiggwilson/songtool/pkg/theory2/note"
)

func GenerateScale(name string, root note.Note, intervals ...interval.Interval) Scale {
	scale := Scale{
		name:  name,
		notes: make([]note.Note, len(intervals)),
	}

	for i, interval := range intervals {
		scale.notes[i] = root.Transpose(interval)
	}

	return scale
}

func NewScale(name string, notes ...note.Note) Scale {
	return Scale{name, notes}
}

type Scale struct {
	name  string
	notes []note.Note
}

func (s Scale) Name() string {
	return s.name
}

func (s Scale) Notes() []note.Note {
	return s.notes
}
