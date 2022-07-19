package key

import (
	"github.com/craiggwilson/songtools/theory/note"
)

type Kind int

const (
	Major Kind = iota + 1
	Minor
)

type Key struct {
	Note note.Note
	Kind Kind
}
