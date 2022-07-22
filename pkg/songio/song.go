package songio

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
