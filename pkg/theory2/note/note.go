package note

import (
	"encoding/json"

	"github.com/craiggwilson/songtool/pkg/theory2/interval"
)

var (
	CFlat  = New(0, 11)
	C      = New(0, 0)
	CSharp = New(0, 1)
	DFlat  = New(1, 1)
	D      = New(1, 2)
	DSharp = New(1, 3)
	EFlat  = New(2, 3)
	E      = New(2, 4)
	ESharp = New(2, 5)
	FFlat  = New(3, 4)
	F      = New(3, 5)
	FSharp = New(3, 6)
	GFlat  = New(4, 6)
	G      = New(4, 7)
	GSharp = New(4, 8)
	AFlat  = New(5, 8)
	A      = New(5, 9)
	ASharp = New(5, 10)
	BFlat  = New(6, 10)
	B      = New(6, 11)
	BSharp = New(6, 0)
)

func New(degreeClass, pitchClass int) Note {
	return Note{normalizeDegreeClass(degreeClass), normalizePitchClass(pitchClass)}
}

type Note struct {
	degreeClass int
	pitchClass  int
}

func (n Note) DegreeClass() int {
	return n.degreeClass
}

func (n Note) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		DegreeClass int `json:"degreeClass"`
		PitchClass  int `json:"pitchClass"`
	}{n.degreeClass, n.pitchClass})
}

func (n Note) PitchClass() int {
	return n.pitchClass
}

func (n Note) Transpose(by interval.Interval) Note {
	current := interval.NewWithChromatic(n.degreeClass, n.pitchClass)
	next := current.Transpose(by)
	return New(next.Diatonic(), next.Chromatic())
}

func normalizeDegreeClass(degreeClass int) int {
	return (degreeClass + 7) % 7
}

func normalizePitchClass(pitchClass int) int {
	return (pitchClass + 12) % 12
}
