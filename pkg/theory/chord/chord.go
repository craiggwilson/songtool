package chord

import (
	"github.com/craiggwilson/songtool/pkg/theory/interval"
)

func New(intervals ...interval.Interval) Chord {
	return Chord{intervals, determineQuality(intervals)}
}

type Chord struct {
	intervals []interval.Interval
	quality   Quality
}

func (c Chord) Quality() Quality {
	return c.quality
}

func (c Chord) Intervals() []interval.Interval {
	return c.intervals
}

func determineQuality(intervals []interval.Interval) Quality {
	major3rd := false
	minor3rd := false
	diminished5th := false
	perfect5th := false
	augmented5th := false

	for _, interval := range intervals {
		switch interval.Chromatic() {
		case 3:
			minor3rd = true
		case 4:
			major3rd = true
		case 5:
			diminished5th = true
		case 7:
			perfect5th = true
		case 8:
			augmented5th = true
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
