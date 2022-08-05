package key

import (
	"encoding/json"
	"sort"
	"strings"
	"sync"

	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/craiggwilson/songtool/pkg/theory/note"
)

type Kind string

const (
	KindMajor Kind = "Major"
	KindMinor Kind = "Minor"
)

var (
	keys     []Key
	initOnce sync.Once

	degreeClassToPitchClass = [7]int{0, 2, 4, 5, 7, 9, 11}
)

func List() []Key {
	initOnce.Do(func() {
		notes := note.List()
		keys = make([]Key, 0, len(notes)*2)
		for i := 0; i < len(notes); i++ {
			keys = append(keys, Major(notes[i]))
			keys = append(keys, Minor(notes[i]))
		}
	})

	localKeys := make([]Key, len(keys))
	copy(localKeys, keys)
	return localKeys
}

func Major(n note.Note) Key {
	return New(n, KindMajor)
}

func Minor(n note.Note) Key {
	return New(n, KindMinor)
}

func New(n note.Note, k Kind) Key {
	return Key{n, k}
}

func Sort(keys []Key) {
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].CompareTo(keys[j]) < 0
	})
}

type Key struct {
	note note.Note
	kind Kind
}

func (k Key) CompareTo(o Key) int {
	comp := k.note.CompareTo(o.note)
	if comp != 0 {
		return comp
	}

	return strings.Compare(string(k.kind), string(o.kind))
}

func (k Key) Enharmonic() interval.Interval {
	return k.note.Enharmonic()
}

func (k Key) Kind() Kind {
	return k.kind
}

func (k Key) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Note note.Note `json:"note"`
		Kind Kind      `json:"kind"`
	}{k.note, k.kind})
}

func (k Key) Name(namer Namer) string {
	return namer.NameKey(k)
}

func (k Key) Note() note.Note {
	return k.note
}

func (k Key) Step(step int) interval.Interval {
	return k.note.Step(step)
}

func (k Key) Transpose(by interval.Interval) Key {
	newNote := k.note.Transpose(by)

	return New(
		newNote,
		k.kind,
	)
}
