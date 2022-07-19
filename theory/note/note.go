package note

type DegreeClass int
type PitchClass int

type Note struct {
	Name        string
	DegreeClass DegreeClass
	PitchClass  PitchClass
	Accidentals int
}
