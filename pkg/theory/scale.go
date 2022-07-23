package theory

type Scale struct {
	Name  string
	Notes []Note
}

func GenerateScale(cfg *Config, name string, root Note, intervals []int) Scale {
	scale := Scale{
		Name:  name,
		Notes: make([]Note, 1, len(intervals)+1),
	}

	scale.Notes[0] = root
	prevNote := root
	for _, interval := range intervals {
		nextNote := TransposeNote(cfg, prevNote, IntervalFromStep(cfg, prevNote, interval, EnharmonicSharp))
		scale.Notes = append(scale.Notes, nextNote)
		prevNote = nextNote
	}

	return scale
}

// type ChordScaleKind string

// const (
// 	ChordScaleMajor ChordScaleKind = "major"
// )

// type ChordScale struct {
// 	Name   string
// 	Chords []Chord
// }

// func GenerateMajorChordScale(cfg *Config, root Note) ChordScale {
// 	if cfg == nil {
// 		cfg = &defaultConfig
// 	}

// }
