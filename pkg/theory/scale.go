package theory

type Scale struct {
	Name  string
	Notes []Note
}

func GenerateScale(name string, root Note, intervals []int) Scale {
	return defaultTheory.GenerateScale(name, root, intervals)
}

func (t *Theory) GenerateScale(name string, root Note, intervals []int) Scale {
	scale := Scale{
		Name:  name,
		Notes: make([]Note, 1, len(intervals)+1),
	}

	scale.Notes[0] = root
	prevNote := root
	for _, interval := range intervals {
		nextNote := t.TransposeNote(prevNote, t.IntervalFromStep(prevNote, interval, Sharp))
		scale.Notes = append(scale.Notes, nextNote)
		prevNote = nextNote
	}

	return scale
}
