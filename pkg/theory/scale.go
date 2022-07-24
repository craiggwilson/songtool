package theory

type Scale struct {
	Name  string
	Notes []Note
}

func GenerateDiatonicScale(name string, root Note, intervals []Interval) Scale {
	return std.GenerateDiatonicScale(name, root, intervals)
}

func (t *Theory) GenerateDiatonicScale(name string, root Note, intervals []Interval) Scale {
	scale := Scale{
		Name:  name,
		Notes: make([]Note, len(intervals)),
	}

	for i, interval := range intervals {
		scale.Notes[i] = t.TransposeNote(root, interval)
	}

	return scale
}
