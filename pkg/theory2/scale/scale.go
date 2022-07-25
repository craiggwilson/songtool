package scale

import (
	"github.com/craiggwilson/songtool/pkg/theory2/interval"
)

var (
	Chromatic = Scale{
		intervals: []interval.Interval{
			interval.Perfect(0),
			interval.Minor(1),
			interval.Major(1),
			interval.Minor(2),
			interval.Major(2),
			interval.Perfect(3),
			interval.Diminished(5, 1),
			interval.Perfect(4),
			interval.Minor(5),
			interval.Major(5),
			interval.Minor(6),
			interval.Major(6),
		},
	}
	Ionian = Scale{
		intervals: []interval.Interval{
			interval.Perfect(0),
			interval.Major(1),
			interval.Major(2),
			interval.Perfect(3),
			interval.Perfect(4),
			interval.Major(5),
			interval.Major(6),
		},
	}
)

type Scale struct {
	intervals []interval.Interval
}

func (s Scale) Intervals() []interval.Interval {
	return s.intervals
}
