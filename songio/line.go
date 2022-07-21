package songio

type Line interface {
	line()
}

type EmptyLine struct{}

func (EmptyLine) line() {}

type ChordLine struct {
	Segments []ChordSegment
}

func (l *ChordLine) line() {}

type LyricLine struct {
	Segments []LyricSegment
}

func (l *LyricLine) line() {}
