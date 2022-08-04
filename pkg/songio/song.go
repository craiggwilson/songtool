package songio

import "strings"

type Song interface {
	Next() (Line, bool)
	Err() error
}

func ReadAllLines(s Song) ([]Line, error) {
	var lines []Line
	for nl, ok := s.Next(); ok; nl, ok = s.Next() {
		lines = append(lines, nl)
	}

	return lines, s.Err()
}

func RemoveChords(src Song) Song {
	return &chordsRemover{src: src}
}

type chordsRemover struct {
	src Song
}

func (s *chordsRemover) Next() (Line, bool) {
	line, ok := s.src.Next()
	if !ok {
		return line, false
	}

	switch tl := line.(type) {
	case *ChordLine:
		return s.Next()
	case *TextLine:
		tl.Text = strings.TrimSpace(tl.Text)
		return tl, true
	default:
		return line, true
	}
}

func (s *chordsRemover) Err() error {
	return s.src.Err()
}
