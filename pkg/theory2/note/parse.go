package note

type Parser interface {
	ParseNote(string) (Note, error)
}
