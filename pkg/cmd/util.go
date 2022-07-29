package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/craiggwilson/songtool/pkg/songio"
	"github.com/craiggwilson/songtool/pkg/theory/note"
)

func WriteChordsOverLyricsWithHighlighter(styles *Styles, noteNamer note.Namer, src songio.Song, w io.Writer) (int, error) {
	n := 0
	i := 0
	var sb strings.Builder
	for line, ok := src.Next(); ok; line, ok = src.Next() {
		sb.Reset()
		switch tl := line.(type) {
		case *songio.SectionStartDirectiveLine:
			sb.WriteString(styles.Section.Renderf("[%s]", tl.Name))
		case *songio.KeyDirectiveLine:
			sb.WriteString(styles.Directive.Render("#key="))
			sb.WriteString(styles.Chord.Render(tl.Key.Name(noteNamer)))
		case *songio.TitleDirectiveLine:
			sb.WriteString(styles.Directive.Renderf("#title=%s", tl.Title))
		case *songio.UnknownDirectiveLine:
			sb.WriteString(styles.Directive.Renderf("#%s", tl.Name))
			if len(tl.Value) > 0 {
				sb.WriteString(styles.Directive.Renderf("=%s", tl.Value))
			}
		case *songio.TextLine:
			sb.WriteString(styles.Text.Render(tl.Text))
		case *songio.ChordLine:
			currentOffset := 0
			for _, chordOffset := range tl.Chords {
				offsetDiff := chordOffset.Offset - currentOffset
				if offsetDiff > 0 {
					sb.WriteString(strings.Repeat(" ", offsetDiff))
					currentOffset += offsetDiff
				}

				chordName := chordOffset.Chord.Name(noteNamer)
				sb.WriteString(styles.Chord.Render(chordName))
				currentOffset += len(chordName)
			}
		}

		sb.WriteByte('\n')

		w, err := io.WriteString(w, sb.String())
		n += w
		if err != nil {
			return n + w, fmt.Errorf("writing line %d: %w", i, err)
		}

		i++
	}

	return n, nil
}
