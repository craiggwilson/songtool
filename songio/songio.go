package songio

import "github.com/craiggwilson/songtools/theory"

type SongReader interface {
	NextLine() (Line, bool)
	Err() error
}

func Transpose(cfg *theory.Config, r SongReader, degreeClassInterval int, pitchClassInterval int) *TransposingSongReader {
	return &TransposingSongReader{
		cfg:                 cfg,
		r:                   r,
		degreeClassInterval: degreeClassInterval,
		pitchClassInterval:  pitchClassInterval,
	}
}

type TransposingSongReader struct {
	cfg                 *theory.Config
	r                   SongReader
	degreeClassInterval int
	pitchClassInterval  int
}

func (r *TransposingSongReader) NextLine() (Line, bool) {
	nl, ok := r.r.NextLine()
	if !ok {
		return nl, false
	}

	switch tnl := nl.(type) {
	case *DirectiveLine:
		switch td := tnl.Directive.(type) {
		case *KeyDirective:
			newKey := theory.TransposeKey(r.cfg, td.Key, r.degreeClassInterval, r.pitchClassInterval)
			td.Key = newKey
		}
	case *ChordLine:
		for _, seg := range tnl.Segments {
			seg.Chord = theory.TransposeChord(r.cfg, seg.Chord, r.degreeClassInterval, r.pitchClassInterval)
			// TODO: deal with length
		}
	}

	return nl, ok
}

func (r *TransposingSongReader) Err() error {
	return r.r.Err()
}
