package songio

import (
	"fmt"
	"io"
	"strings"

	"github.com/craiggwilson/songtool/pkg/theory"
)

func WriteChordsOverLyrics(cfg *theory.Config, src LineIter, w io.Writer) (int, error) {
	n := 0
	i := 0
	var sb strings.Builder
	for line, ok := src.Next(); ok; line, ok = src.Next() {
		sb.Reset()
		switch tl := line.(type) {
		case *SectionStartDirectiveLine:
			sb.WriteString("[")
			sb.WriteString(tl.Name)
			sb.WriteString("]")
		case *KeyDirectiveLine:
			sb.WriteString("#key=")
			sb.WriteString(tl.Key.Name())
		case *TitleDirectiveLine:
			sb.WriteString("#title=")
			sb.WriteString(tl.Title)
		case *UnknownDirectiveLine:
			sb.WriteString("#")
			sb.WriteString(tl.Name)
			if len(tl.Value) > 0 {
				sb.WriteString("=")
				sb.WriteString(tl.Value)
			}
		case *TextLine:
			sb.WriteString(tl.Text)
		case *ChordLine:
			currentOffset := 0
			for _, chord := range tl.Chords {
				offsetDiff := chord.Offset - currentOffset
				if offsetDiff > 0 {
					sb.WriteString(strings.Repeat(" ", offsetDiff))
					currentOffset += offsetDiff
				}

				chordName := chord.Name()
				sb.WriteString(chordName)
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
