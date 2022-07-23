package theory

type Interval struct {
	DegreeClass int
	PitchClass  int
}

func IntervalFromStep(note Note, step int, enharmonic Enharmonic) Interval {
	return std.IntervalFromStep(note, step, enharmonic)
}

func (t *Theory) IntervalFromStep(note Note, step int, enharmonic Enharmonic) Interval {
	newDegreeClass := t.Config.DegreeClassFromPitchClass(t.Config.AdjustPitchClass(note.PitchClass, step), enharmonic)

	return Interval{
		DegreeClass: int(newDegreeClass) - int(note.DegreeClass),
		PitchClass:  step,
	}
}

func IntervalFromDiff(a, b Note) Interval {
	return std.IntervalFromDiff(a, b)
}

func (t *Theory) IntervalFromDiff(a, b Note) Interval {
	return Interval{
		DegreeClass: int(b.DegreeClass) - int(a.DegreeClass),
		PitchClass:  int(b.PitchClass) - int(a.PitchClass),
	}
}
