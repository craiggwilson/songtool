package interval

import (
	"fmt"
	"sort"
	"strconv"
	"unicode"
)

var (
	diatonicToChromatic = [7]int{0, 2, 4, 5, 7, 9, 11}
)

func Augmented(diatonic, size int) Interval {
	i, err := augmentedErr(diatonic, size)
	if err != nil {
		panic(err)
	}
	return i
}

func Diminished(diatonic, size int) Interval {
	i, err := diminishedErr(diatonic, size)
	if err != nil {
		panic(err)
	}

	return i
}

func Major(diatonic int) Interval {
	i, err := majorErr(diatonic)
	if err != nil {
		panic(err)
	}

	return i
}

func Minor(diatonic int) Interval {
	i, err := minorErr(diatonic)
	if err != nil {
		panic(err)
	}

	return i
}

func Must(interval Interval, err error) Interval {
	if err != nil {
		panic(err)
	}

	return interval
}

func New(diatonic int, q Quality) Interval {
	return Interval{normalizeDiatonic(diatonic), q}
}

func NewWithChromatic(diatonic, chromatic int) Interval {
	diatonic = normalizeDiatonic(diatonic)
	return Interval{diatonic, qualityFromDiatonicAndChromatic(diatonic, chromatic)}
}

func Parse(text string) (Interval, error) {
	if len(text) < 2 {
		return Interval{}, fmt.Errorf("intervals must contain at least 2 characters, but had %d", len(text))
	}

	if !unicode.IsDigit(rune(text[0])) {
		return Interval{}, fmt.Errorf("expected number as pos 0, but got %q", text[0])
	}

	digits := 1
	if unicode.IsDigit(rune(text[1])) {
		digits++

		if len(text) < 2 {
			return Interval{}, fmt.Errorf("intervals must contain a quality, but had none")
		}
	}

	diatonic, _ := strconv.Atoi(string(text[:digits]))
	if diatonic < 1 || diatonic > 13 {
		return Interval{}, fmt.Errorf("expected a number between 1 and 13, but got %d", diatonic)
	}
	switch text[digits] {
	case 'P':
		if !canDiatonicBePerfect(diatonic - 1) {
			return Interval{}, fmt.Errorf("only 1, 4, 5, and 11 can be perfect, but got %d", diatonic)
		}
		if len(text) > digits+1 {
			return Interval{}, fmt.Errorf("perfect interval quality has no size")
		}
		return perfectErr(diatonic - 1)
	case 'M':
		if !canDiatonicBeMajorMinor(diatonic - 1) {
			return Interval{}, fmt.Errorf("only 2, 3, 6, 7, 9, 10, and 13 can be major, but got %d", diatonic)
		}
		if len(text) > digits+1 {
			return Interval{}, fmt.Errorf("major interval quality has no size")
		}
		return majorErr(diatonic - 1)
	case 'm':
		if !canDiatonicBeMajorMinor(diatonic - 1) {
			return Interval{}, fmt.Errorf("only 2, 3, 6, 7, 9, 10, and 13 can be minor, but got %d", diatonic)
		}
		if len(text) > digits+1 {
			return Interval{}, fmt.Errorf("minor interval quality has no size")
		}
		return minorErr(diatonic - 1)
	case 'a':
		size := 1
		for i := digits + 1; i < len(text); i++ {
			if text[i] == text[digits] {
				size++
			} else {
				return Interval{}, fmt.Errorf("cannot mix interval qualities; expected %q, but got %q at pos %d", text[1], text[i], i)
			}
		}
		return augmentedErr(diatonic-1, size)
	case 'd':
		size := 1
		for i := digits + 1; i < len(text); i++ {
			if text[i] == text[digits] {
				size++
			} else {
				return Interval{}, fmt.Errorf("cannot mix interval qualities; expected %q, but got %q at pos %d", text[1], text[i], i)
			}
		}
		return diminishedErr(diatonic-1, size)
	default:
		return Interval{}, fmt.Errorf(`unknown interval quality at pos 1, expected ["P" "M" "m" "a" "d"], but got %q`, text[1])
	}
}

func Perfect(diatonic int) Interval {
	i, err := perfectErr(diatonic)
	if err != nil {
		panic(err)
	}

	return i
}

func Steps(intervals []Interval) [22]bool {
	var r [22]bool

	for _, iv := range intervals {
		r[iv.Chromatic()] = true
	}

	return r
}

func Sort(intervals []Interval) {
	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].diatonic < intervals[j].diatonic {
			return true
		} else if intervals[i].diatonic > intervals[j].diatonic {
			return false
		}

		return intervals[i].String() < intervals[j].String()
	})
}

type Interval struct {
	diatonic int
	quality  Quality
}

func (i Interval) Chromatic() int {
	c := chromaticFromDiatonicAndQuality(i.diatonic, i.quality)
	if c == 12 {
		c = 0
	}

	return c
}

func (i Interval) Diatonic() int {
	return i.diatonic
}

func (i *Interval) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

func (i Interval) Quality() Quality {
	return i.quality
}

func (i Interval) String() string {
	return fmt.Sprintf("%d%s", i.diatonic+1, i.quality)
}

func (i Interval) Transpose(other Interval) Interval {
	newDiatonic := normalizeDiatonic(i.diatonic+other.diatonic) % 7
	newChromatic := normalizeChromatic(i.Chromatic()+other.Chromatic()) % 12
	return New(newDiatonic, qualityFromDiatonicAndChromatic(newDiatonic, newChromatic))
}

func (i *Interval) UnmarshalText(text []byte) error {
	intval, err := Parse(string(text))
	if err != nil {
		return err
	}

	*i = intval
	return nil
}

func canDiatonicBeMajorMinor(diatonic int) bool {
	switch diatonic {
	case 1, 2, 5, 6, 8, 9, 11, 12:
		return true
	default:
		return false
	}
}

func canDiatonicBePerfect(diatonic int) bool {
	switch diatonic {
	case 0, 3, 4, 10:
		return true
	default:
		return false
	}
}

func chromaticFromDiatonicAndQuality(diatonic int, q Quality) int {
	chromatic := diatonicToChromatic[diatonic%7]
	if diatonic > 7 {
		chromatic += 12
	}

	switch q.Kind() {
	case QualityKindPerfect, QualityKindMajor:
		return chromatic
	case QualityKindMinor:
		return chromatic - 1
	case QualityKindDiminished:
		value := chromatic - q.Size()
		if canDiatonicBeMajorMinor(diatonic) {
			value--
		}
		return value
	case QualityKindAugmented:
		return chromatic + q.Size()
	default:
		panic(fmt.Sprintf("unknown interval quality kind, %d", q.Kind()))
	}
}

func augmentedErr(diatonic, size int) (Interval, error) {
	q, err := AugmentedQuality(size)
	if err != nil {
		return Interval{}, err
	}
	return New(diatonic, q), nil
}

func diminishedErr(diatonic, size int) (Interval, error) {
	q, err := DiminishedQuality(size)
	if err != nil {
		return Interval{}, err
	}

	return New(diatonic, q), nil
}

func majorErr(diatonic int) (Interval, error) {
	if !canDiatonicBeMajorMinor(diatonic) {
		return Interval{}, fmt.Errorf("only 1, 2, 5, 6, 9, 10, and 13 can be major, but got %d", diatonic)
	}

	return New(diatonic, MajorQuality()), nil
}

func minorErr(diatonic int) (Interval, error) {
	if !canDiatonicBeMajorMinor(diatonic) {
		return Interval{}, fmt.Errorf("only 1, 2, 5, 6, 9, 10, 13 can be minor, but got %d", diatonic)
	}

	return New(diatonic, MinorQuality()), nil
}

func normalizeDiatonic(diatonic int) int {
	if diatonic < 0 {
		return (diatonic + 7) % 7
	}

	if diatonic > 7 {
		return 7 + (diatonic % 7)
	}

	return diatonic % 7
}

func normalizeChromatic(chromatic int) int {
	if chromatic < 0 {
		return (chromatic + 12) % 12
	}

	if chromatic > 12 {
		return 12 + (chromatic % 12)
	}

	return chromatic % 12
}

func qualityFromDiatonicAndChromatic(diatonic, chromatic int) Quality {
	if diatonic > 7 {
		diatonic %= 7
	}

	diff := (chromatic % 12) - diatonicToChromatic[diatonic]
	if diff > 6 {
		diff -= 12
	} else if diff < -6 {
		diff += 12
	}

	switch {
	case canDiatonicBePerfect(diatonic):
		if diff == 0 {
			return PerfectQuality()
		}
	default:
		if diff == 0 {
			return MajorQuality()
		}
		if diff == -1 {
			return MinorQuality()
		}

		if diff < -1 {
			diff++
		}
	}

	if diff > 0 {
		q, _ := AugmentedQuality(diff)
		return q
	}

	q, _ := DiminishedQuality(-diff)
	return q
}

func perfectErr(diatonic int) (Interval, error) {
	if !canDiatonicBePerfect(diatonic) {
		return Interval{}, fmt.Errorf("only 0, 3, 4, and 10 can be perfect, but got %d", diatonic)
	}

	return New(diatonic, PerfectQuality()), nil
}
