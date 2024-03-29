package note

import (
	"encoding/json"
	"sync"

	"github.com/craiggwilson/songtool/pkg/theory/interval"
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

	degreeClassToPitchClass = [7]int{0, 2, 4, 5, 7, 9, 11}
	notes                   []Note
	initOnce                sync.Once
)

func List() []Note {
	initOnce.Do(func() {
		notes = make([]Note, 0, 21)
		for i := 0; i < 7; i++ {
			pc := degreeClassToPitchClass[i]
			notes = append(notes, New(i, pc-1))
			notes = append(notes, New(i, pc))
			notes = append(notes, New(i, pc+1))
		}
	})

	localNotes := make([]Note, len(notes))
	copy(localNotes, notes)
	return localNotes
}

func New(degreeClass, pitchClass int) Note {
	return Note{normalizeDegreeClass(degreeClass), normalizePitchClass(pitchClass)}
}

type Note struct {
	degreeClass int
	pitchClass  int
}

func (n Note) Accidentals() int {
	stdPitchClass := degreeClassToPitchClass[n.degreeClass]
	return n.pitchClass - stdPitchClass
}

func (n Note) CompareTo(o Note) int {
	if n.degreeClass < o.degreeClass {
		return -1
	}
	if n.degreeClass > o.degreeClass {
		return 1
	}

	if n.pitchClass < o.pitchClass {
		return -1
	}
	if n.pitchClass > o.pitchClass {
		return 1
	}

	return 0
}

func (n Note) DegreeClass() int {
	return n.degreeClass
}

func (n Note) Enharmonic() interval.Interval {
	accidentals := n.Accidentals()
	switch {
	case accidentals > 0:
		return interval.NewWithChromatic(1, 0)
	case accidentals < 0:
		return interval.NewWithChromatic(-1, 0)
	default:
		return interval.NewWithChromatic(0, 0)
	}
}

func (n Note) Interval(other Note) interval.Interval {
	return interval.NewWithChromatic(other.degreeClass-n.degreeClass, other.pitchClass-n.pitchClass)
}

func (n Note) Step(step int) interval.Interval {
	switch {
	case step < 0:
		return stepDown(n.degreeClass, n.pitchClass, step)
	case step > 0:
		return stepUp(n.degreeClass, n.pitchClass, step)
	default:
		return interval.NewWithChromatic(0, 0)
	}
}

func (n Note) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		DegreeClass int `json:"degreeClass"`
		PitchClass  int `json:"pitchClass"`
	}{n.degreeClass, n.pitchClass})
}

func (n Note) Name(namer Namer) string {
	return namer.NameNote(n)
}

func (n Note) PitchClass() int {
	return n.pitchClass
}

func (n Note) Transpose(by interval.Interval) Note {
	current := interval.NewWithChromatic(n.degreeClass, n.pitchClass)
	next := current.Transpose(by)
	return New(next.Diatonic(), next.Chromatic())
}

func stepDown(degreeClass, pitchClass, step int) interval.Interval {
	diatonicSteps := 0
	chromaticSteps := 0
	for step < 0 {
		stdPitchClass := degreeClassToPitchClass[normalizeDegreeClass(degreeClass+diatonicSteps)]
		accidentals := pitchClass - stdPitchClass + chromaticSteps
		prevStdPitchClass := degreeClassToPitchClass[normalizeDegreeClass(degreeClass-1+diatonicSteps)]
		normalizedNextPitchClass := normalizePitchClass(stdPitchClass + accidentals - 1)
		if prevStdPitchClass-normalizedNextPitchClass > 6 {
			prevStdPitchClass -= 12
		}
		if normalizedNextPitchClass <= prevStdPitchClass {
			diatonicSteps--
		}
		chromaticSteps--
		step++
	}

	return interval.NewWithChromatic(diatonicSteps, chromaticSteps)
}

func stepUp(degreeClass, pitchClass, step int) interval.Interval {
	diatonicSteps := 0
	chromaticSteps := 0
	for step > 0 {
		stdPitchClass := degreeClassToPitchClass[normalizeDegreeClass(degreeClass+diatonicSteps)]
		accidentals := pitchClass - stdPitchClass + chromaticSteps
		nextStdPitchClass := degreeClassToPitchClass[normalizeDegreeClass(degreeClass+diatonicSteps+1)]
		normalizedPitchClass := normalizePitchClass(stdPitchClass + accidentals + 1)
		if normalizedPitchClass-nextStdPitchClass > 6 {
			nextStdPitchClass += 12
		}
		if normalizedPitchClass >= nextStdPitchClass {
			diatonicSteps++
		}
		chromaticSteps++
		step--
	}

	return interval.NewWithChromatic(diatonicSteps, chromaticSteps)
}

func normalizeDegreeClass(degreeClass int) int {
	return (degreeClass + 7) % 7
}

func normalizePitchClass(pitchClass int) int {
	return (pitchClass + 12) % 12
}
