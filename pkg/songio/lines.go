package songio

func FromLines(lines []Line) *Lines {
	return &Lines{
		lines: lines,
	}
}

type Lines struct {
	lines []Line
	i     int
}

func (s *Lines) Lines() []Line {
	return s.lines
}

func (s *Lines) Next() (Line, bool) {
	if s.i < len(s.lines) {
		s.i++
		return s.lines[s.i-1], true
	}

	return nil, false
}

func (s *Lines) Err() error {
	return nil
}

func (s *Lines) Rewind() {
	s.i = 0
}
