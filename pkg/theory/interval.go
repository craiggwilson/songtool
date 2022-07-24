package theory

import "fmt"

var IonianIntervals = []Interval{
	{0, 0},
	{1, 2},
	{2, 4},
	{3, 5},
	{4, 7},
	{5, 9},
	{6, 11},
}

var ChromaticIntervals = []Interval{
	{0, 0},
	{1, 1},
	{1, 2},
	{2, 3},
	{2, 4},
	{3, 5},
	{3, 6},
	{4, 6},
	{5, 8},
	{5, 9},
	{6, 10},
	{6, 11},
}

type IntervalStr string

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
		DegreeClass: int(t.Config.AdjustDegreeClass(b.DegreeClass, -int(a.DegreeClass))),
		PitchClass:  int(t.Config.AdjustPitchClass(b.PitchClass, -int(a.PitchClass))),
	}
}

func MustInterval(interval Interval, err error) Interval {
	if err != nil {
		panic(err)
	}

	return interval
}

func ParseInterval(interval string) (Interval, error) {
	return std.ParseInterval(interval)
}

func (t *Theory) ParseInterval(interval string) (Interval, error) {
	switch interval {
	case "1P":
		return Interval{0, 0}, nil
	case "2m":
		return Interval{1, 1}, nil
	case "2M":
		return Interval{1, 2}, nil
	case "3m":
		return Interval{2, 3}, nil
	case "3M":
		return Interval{2, 4}, nil
	case "4P":
		return Interval{3, 5}, nil
	case "4a":
		return Interval{3, 6}, nil
	case "5d":
		return Interval{4, 6}, nil
	case "5P":
		return Interval{4, 7}, nil
	case "6m":
		return Interval{5, 8}, nil
	case "6M":
		return Interval{5, 9}, nil
	case "7m":
		return Interval{6, 10}, nil
	case "7M":
		return Interval{6, 11}, nil
	default:
		return Interval{}, fmt.Errorf("unknown interval %q", interval)
	}
}
