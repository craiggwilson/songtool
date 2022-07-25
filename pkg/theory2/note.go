package theory2

import "github.com/craiggwilson/songtool/pkg/theory2/note"

func NameNote(n note.Note) string {
	return std.NameNote(n)
}

type NoteNamer interface {
	NameNote(note.Note) string
}

type NoteNamerFunc func(note.Note) string

func (f NoteNamerFunc) NameNote(n note.Note) string {
	return f(n)
}
