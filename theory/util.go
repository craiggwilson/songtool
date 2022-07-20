package theory

import (
	"fmt"
	"math"
)

func adjustDegreeClass(cfg *Config, degreeClass DegreeClass, by int) DegreeClass {
	return (degreeClass + DegreeClass(by) + DegreeClass(len(cfg.NaturalNoteNames))) % DegreeClass(len(cfg.NaturalNoteNames))
}

func adjustPitchClass(cfg *Config, pitchClass PitchClass, by int) PitchClass {
	return (pitchClass + PitchClass(by) + PitchClass(cfg.PitchClassCount)) % PitchClass(cfg.PitchClassCount)
}

func classDelta(a, b, count int) int {
	d1 := b - a
	d2 := b - a + count
	if math.Abs(float64(d1)) < math.Abs(float64(d2)) {
		return d1
	}

	return d2
}

func degreeClassDelta(cfg *Config, a, b DegreeClass) int {
	return classDelta(int(a), int(b), len(cfg.NaturalNoteNames))
}

func degreeClassFromNaturalNoteName(cfg *Config, naturalNoteName rune) DegreeClass {
	for i, nn := range cfg.NaturalNoteNames {
		if nn == naturalNoteName {
			return DegreeClass(i)
		}
	}

	panic(fmt.Sprintf("natural note name %q does not map to a degree class", naturalNoteName))
}

func degreeClassFromPitchClass(cfg *Config, pitchClass PitchClass, enharmonic Enharmonic) DegreeClass {
	switch enharmonic {
	case EnharmonicSharp:
		for i := len(cfg.DegreeClassToPitchClass) - 1; i >= 0; i-- {
			if pitchClass >= cfg.DegreeClassToPitchClass[i] {
				return DegreeClass(i)
			}
		}
	case EnharmonicFlat:
		for i := 0; i < len(cfg.DegreeClassToPitchClass); i++ {
			if pitchClass <= cfg.DegreeClassToPitchClass[i] {
				return DegreeClass(i)
			}
		}
	default:
		panic(fmt.Sprintf("invalid enharmonic %d", enharmonic))
	}

	panic(fmt.Sprintf("invalid pitch class %d", pitchClass))
}

func normalize(v int, count int) int {
	switch {
	case v > count/2:
		return -count + v
	case v < -count/2:
		return count + v
	default:
		return v
	}
}

func normalizeAccidentals(cfg *Config, accidentals int) int {
	return normalize(accidentals, cfg.PitchClassCount)
}

func pitchClassFromDegreeClass(cfg *Config, degreeClass DegreeClass) PitchClass {
	if int(degreeClass) < len(cfg.DegreeClassToPitchClass) {
		return cfg.DegreeClassToPitchClass[int(degreeClass)]
	}

	panic(fmt.Sprintf("degree class %d does not map to a pitch class", degreeClass))
}

func pitchClassDelta(cfg *Config, a, b PitchClass) int {
	return classDelta(int(a), int(b), cfg.PitchClassCount)
}
