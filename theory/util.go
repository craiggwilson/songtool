package theory

import "fmt"

const (
	degreeClassCount = 7
	pitchClassCount  = 12
)

var degreeClassToPitchClass = []PitchClass{0, 2, 4, 5, 7, 9, 11}

func AdjustPitchClass(_ *Config, pitchClass PitchClass, by int) PitchClass {
	return (pitchClass + PitchClass(by) + pitchClassCount) % pitchClassCount
}

func DegreeClassFromNaturalNoteName(cfg *Config, naturalNoteName rune) DegreeClass {
	for i, nn := range cfg.NaturalNoteNames {
		if nn == naturalNoteName {
			return DegreeClass(i)
		}
	}

	panic(fmt.Sprintf("natural note name %q does not map to a degree class", naturalNoteName))
}

func PitchClassFromDegreeClass(cfg *Config, degreeClass DegreeClass) PitchClass {
	if int(degreeClass) < len(degreeClassToPitchClass) {
		return degreeClassToPitchClass[int(degreeClass)]
	}

	panic(fmt.Sprintf("degree class %d does not map to a pitch class", degreeClass))
}
