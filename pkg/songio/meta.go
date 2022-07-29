package songio

import (
	"github.com/craiggwilson/songtool/pkg/theory/chord"
	"github.com/craiggwilson/songtool/pkg/theory/key"
	"github.com/craiggwilson/songtool/pkg/theory/note"
)

type Meta struct {
	Title    string        `json:"title"`
	Key      *key.Named    `json:"key"`
	Sections []string      `json:"sections"`
	Chords   []chord.Named `json:"chords"`
}

func ReadMeta(noteNamer note.Namer, src Song, full bool) (Meta, error) {
	var meta Meta

	chordSet := make(map[string]struct{})
Loop:
	for line, ok := src.Next(); ok; line, ok = src.Next() {
		switch tl := line.(type) {
		case *KeyDirectiveLine:
			meta.Key = &tl.Key
		case *TitleDirectiveLine:
			meta.Title = tl.Title
		case *ChordLine:
			if !full && meta.Key != nil {
				break Loop
			}

			for _, chordOffset := range tl.Chords {
				if meta.Key == nil {
					kind := key.KindMajor
					suffix := ""
					if chordOffset.Chord.Quality() == chord.QualityMinor {
						kind = key.KindMinor
						// TODO: minor kinda has to happen first in the name, so...?
						// I think we need a ParseKey that doesn't error if the full text isn't a key...
						suffix = chordOffset.Chord.Suffix[:1]
					}

					metaKey := key.Parsed{
						Key:    key.New(chordOffset.Chord.Root(), kind),
						Suffix: suffix,
					}

					meta.Key = &key.Named{
						Parsed: metaKey,
						Name:   metaKey.Name(noteNamer),
					}
					if !full {
						break Loop
					}
				}

				name := chordOffset.Chord.Name
				if _, ok := chordSet[name]; !ok {
					meta.Chords = append(meta.Chords, chordOffset.Chord)
					chordSet[name] = struct{}{}
				}
			}
		case *SectionStartDirectiveLine:
			if !full && meta.Key != nil {
				break Loop
			}

			meta.Sections = append(meta.Sections, tl.Name)
		case *TextLine, *SectionEndDirectiveLine:
			if !full && meta.Key != nil {
				break Loop
			}
		}
	}

	return meta, src.Err()
}
