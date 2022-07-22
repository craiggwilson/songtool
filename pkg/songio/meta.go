package songio

import "github.com/craiggwilson/songtool/pkg/theory"

type Meta struct {
	Title    string         `json:"title"`
	Key      theory.Key     `json:"key"`
	Sections []string       `json:"sections"`
	Chords   []theory.Chord `json:"chords"`
}

func ReadMeta(cfg *theory.Config, src Song, full bool) (Meta, error) {
	if cfg == nil {
		defCfg := theory.DefaultConfig()
		cfg = &defCfg
	}
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

			for _, chord := range tl.Chords {
				if !meta.Key.Note.IsValid() {
					kind := theory.KeyMajor
					suffix := ""
					if chord.IsMinor() {
						kind = theory.KeyMinor
						suffix = string(cfg.MinorKeySymbols[0])
					}

					meta.Key = theory.Key{
						Note:   chord.Root,
						Suffix: suffix,
						Kind:   kind,
					}
					if !full {
						break Loop
					}
				}

				name := chord.Name()
				if _, ok := chordSet[name]; !ok {
					meta.Chords = append(meta.Chords, chord.Chord)
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
