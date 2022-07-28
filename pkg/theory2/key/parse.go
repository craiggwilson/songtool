package key

import (
	"encoding/json"

	"github.com/craiggwilson/songtool/pkg/theory2/note"
)

type KeyParser interface {
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
