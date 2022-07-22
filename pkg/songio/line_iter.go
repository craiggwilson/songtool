package songio

type LineIter interface {
	Next() (Line, bool)
	Err() error
}

func ReadAllLines(it LineIter) ([]Line, error) {
	var lines []Line
	for nl, ok := it.Next(); ok; nl, ok = it.Next() {
		lines = append(lines, nl)
	}

	return lines, it.Err()
}
