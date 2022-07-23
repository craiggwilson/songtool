package songio

import "github.com/craiggwilson/songtool/pkg/theory"

type Meta struct {
	Title    string         `json:"title"`
	Key      theory.Key     `json:"key"`
	Sections []string       `json:"sections"`
	Chords   []theory.Chord `json:"chords"`
}

func ReadMeta(t *theory.Theory, src Song, full bool) (Meta, error) {
	var meta Meta

	chordSet := make(map[string]struct{})
Loop:
	for line, ok := src.Next(); ok; line, ok = src.Next() {
		switch tl := line.(type) {
		case *KeyDirectiveLine:
			meta.Key = tl.Key
		case *TitleDirectiveLine:
			meta.Title = tl.Title
		case *ChordLine:
			if !full && meta.Key.Note.IsValid() {
				break Loop
			}

			for _, chordOffset := range tl.Chords {
				if !meta.Key.Note.IsValid() {
					kind := theory.KeyMajor
					suffix := ""
					if chordOffset.Chord.IsMinor() {
						kind = theory.KeyMinor
						suffix = string(t.Config.MinorKeySymbols[0])
					}

					meta.Key = theory.Key{
						Note:   chordOffset.Chord.Root,
						Suffix: suffix,
						Kind:   kind,
					}
					if !full {
						break Loop
					}
				}

				name := chordOffset.Chord.Name()
				if _, ok := chordSet[name]; !ok {
					meta.Chords = append(meta.Chords, chordOffset.Chord)
					chordSet[name] = struct{}{}
				}
			}
		case *SectionStartDirectiveLine:
			if !full && meta.Key.Note.IsValid() {
				break Loop
			}

			meta.Sections = append(meta.Sections, tl.Name)
		case *TextLine, *SectionEndDirectiveLine:
			if !full && meta.Key.Note.IsValid() {
				break Loop
			}
		}
	}

	return meta, src.Err()
}
