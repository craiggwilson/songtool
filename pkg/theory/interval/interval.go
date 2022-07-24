package interval

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/craiggwilson/songtool/pkg/theory/intervalquality"
)

type Interval struct {
	Diatonic  int
	Chromatic int
}

func (i Interval) Quality() intervalquality.IntervalQuality {
	diff := i.Chromatic - diatonicToChromatic[i.Diatonic]

	switch i.Diatonic {
	case 0:
		fallthrough
	case 3:
		fallthrough
	case 4:
		if diff == 0 {
			return intervalquality.Perfect()
		}
	default:
		if diff == 0 {
			return intervalquality.Major()
		}

		if diff == -1 {
			return intervalquality.Minor()
		}
	}

	if diff > 6 {
		diff = 12 - diff
	} else if diff < -6 {
		diff = 12 + diff
	}

	if diff > 0 {
		return intervalquality.Augmented(diff)
	}

	return intervalquality.Diminished(-diff)
}

func (i Interval) String() string {
	return fmt.Sprintf("%d%s", i.Diatonic+1, i.Quality())
}

func Must(interval Interval, err error) Interval {
	if err != nil {
		panic(err)
	}

	return interval
}

func Parse(text string) (Interval, error) {
	if len(text) < 2 {
		return Interval{}, fmt.Errorf("intervals must contain at least 2 characters, but had %d", len(text))
	}

	if !unicode.IsDigit(rune(text[0])) {
		return Interval{}, fmt.Errorf("expected number as pos 0, but got %q", text[0])
	}

	diatonic, _ := strconv.Atoi(string(text[0]))
	if diatonic < 1 || diatonic > 7 {
		return Interval{}, fmt.Errorf("expected a number between 1 and 7, but got %d", diatonic)
	}
	diatonic-- // normalizing for 0-based

	var q intervalquality.IntervalQuality
	switch text[1] {
	case 'P':
		if diatonic != 0 && diatonic != 3 && diatonic != 4 {
			return Interval{}, fmt.Errorf("only 1, 4, and 5 can be perfect, but got %d", diatonic+1)
		}
		q = intervalquality.Perfect()
		if len(text) > 2 {
			return Interval{}, fmt.Errorf("perfect interval quality has no size")
		}
	case 'M':
		if diatonic != 1 && diatonic != 2 && diatonic != 5 && diatonic != 6 {
			return Interval{}, fmt.Errorf("only 2, 3, 6, and 7 can be major, but got %d", diatonic+1)
		}
		q = intervalquality.Major()
		if len(text) > 2 {
			return Interval{}, fmt.Errorf("major interval quality has no size")
		}
	case 'm':
		if diatonic != 1 && diatonic != 2 && diatonic != 5 && diatonic != 6 {
			return Interval{}, fmt.Errorf("only 2, 3, 6, and 7 can be minor, but got %d", diatonic+1)
		}
		q = intervalquality.Minor()
		if len(text) > 2 {
			return Interval{}, fmt.Errorf("minor interval quality has no size")
		}
	case 'a':
		size := 1
		for i := 2; i < len(text); i++ {
			if text[i] == text[1] {
				size++
			} else {
				return Interval{}, fmt.Errorf("cannot mix interval qualities; expected %q, but got %q at pos %d", text[1], text[i], i)
			}
		}
		q = intervalquality.Augmented(size)
	case 'd':
		size := 1
		for i := 2; i < len(text); i++ {
			if text[i] == text[1] {
				size++
			} else {
				return Interval{}, fmt.Errorf("cannot mix interval qualities; expected %q, but got %q at pos %d", text[1], text[i], i)
			}
		}
		q = intervalquality.Diminished(size)
	default:
		return Interval{}, fmt.Errorf(`unknown interval quality at pos 1, expected ["P" "M" "m" "a" "d"], but got %q`, text[1])
	}

	chromatic := normalizeChromatic(chromaticFromDiatonicAndQuality(diatonic, q))

	return Interval{diatonic, chromatic}, nil
}

var (
	diatonicToChromatic = [7]int{0, 2, 4, 5, 7, 9, 11}
)

func chromaticFromDiatonicAndQuality(diatonic int, q intervalquality.IntervalQuality) int {
	chromatic := diatonicToChromatic[diatonic]

	switch q.Kind {
	case intervalquality.PerfectKind, intervalquality.MajorKind:
		return chromatic
	case intervalquality.MinorKind:
		return chromatic - 1
	case intervalquality.DiminishedKind:
		return chromatic - (q.Size + 1)
	case intervalquality.AugmentedKind:
		return chromatic + (q.Size + 1)
	default:
		panic(fmt.Sprintf("unknown interval quality kind, %d", q.Kind))
	}
}

func normalizeChromatic(chromatic int) int {
	return (chromatic + 12) % 12
}
