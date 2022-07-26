package theory2

import "github.com/craiggwilson/songtool/pkg/theory2/note"

type NoteNamer interface {
	NameNote(note.Note) string
}

func NameNote(n note.Note) string {
	return std.NameNote(n)
}

func ParseNote(text string) (note.Note, error) {
	return std.ParseNote(text)
}

type NoteParser interface {
	ParseNote(string) (note.Note, error)
}
