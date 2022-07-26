package key

import (
	"encoding/json"
	"sort"
	"strings"
	"sync"

	"github.com/craiggwilson/songtool/pkg/theory2/note"
)

type Kind string

const (
	Major Kind = "Major"
	Minor Kind = "Minor"
)

var (
	C           Key = New(note.C, Major)
	CMinor      Key = New(note.C, Minor)
	CSharp      Key = New(note.CSharp, Major)
	CSharpMinor Key = New(note.CSharp, Minor)
	DFlat       Key = New(note.DFlat, Major)
	DFlatMinor  Key = New(note.DFlat, Minor)
	D           Key = New(note.D, Major)
	DMinor      Key = New(note.D, Minor)
	DSharp      Key = New(note.DSharp, Major)
	DSharpMinor Key = New(note.DSharp, Minor)
	EFlat       Key = New(note.EFlat, Major)
	EFlatMinor  Key = New(note.EFlat, Minor)
	E           Key = New(note.E, Major)
	EMinor      Key = New(note.E, Minor)
	ESharp      Key = New(note.ESharp, Major)
	ESharpMinor Key = New(note.ESharp, Minor)
	FFlat       Key = New(note.FFlat, Major)
	FFlatMinor  Key = New(note.FFlat, Minor)
	F           Key = New(note.F, Major)
	FMinor      Key = New(note.F, Minor)
	FSharp      Key = New(note.FSharp, Major)
	FSharpMinor Key = New(note.FSharp, Minor)
	GFlat       Key = New(note.GFlat, Major)
	GFlatMinor  Key = New(note.GFlat, Minor)
	G           Key = New(note.G, Major)
	GMinor      Key = New(note.G, Minor)
	GSharp      Key = New(note.GSharp, Major)
	GSharpMinor Key = New(note.GSharp, Minor)
	AFlat       Key = New(note.AFlat, Major)
	AFlatMinor  Key = New(note.AFlat, Minor)
	A           Key = New(note.A, Major)
	AMinor      Key = New(note.A, Minor)
	ASharp      Key = New(note.ASharp, Major)
	ASharpMinor Key = New(note.ASharp, Minor)
	BFlat       Key = New(note.BFlat, Major)
	BFlatMinor  Key = New(note.BFlat, Minor)
	B           Key = New(note.B, Major)
	BMinor      Key = New(note.B, Minor)
	BSharp      Key = New(note.BSharp, Major)
	BSharpMinor Key = New(note.BSharp, Minor)

	keys     []Key
	initOnce sync.Once
)

func List() []Key {
	initOnce.Do(func() {
		notes := note.List()
		keys = make([]Key, 0, len(notes)*2)
		for i := 0; i < len(notes); i++ {
			keys = append(keys, New(
				notes[i],
				Major,
			))
			keys = append(keys, New(
				notes[i],
				Minor,
			))
		}
	})

	localKeys := make([]Key, len(keys))
	copy(localKeys, keys)
	return localKeys
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

func (k Key) Kind() Kind {
	return k.kind
}

func (k Key) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Note note.Note `json:"note"`
		Kind Kind      `json:"kind"`
	}{k.note, k.kind})
}

func (k Key) Note() note.Note {
	return k.note
}
