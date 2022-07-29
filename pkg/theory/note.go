package theory

import "github.com/craiggwilson/songtool/pkg/theory/note"

func NameNote(n note.Note) string {
	return std.NameNote(n)
}

func ParseNote(text string) (note.Note, error) {
	return std.ParseNote(text)
}
