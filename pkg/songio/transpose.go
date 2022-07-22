package songio

import (
	"github.com/craiggwilson/songtool/pkg/theory"
)

func Transpose(cfg *theory.Config, src Song, interval theory.Interval) *SongTransposer {
	return &SongTransposer{
		cfg:      cfg,
		src:      src,
		interval: interval,
	}
}

type SongTransposer struct {
	cfg      *theory.Config
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
		newKey := theory.TransposeKey(s.cfg, tnl.Key, s.interval)
		tnl.Key = newKey
	case *ChordLine:
		for _, seg := range tnl.Chords {
			seg.Chord = theory.TransposeChord(s.cfg, seg.Chord, s.interval)
		}
	}

	return nl, ok
}

func (s *SongTransposer) Err() error {
	return s.src.Err()
}
