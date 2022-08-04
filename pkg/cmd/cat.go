package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/craiggwilson/songtool/pkg/songio"
)

type CatCmd struct {
	songCmd

	NoChords bool  `name:"no-chords" help:"Hides chords from the output."`
	JSON     bool  `name:"json" xor:"json" help:"Prints the output as JSON."`
	Color    color `name:"color" xor:"json" default:"${color}" negatable:"" help:"Indicates whether to use color"`
}

func (cmd *CatCmd) Run(cfg *Config) error {
	defer cmd.ensurePath().Close()

	song := cmd.openSong(cfg)

	if cmd.NoChords {
		song = songio.RemoveChords(song)
	}

	if cmd.JSON {
		return cmd.printSongJSON(song)
	}

	return cmd.printSong(cfg, song)
}

func (cmd *songCmd) printSong(cfg *Config, song songio.Song) error {
	for line, ok := song.Next(); ok; line, ok = song.Next() {
		switch tl := line.(type) {
		case *songio.TitleDirectiveLine:
			fmt.Println(directiveStyle.Render(fmt.Sprintf("#title=%s", tl.Title)))
		case *songio.KeyDirectiveLine:
			fmt.Println(directiveStyle.Render(fmt.Sprintf("#key=%s", chordStyle.Render(tl.Key.Name))))
		case *songio.UnknownDirectiveLine:
			fmt.Printf(directiveStyle.Render(fmt.Sprintf("#%s", tl.Name)))
			if len(tl.Value) > 0 {
				fmt.Printf(directiveStyle.Render(fmt.Sprintf("=%s", tl.Value)))
			}
			fmt.Println()
		case *songio.SectionStartDirectiveLine:
			fmt.Println(sectionNameStyle.Render(tl.Name))
		case *songio.SectionEndDirectiveLine:
			fmt.Println()
		case *songio.TextLine:
			fmt.Println(lyricsStyle.Render(tl.Text))
		case *songio.ChordLine:
			row := ""
			currentOffset := 0
			for _, chordOffset := range tl.Chords {
				offsetDiff := chordOffset.Offset - currentOffset
				if offsetDiff > 0 {
					row += strings.Repeat(" ", offsetDiff)
					currentOffset += offsetDiff
				}

				chordName := chordOffset.Chord.Name
				row += chordStyle.Render(chordName)
				currentOffset += len(chordName)
			}

			fmt.Println(row)
		default:
			fmt.Println()
		}
	}

	return nil
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
