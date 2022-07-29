package key

import (
	"encoding/json"

	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/craiggwilson/songtool/pkg/theory/note"
)

type Namer interface {
	NameKey(Key) string
}

type Named struct {
	Parsed

	Name string
}

func (n Named) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name   string    `json:"name"`
		Note   note.Note `json:"note"`
		Kind   Kind      `json:"kind"`
		Suffix string    `json:"suffix"`
	}{n.Name, n.Note(), n.Kind(), n.Suffix})
}

func (n Named) Transpose(noteNamer note.Namer, by interval.Interval) Named {
	parsed := Parsed{
		Key:    n.Key.Transpose(by),
		Suffix: n.Suffix,
	}

	return Named{
		Parsed: parsed,
		Name:   parsed.Name(noteNamer),
	}
}
