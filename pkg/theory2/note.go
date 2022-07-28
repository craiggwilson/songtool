package theory2

import "github.com/craiggwilson/songtool/pkg/theory2/note"

func NameNote(n note.Note) string {
	return std.NameNote(n)
}

func ParseNote(text string) (note.Note, error) {
	return std.ParseNote(text)
}
