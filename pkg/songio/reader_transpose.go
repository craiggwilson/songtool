package songio

import "github.com/craiggwilson/songtools/pkg/theory"

func Transpose(cfg *theory.Config, src SongReader, degreeClassInterval int, pitchClassInterval int) *TransposingSongReader {
	return &TransposingSongReader{
		cfg:                 cfg,
		src:                 src,
		degreeClassInterval: degreeClassInterval,
		pitchClassInterval:  pitchClassInterval,
	}
}

type TransposingSongReader struct {
	cfg                 *theory.Config
	src                 SongReader
	degreeClassInterval int
	pitchClassInterval  int
}

func (r *TransposingSongReader) NextLine() (Line, bool) {
	nl, ok := r.src.NextLine()
	if !ok {
		return nl, false
	}

	switch tnl := nl.(type) {
	case *KeyDirectiveLine:
		newKey := theory.TransposeKey(r.cfg, tnl.Key, r.degreeClassInterval, r.pitchClassInterval)
		tnl.Key = newKey
	case *ChordLine:
		for _, seg := range tnl.Chords {
			seg.Chord = theory.TransposeChord(r.cfg, seg.Chord, r.degreeClassInterval, r.pitchClassInterval)
			// TODO: deal with length
		}
	}

	return nl, ok
}

func (r *TransposingSongReader) Err() error {
	return r.src.Err()
}
