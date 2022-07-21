package songio

import "github.com/craiggwilson/songtools/theory"

type Segment interface {
	Len() int
	Offset() int

	seg()
}

type ChordSegment struct {
	Chord theory.Chord

	offset int
}

func (s *ChordSegment) seg() {}

func (s *ChordSegment) Len() int {
	return len(s.Chord.Name())
}

func (s *ChordSegment) Offset() int {
	return s.offset
}

type LyricSegment struct {
	Lyric string

	offset int
}

func (s *LyricSegment) seg() {}

func (s *LyricSegment) Len() int {
	return len(s.Lyric)
}

func (s *LyricSegment) Offset() int {
	return s.offset
}
