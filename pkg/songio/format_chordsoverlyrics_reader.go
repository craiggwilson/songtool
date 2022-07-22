package songio

import (
	"bufio"
	"io"
	"strings"
	"unicode"

	"github.com/craiggwilson/songtool/pkg/theory"
)

type ChordsOverLyricsReader struct {
	cfg     *theory.Config
	scanner *bufio.Scanner

	blankLineCount     int
	currentSectionName string

	saveLine Line

	err error
}

func ReadChordsOverLyrics(cfg *theory.Config, src io.Reader) *ChordsOverLyricsReader {
	return &ChordsOverLyricsReader{
		cfg:     cfg,
		scanner: bufio.NewScanner(src),
	}
}

func (r *ChordsOverLyricsReader) Next() (Line, bool) {
	if r.saveLine != nil {
		for r.blankLineCount > 0 {
			r.blankLineCount--
			return EmptyLine{}, true
		}

		line := r.saveLine
		r.saveLine = nil

		if ssdl, ok := line.(*SectionStartDirectiveLine); ok {
			r.currentSectionName = ssdl.Name
		}

		return line, true
	}

	if !r.scanner.Scan() {
		if len(r.currentSectionName) > 0 {
			if r.blankLineCount > 0 {
				r.blankLineCount--
			}

			currentSectionName := r.currentSectionName
			r.currentSectionName = ""
			return &SectionEndDirectiveLine{
				Name: currentSectionName,
			}, true
		}

		for r.blankLineCount > 0 {
			r.blankLineCount--
			return EmptyLine{}, true
		}

		return nil, false
	}

	line := r.parseLine(r.scanner.Text())

	switch tl := line.(type) {
	case EmptyLine:
		if len(r.currentSectionName) > 0 {
			if r.blankLineCount == 0 {
				r.blankLineCount++
				return r.Next()
			}

			currentSectionName := r.currentSectionName
			r.blankLineCount--
			r.currentSectionName = ""
			return &SectionEndDirectiveLine{
				Name: currentSectionName,
			}, true
		}
	case *SectionStartDirectiveLine:
		if len(r.currentSectionName) > 0 {
			r.saveLine = line
			if r.blankLineCount > 0 {
				r.blankLineCount--
			}

			currentSectionName := r.currentSectionName
			r.currentSectionName = ""

			return &SectionEndDirectiveLine{
				Name: currentSectionName,
			}, true
		} else {
			r.currentSectionName = tl.Name
		}
	}

	return line, true
}

func (r *ChordsOverLyricsReader) Err() error {
	return r.scanner.Err()
}

func (r *ChordsOverLyricsReader) parseContent(text string) Line {
	var chordSegments []*ChordOffset

	wordStartIdx := -1
	for i, n := range text {
		if unicode.IsSpace(n) {
			if wordStartIdx > -1 {
				chord, err := theory.ParseChord(r.cfg, text[wordStartIdx:i])
				if err != nil {
					return &TextLine{
						Text: text,
					}
				}

				chordSegments = append(chordSegments, &ChordOffset{
					Chord:  chord,
					Offset: wordStartIdx,
				})
				wordStartIdx = -1
			}
		} else if wordStartIdx == -1 {
			wordStartIdx = i
		}
	}

	if wordStartIdx > -1 {
		chord, err := theory.ParseChord(r.cfg, text[wordStartIdx:])
		if err != nil {
			return &TextLine{
				Text: text,
			}
		}

		chordSegments = append(chordSegments, &ChordOffset{
			Chord:  chord,
			Offset: wordStartIdx,
		})
	}

	return &ChordLine{
		Chords: chordSegments,
	}
}

func (r *ChordsOverLyricsReader) parseDirective(text string) Line {
	idx := strings.IndexRune(text, '=')
	if idx < 0 {
		return &UnknownDirectiveLine{
			Name: text[1:],
		}
	}

	name := text[1:idx]
	value := strings.TrimSpace(text[idx+1:])

	switch name {
	case "title":
		return &TitleDirectiveLine{
			Title: value,
		}
	case "key":
		if key, err := theory.ParseKey(r.cfg, value); err == nil {
			return &KeyDirectiveLine{
				Key: key,
			}
		}
	}

	return &UnknownDirectiveLine{
		Name:  name,
		Value: value,
	}
}

func (r *ChordsOverLyricsReader) parseLine(text string) Line {
	if isEmptyOrWhitespace(text) {
		return EmptyLine{}
	}

	switch text[0] {
	case '#':
		return r.parseDirective(text)
	case '[':
		return r.parseSectionStart(text)
	default:
		return r.parseContent(text)
	}
}

func (r *ChordsOverLyricsReader) parseSectionStart(text string) Line {
	idx := strings.IndexRune(text, ']')
	if idx < 0 || !isEmptyOrWhitespace(text[idx+1:]) {
		return r.parseContent(text)
	}

	return &SectionStartDirectiveLine{
		Name: text[1:idx],
	}
}
