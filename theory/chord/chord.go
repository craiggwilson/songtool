package chord

import (
	"github.com/craiggwilson/songtools/theory/note"
)

type Chord struct {
	Root      note.Note
	Intervals []int
	Suffix    string
	BaseNote  note.Note
}
