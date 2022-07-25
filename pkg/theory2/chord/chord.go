package chord

import (
	"github.com/craiggwilson/songtool/pkg/theory2/interval"
)

func New(intervals ...interval.Interval) Chord {
	return Chord{intervals}
}

type Chord struct {
	intervals []interval.Interval
}

func (c Chord) Quality() Quality {
	major3rd := false
	minor3rd := false
	diminished5th := false
	perfect5th := false
	augmented5th := false

	for _, ival := range c.intervals {
		q := ival.Quality()
		switch q.Kind() {
		case interval.QualityKindAugmented:
			augmented5th = augmented5th || (ival.Diatonic() == 4 && q.Size() == 1)
		case interval.QualityKindDiminished:
			diminished5th = diminished5th || (ival.Diatonic() == 4 && q.Size() == 1)
		case interval.QualityKindMajor:
			major3rd = major3rd || ival.Diatonic() == 2
		case interval.QualityKindMinor:
			minor3rd = minor3rd || ival.Diatonic() == 2
		case interval.QualityKindPerfect:
			perfect5th = perfect5th || ival.Diatonic() == 4
		}
	}

	var qualities []Quality
	if major3rd && perfect5th {
		qualities = append(qualities, MajorQuality)
	}
	if major3rd && augmented5th {
		qualities = append(qualities, AugmentedQuality)
	}
	if minor3rd && perfect5th {
		qualities = append(qualities, MinorQuality)
	}
	if minor3rd && diminished5th {
		qualities = append(qualities, DiminishedQuality)
	}

	quality := IndeterminateQuality
	if len(qualities) == 1 {
		quality = qualities[0]
	}

	return quality
}

func (c Chord) Intervals() []interval.Interval {
	return c.intervals
}
