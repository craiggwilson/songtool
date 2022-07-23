package songio

import (
	"github.com/craiggwilson/songtool/pkg/theory"
)

func Transpose(t *theory.Theory, src Song, interval theory.Interval) *SongTransposer {
	return &SongTransposer{
		t:        t,
		src:      src,
		interval: interval,
	}
}

type SongTransposer struct {
	t        *theory.Theory
	src      Song
	interval theory.Interval
}

func (s *SongTransposer) Next() (Line, bool) {
	nl, ok := s.src.Next()
	if !ok {
		return nl, false
	}

	switch tnl := nl.(type) {
	case *KeyDirectiveLine:
		newKey := s.t.TransposeKey(tnl.Key, s.interval)
		tnl.Key = newKey
	case *ChordLine:
		for _, seg := range tnl.Chords {
			seg.Chord = s.t.TransposeChord(seg.Chord, s.interval)
		}
	}

	return nl, ok
}

func (s *SongTransposer) Err() error {
	return s.src.Err()
}
