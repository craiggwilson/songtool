package songio

func NewMemory(src Song) *Memory {
	return &Memory{
		src: src,
	}
}

type Memory struct {
	src   Song
	lines []Line
	i     int

	err error
}

func (s *Memory) Next() (Line, bool) {
	if s.i < len(s.lines) {
		s.i++
		return s.lines[s.i-1], true
	}

	line, ok := s.src.Next()
	if ok {
		s.lines = append(s.lines, line)
		s.i++
	}

	return line, ok
}

func (s *Memory) Err() error {
	return s.src.Err()
}

func (s *Memory) Rewind() {
	s.i = 0
}
