package theory2

import (
	"fmt"
	"strings"

	"github.com/craiggwilson/songtool/pkg/theory2/note"
	"github.com/craiggwilson/songtool/pkg/theory2/scale"
)

var std = func() *Theory {
	cfg := DefaultConfig()

	return New(cfg, cfg, cfg)
}()

func Default() *Theory {
	return std
}

func New(noteNamer NoteNamer, noteParser NoteParser, scaleProvider ScaleProvider) *Theory {
	return &Theory{
		noteNamer:     noteNamer,
		noteParser:    noteParser,
		scaleProvider: scaleProvider,
	}
}

type Theory struct {
	noteNamer     NoteNamer
	noteParser    NoteParser
	scaleProvider ScaleProvider
}

func (t *Theory) ListScales() []ScaleMeta {
	return t.scaleProvider.ListScales()
}

func (t *Theory) LookupScale(name string) (ScaleMeta, bool) {
	return t.scaleProvider.LookupScale(name)
}

func (t *Theory) NameNote(n note.Note) string {
	return t.noteNamer.NameNote(n)
}

func (t *Theory) ParseNote(text string) (note.Note, error) {
	return t.noteParser.ParseNote(text)
}

func (t *Theory) ParseScale(text string) (scale.Scale, error) {
	parts := strings.SplitN(text, " ", 2)

	root, err := t.ParseNote(parts[0])
	if err != nil {
		return scale.Scale{}, err
	}

	scaleName := "Major"
	if len(parts) == 2 {
		scaleName = parts[1]
	}

	meta, ok := t.LookupScale(scaleName)
	if !ok {
		return scale.Scale{}, fmt.Errorf("unknown scale name %q", scaleName)
	}

	return scale.Generate(fmt.Sprintf("%s %s", parts[0], meta.Name), root, meta.Intervals...), nil
}
