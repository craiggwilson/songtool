package theory2

import (
	"fmt"
	"io"
	"strings"

	"github.com/craiggwilson/songtool/pkg/theory2/chord"
	"github.com/craiggwilson/songtool/pkg/theory2/interval"
	"github.com/craiggwilson/songtool/pkg/theory2/key"
	"github.com/craiggwilson/songtool/pkg/theory2/note"
	"github.com/craiggwilson/songtool/pkg/theory2/scale"
)

var std = func() *Theory {
	cfg := DefaultConfig()

	return New(cfg)
}()

func Default() *Theory {
	return std
}

func New(cfg *Config) *Theory {
	return &Theory{cfg: cfg}
}

type Theory struct {
	cfg *Config
}

func (t *Theory) ListScales() []ScaleMeta {
	result := make([]ScaleMeta, 0, len(t.cfg.Scales))
	for k, v := range t.cfg.Scales {
		result = append(result, ScaleMeta{k, v})
	}

	return result
}

func (t *Theory) LookupScale(name string) (ScaleMeta, bool) {
	intervals, ok := t.cfg.Scales[name]
	if !ok {
		return ScaleMeta{}, false
	}

	return ScaleMeta{name, intervals}, true
}

func (t *Theory) NameChord(c chord.Chord) string {
	panic("not implemented")
}

func (t *Theory) NameKey(k key.Key) string {
	name := t.NameNote(k.Note())
	if k.Kind() == key.KindMajor && len(t.cfg.MajorKeySymbols) > 0 {
		name += t.cfg.MajorKeySymbols[0]
	}
	if k.Kind() == key.KindMinor && len(t.cfg.MinorKeySymbols) > 0 {
		name += t.cfg.MinorKeySymbols[0]
	}

	return name
}

func (t *Theory) NameNote(n note.Note) string {
	degreeClass := n.DegreeClass()
	pitchClass := degreeClassToPitchClass[degreeClass]
	accidentals := n.PitchClass() - pitchClass

	if accidentals > 6 {
		accidentals -= 12
	} else if accidentals < -6 {
		accidentals += 12
	}

	accidentalStr := ""
	if accidentals > 0 {
		accidentalStr = strings.Repeat(t.cfg.SharpSymbols[0], accidentals)
	} else if accidentals < 0 {
		accidentalStr = strings.Repeat(t.cfg.FlatSymbols[0], -accidentals)
	}

	natural := t.cfg.NaturalNoteNames[degreeClass]
	return natural + accidentalStr
}

func (t *Theory) ParseChord(text string) (chord.Parsed, error) {
	root, pos, err := t.parseNote(text, 0)
	if err != nil {
		return chord.Parsed{}, err
	}

	intervalMap := make(map[interval.Interval]struct{})

	suffix := text[pos:]
	suffixPos := pos
	for _, m := range t.cfg.ChordModifiers {
		matched := false
		if m.Match == nil {
			matched = true
		} else {
			match := m.Match.FindStringSubmatch(suffix)
			if len(match) > 0 && (m.Except == nil || !m.Except.MatchString(suffix)) {
				matched = true
				if len(match) == 1 {
					// If there are no groups, use the full match.
					pos += len(match[0])
				} else {
					// If there are named capture groups, there must be one called "mod" which will be used.
					captureNames := m.Match.SubexpNames()
					found := false
					for i, name := range captureNames {
						if name == "mod" {
							pos += len(match[i])
							found = true
							break
						}
					}
					if !found {
						pos += len(match[1])
					}
				}
			}
		}
		if matched {
			for _, add := range m.Add {
				intervalMap[add] = struct{}{}
			}
			for _, rm := range m.Remove {
				delete(intervalMap, rm)
			}
		}
	}

	intervals := make([]interval.Interval, 0, len(intervalMap))
	for k := range intervalMap {
		intervals = append(intervals, k)
	}

	interval.Sort(intervals)

	delimiterPos := pos
	base, delim, pos, _ := t.parseBaseNote(text, pos)

	if len(text) > pos {
		return chord.Parsed{}, fmt.Errorf("expected EOF at position %d, but had %s", pos, text[pos:])
	}

	return chord.Parsed{
		Chord:             chord.New(root, base, intervals...),
		Suffix:            text[suffixPos:delimiterPos],
		BaseNoteDelimiter: delim,
	}, nil
}

func (t *Theory) ParseKey(text string) (key.Parsed, error) {
	found := false
	kind := key.KindMajor
	suffix := ""
	for _, sym := range t.cfg.MajorKeySymbols {
		idx := strings.Index(text, sym)
		if idx > 0 {
			text = text[:idx]
			kind = key.KindMajor
			suffix = sym
			found = true
			break
		}
	}

	if !found {
		for _, sym := range t.cfg.MinorKeySymbols {
			idx := strings.Index(text, sym)
			if idx > 0 {
				text = text[:idx]
				kind = key.KindMinor
				suffix = sym
				break
			}
		}
	}

	n, err := t.ParseNote(text)
	if err != nil {
		return key.Parsed{}, err
	}

	return key.Parsed{
		Key:    key.New(n, kind),
		Suffix: suffix,
	}, nil
}

func (t *Theory) ParseNote(text string) (note.Note, error) {
	n, pos, err := t.parseNote(text, 0)
	if err != nil {
		return note.Note{}, err
	}

	if len(text) != pos {
		return note.Note{}, fmt.Errorf("expected EOF at position %d, but had %q", pos, text[pos:])
	}

	return n, err
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

func (t *Theory) degreeClassFromNaturalNoteName(naturalNoteName string) int {
	for i, nn := range t.cfg.NaturalNoteNames {
		if nn == naturalNoteName {
			return i
		}
	}

	panic(fmt.Sprintf("natural note name %q does not map to a degree class", naturalNoteName))
}

func (t *Theory) parseAccidentals(text string, pos int) (int, int) {
	if len(text) <= pos {
		return 0, pos
	}

	accidentals, pos := t.parseSharps(text, pos)
	if accidentals != 0 {
		return accidentals, pos
	}

	return t.parseFlats(text, pos)
}

func (t *Theory) parseBaseNote(text string, pos int) (*note.Note, string, int, error) {
	if len(text) <= pos {
		return nil, "", pos, io.ErrUnexpectedEOF
	}

	for _, r := range t.cfg.BaseNoteDelimiters {
		if strings.HasPrefix(text[pos:], r) {
			note, pos, err := t.parseNote(text, pos+len(r))
			return &note, r, pos, err
		}
	}

	return nil, "", pos, fmt.Errorf("expected one of %q, but got %q", t.cfg.BaseNoteDelimiters, text[pos:])
}

func (t *Theory) parseNote(text string, pos int) (note.Note, int, error) {
	naturalNoteName, newPos, err := t.parseNaturalNoteName(text, pos)
	if err != nil {
		return note.Note{}, pos, fmt.Errorf("expected natural note name at position %d: %w", newPos, err)
	}

	degreeClass := t.degreeClassFromNaturalNoteName(naturalNoteName)
	pitchClass := degreeClassToPitchClass[degreeClass]
	accidentals, newPos := t.parseAccidentals(text, newPos)

	return note.New(degreeClass, pitchClass+accidentals), newPos, nil
}

func (t *Theory) parseNaturalNoteName(text string, pos int) (string, int, error) {
	if len(text) <= pos {
		return "", pos, io.ErrUnexpectedEOF
	}

	for _, nn := range t.cfg.NaturalNoteNames {
		if strings.HasPrefix(text[pos:], nn) {
			return nn, pos + len(nn), nil
		}
	}

	return "", pos, fmt.Errorf("expected one of %q, but got %q", t.cfg.NaturalNoteNames, text[pos:])
}

func (t *Theory) parseSharps(text string, pos int) (int, int) {
	if len(text) <= pos {
		return 0, pos
	}

	accidentals := 0

	changed := true
	for changed {
		changed = false
		for _, sym := range t.cfg.SharpSymbols {
			if strings.HasPrefix(text[pos:], sym) {
				accidentals++
				pos += len(sym)
				changed = true
				break
			}
		}
	}

	return accidentals, pos
}

func (t *Theory) parseFlats(text string, pos int) (int, int) {
	if len(text) <= pos {
		return 0, pos
	}

	accidentals := 0

	changed := true
	for changed {
		changed = false
		for _, sym := range t.cfg.FlatSymbols {
			if strings.HasPrefix(text[pos:], sym) {
				accidentals--
				pos += len(sym)
				changed = true
				break
			}
		}
	}

	return accidentals, pos
}
