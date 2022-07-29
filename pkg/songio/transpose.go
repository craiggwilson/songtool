package songio

import (
	"github.com/craiggwilson/songtool/pkg/theory/interval"
)

func Transpose(src Song, interval interval.Interval) *SongTransposer {
	return &SongTransposer{
		src:      src,
		interval: interval,
	}
}

type SongTransposer struct {
	src      Song
	interval interval.Interval
}

func (s *SongTransposer) Next() (Line, bool) {
	nl, ok := s.src.Next()
	if !ok {
		return nl, false
	}

	switch tnl := nl.(type) {
	case *KeyDirectiveLine:
		tnl.Key = tnl.Key.Transpose(s.interval)
	case *ChordLine:
		for _, seg := range tnl.Chords {
			seg.Chord = seg.Chord.Transpose(s.interval)
		}
	}

	return nl, ok
}

func (s *SongTransposer) Err() error {
	return s.src.Err()
}
