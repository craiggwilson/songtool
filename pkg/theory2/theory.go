package theory2

var std = func() *Theory {
	cfg := DefaultConfig()

	return New(cfg)
}()

func New(noteNamer NoteNamer) *Theory {
	return &Theory{
		NoteNamer: noteNamer,
	}
}

func NewDefault() *Theory {
	return New(DefaultConfig())
}

type Theory struct {
	NoteNamer
}
