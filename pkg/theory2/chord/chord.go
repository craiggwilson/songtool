package chord

import (
	"encoding/json"

	"github.com/craiggwilson/songtool/pkg/theory2/interval"
	"github.com/craiggwilson/songtool/pkg/theory2/note"
)

func New(root note.Note, base *note.Note, intervals ...interval.Interval) Chord {
	return Chord{root, intervals, base}
}

type Chord struct {
	root      note.Note
	intervals []interval.Interval
	base      *note.Note
}

func (c Chord) Base() *note.Note {
	return c.base
}

func (c Chord) Intervals() []interval.Interval {
	return c.intervals
}

func (c Chord) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Root      note.Note           `json:"root"`
		Intervals []interval.Interval `json:"intervals"`
		Base      *note.Note          `json:"base"`
	}{c.root, c.intervals, c.base})
}

func (c Chord) Name(namer Namer) string {
	return namer.NameChord(c)
}

func (c Chord) Quality() Quality {
	major3rd := false
	minor3rd := false
	diminished5th := false
	perfect5th := false
	augmented5th := false

	for _, ival := range c.intervals {
		q := ival.Quality()
		switch q.Kind() {
		case interval.QualityKindAugmented:
			augmented5th = augmented5th || (ival.Diatonic() == 4 && q.Size() == 1)
		case interval.QualityKindDiminished:
			diminished5th = diminished5th || (ival.Diatonic() == 4 && q.Size() == 1)
		case interval.QualityKindMajor:
			major3rd = major3rd || ival.Diatonic() == 2
		case interval.QualityKindMinor:
			minor3rd = minor3rd || ival.Diatonic() == 2
		case interval.QualityKindPerfect:
			perfect5th = perfect5th || ival.Diatonic() == 4
		}
	}

	var qualities []Quality
	if major3rd && perfect5th {
		qualities = append(qualities, QualityMajor)
	}
	if major3rd && augmented5th {
		qualities = append(qualities, QualityAugmented)
	}
	if minor3rd && perfect5th {
		qualities = append(qualities, QualityMinor)
	}
	if minor3rd && diminished5th {
		qualities = append(qualities, QualityDiminished)
	}

	quality := QualityIndeterminate
	if len(qualities) == 1 {
		quality = qualities[0]
	}

	return quality
}

func (c Chord) Root() note.Note {
	return c.root
}
