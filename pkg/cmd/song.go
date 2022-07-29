package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/craiggwilson/songtool/pkg/songio"
)

type songCmd struct {
	Format string   `name:"format" enum:"auto,chordsOverLyrics" default:"auto" help:"The format of the song; defaults to 'auto'."`
	Path   *os.File `arg:"" optional:"" help:"The path to the song; '-' can be used for stdin."`
}

func (cmd *songCmd) ensurePath() *os.File {
	if cmd.Path == nil {
		cmd.Path = os.Stdin
	}
	return cmd.Path
}

func (cmd *songCmd) openSong(cfg *Config) songio.Song {
	return songio.ReadChordsOverLyrics(cfg.Theory, cfg.Theory, cmd.Path)
}

func (cmd *songCmd) printSong(styles *Styles, song songio.Song) error {
	_, err := writeChordsOverLyricsWithHighlighter(styles, song, os.Stdout)
	return err
}

func (cmd *songCmd) printSongJSON(song songio.Song) error {
	lines, err := songio.ReadAllLines(song)
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(lines, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}

func writeChordsOverLyricsWithHighlighter(styles *Styles, src songio.Song, w io.Writer) (int, error) {
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
			sb.WriteString(styles.Chord.Render(tl.Key.Name))
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

				chordName := chordOffset.Chord.Name
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
