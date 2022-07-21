package songio

import "github.com/craiggwilson/songtools/theory"

type Segment interface {
	Len() int
	Offset() int

	seg()
}

type ChordSegment struct {
	Chord theory.Chord

	length int
	offset int
}

func (s *ChordSegment) seg() {}

func (s *ChordSegment) Len() int {
	return s.length
}

func (s *ChordSegment) Offset() int {
	return s.offset
}

type LyricSegment struct {
	Lyric string

	length int
	offset int
}

func (s *LyricSegment) seg() {}

func (s *LyricSegment) Len() int {
	return s.length
}

func (s *LyricSegment) Offset() int {
	return s.offset
}
