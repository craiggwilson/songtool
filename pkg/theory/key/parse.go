package key

import (
	"encoding/json"

	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/craiggwilson/songtool/pkg/theory/note"
)

type Parser interface {
	ParseKey(string) (Parsed, error)
}

type Parsed struct {
	Key
	Suffix string
}

func (k Parsed) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Note   note.Note `json:"note"`
		Kind   Kind      `json:"kind"`
		Suffix string    `json:"suffix"`
	}{k.Note(), k.Kind(), k.Suffix})
}

func (k Parsed) Name(noteNamer note.Namer) string {
	nn := noteNamer.NameNote(k.Note())
	return nn + k.Suffix
}

func (k Parsed) Transpose(by interval.Interval) Parsed {
	return Parsed{
		Key:    k.Key.Transpose(by),
		Suffix: k.Suffix,
	}
}
