package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/craiggwilson/songtool/pkg/songio"
)

type CatCmd struct {
	songCmd

	JSON  bool `name:"json" xor:"json" help:"Prints the output as JSON."`
	Color bool `name:"color" xor:"json" help:"Indicates whether to use color." negatable:""`
}

func (cmd *CatCmd) Run(cfg *Config) error {
	defer cmd.ensurePath().Close()

	song := cmd.openSong()

	if cmd.JSON {
		return cmd.printJSON(song)
	}

	return cmd.print(cfg, song)
}

func (cmd *CatCmd) print(cfg *Config, song songio.Song) error {
	_, err := WriteChordsOverLyricsWithHighlighter(song, os.Stdout, cfg.Highlighter)
	return err
}

func (cmd *CatCmd) printJSON(song songio.Song) error {
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

func WriteChordsOverLyricsWithHighlighter(src songio.Song, w io.Writer, highlighter Highlighter) (int, error) {
	n := 0
	i := 0
	var sb strings.Builder
	for line, ok := src.Next(); ok; line, ok = src.Next() {
		sb.Reset()
		switch tl := line.(type) {
		case *songio.SectionStartDirectiveLine:
			sb.WriteString("[")
			sb.WriteString(highlighter.SectionColor.Sprint(tl.Name))
			sb.WriteString("]")
		case *songio.KeyDirectiveLine:
			sb.WriteString(highlighter.DirectiveColor.Sprint("#key"))
			sb.WriteString("=")
			sb.WriteString(highlighter.ChordColor.Sprint(tl.Key.Name()))
		case *songio.TitleDirectiveLine:
			sb.WriteString(highlighter.DirectiveColor.Sprint("#title"))
			sb.WriteString("=")
			sb.WriteString(tl.Title)
		case *songio.UnknownDirectiveLine:
			sb.WriteString(highlighter.DirectiveColor.Sprint("#"))
			sb.WriteString(highlighter.DirectiveColor.Sprint(tl.Name))
			if len(tl.Value) > 0 {
				sb.WriteString("=")
				sb.WriteString(tl.Value)
			}
		case *songio.TextLine:
			sb.WriteString(tl.Text)
		case *songio.ChordLine:
			currentOffset := 0
			for _, chordOffset := range tl.Chords {
				offsetDiff := chordOffset.Offset - currentOffset
				if offsetDiff > 0 {
					sb.WriteString(strings.Repeat(" ", offsetDiff))
					currentOffset += offsetDiff
				}

				chordName := chordOffset.Chord.Name()
				sb.WriteString(highlighter.ChordColor.Sprint(chordName))
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
