package songio

type SongReader interface {
	NextLine() (Line, bool)
	Err() error
}

func ReadAllLines(r SongReader) ([]Line, error) {
	var lines []Line
	for nl, ok := r.NextLine(); ok; nl, ok = r.NextLine() {
		lines = append(lines, nl)
	}

	return lines, r.Err()
}
