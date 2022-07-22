package theory

type Interval struct {
	DegreeClass int
	PitchClass  int
}

func IntervalFromStep(cfg *Config, note Note, step int, enharmonic Enharmonic) Interval {
	if cfg == nil {
		cfg = &defaultConfig
	}

	newDegreeClass := degreeClassFromPitchClass(cfg, adjustPitchClass(cfg, note.PitchClass, step), enharmonic)

	return Interval{
		DegreeClass: int(newDegreeClass) - int(note.DegreeClass),
		PitchClass:  step,
	}
}

func IntervalFromDiff(a, b Note) Interval {
	return Interval{
		DegreeClass: int(b.DegreeClass) - int(a.DegreeClass),
		PitchClass:  int(b.PitchClass) - int(a.PitchClass),
	}
}
