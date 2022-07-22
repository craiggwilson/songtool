package songio

import (
	"github.com/craiggwilson/songtool/pkg/theory"
)

func Transpose(cfg *theory.Config, src LineIter, interval theory.Interval) *TransposingLineIter {
	return &TransposingLineIter{
		cfg:      cfg,
		src:      src,
		interval: interval,
	}
}

type TransposingLineIter struct {
	cfg      *theory.Config
	src      LineIter
	interval theory.Interval
}

func (it *TransposingLineIter) Next() (Line, bool) {
	nl, ok := it.src.Next()
	if !ok {
		return nl, false
	}

	switch tnl := nl.(type) {
	case *KeyDirectiveLine:
		newKey := theory.TransposeKey(it.cfg, tnl.Key, it.interval)
		tnl.Key = newKey
	case *ChordLine:
		for _, seg := range tnl.Chords {
			seg.Chord = theory.TransposeChord(it.cfg, seg.Chord, it.interval)
		}
	}

	return nl, ok
}

func (r *TransposingLineIter) Err() error {
	return r.src.Err()
}
