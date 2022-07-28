package note

type Namer interface {
	NameNote(Note) string
}
