package songio

func NewMemory(lines []Line) *Memory {
	return &Memory{
		lines: lines,
	}
}

type Memory struct {
	lines []Line
	i     int
}

func (s *Memory) Lines() []Line {
	return s.lines
}

func (s *Memory) Next() (Line, bool) {
	if s.i < len(s.lines) {
		s.i++
		return s.lines[s.i-1], true
	}

	return nil, false
}

func (s *Memory) Err() error {
	return nil
}

func (s *Memory) Rewind() {
	s.i = 0
}
