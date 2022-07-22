package songio

type Rewinder struct {
	src   Song
	lines []Line
}

func NewRewinder(src Song) *Rewinder {
	return &Rewinder{
		src: src,
	}
}

func (s *Rewinder) Next() (Line, bool) {
	line, ok := s.src.Next()
	if ok {
		s.lines = append(s.lines, line)
	}

	return line, ok
}

func (s *Rewinder) Err() error {
	return s.src.Err()
}

func (s *Rewinder) Rewind() Song {
	return &rewoundSong{
		src:   s.src,
		lines: s.lines,
	}
}

type rewoundSong struct {
	src   Song
	lines []Line
	i     int
}

func (s *rewoundSong) Next() (Line, bool) {
	if s.i < len(s.lines) {
		s.i++
		return s.lines[s.i-1], true
	}

	return s.src.Next()
}

func (s *rewoundSong) Err() error {
	return s.src.Err()
}
