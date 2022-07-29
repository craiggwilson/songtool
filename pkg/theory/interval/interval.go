package interval

import (
	"fmt"
	"sort"
	"strconv"
	"unicode"
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

func FromStep(step int) Interval {
	step = normalizeChromatic(step)
	for i := 0; i < len(diatonicToChromatic); i++ {
		if step <= diatonicToChromatic[i] {
			return Interval{i, qualityFromDiatonicAndChromatic(i, step)}
		}
	}

	panic(fmt.Sprintf("impossible step %d", step))
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
	chromatic = normalizeChromatic(chromatic)
	return Interval{diatonic, qualityFromDiatonicAndChromatic(diatonic, chromatic)}
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
	switch text[1] {
	case 'P':
		if diatonic != 1 && diatonic != 4 && diatonic != 5 {
			return Interval{}, fmt.Errorf("only 1, 4, and 5 can be perfect, but got %d", diatonic)
		}
		if len(text) > 2 {
			return Interval{}, fmt.Errorf("perfect interval quality has no size")
		}
		return perfectErr(diatonic - 1)
	case 'M':
		if diatonic != 2 && diatonic != 3 && diatonic != 6 && diatonic != 7 {
			return Interval{}, fmt.Errorf("only 2, 3, 6, and 7 can be major, but got %d", diatonic)
		}
		if len(text) > 2 {
			return Interval{}, fmt.Errorf("major interval quality has no size")
		}
		return majorErr(diatonic - 1)
	case 'm':
		if diatonic != 2 && diatonic != 3 && diatonic != 6 && diatonic != 7 {
			return Interval{}, fmt.Errorf("only 2, 3, 6, and 7 can be minor, but got %d", diatonic)
		}
		if len(text) > 2 {
			return Interval{}, fmt.Errorf("minor interval quality has no size")
		}
		return minorErr(diatonic - 1)
	case 'a':
		size := 1
		for i := 2; i < len(text); i++ {
			if text[i] == text[1] {
				size++
			} else {
				return Interval{}, fmt.Errorf("cannot mix interval qualities; expected %q, but got %q at pos %d", text[1], text[i], i)
			}
		}
		return augmentedErr(diatonic-1, size)
	case 'd':
		size := 1
		for i := 2; i < len(text); i++ {
			if text[i] == text[1] {
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

func Steps(intervals []Interval) [12]bool {
	var r [12]bool

	for _, iv := range intervals {
		r[iv.Chromatic()] = true
	}

	return r
}

func Sort(intervals []Interval) {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].String() < intervals[j].String()
	})
}

type Interval struct {
	diatonic int
	quality  Quality
}

func (i Interval) Chromatic() int {
	return chromaticFromDiatonicAndQuality(i.diatonic, i.quality)
}

func (i Interval) Diatonic() int {
	return i.diatonic
}

func (i Interval) MarshalJSON() ([]byte, error) {
	return []byte("\"" + i.String() + "\""), nil
}

func (i Interval) Quality() Quality {
	return i.quality
}

func (i Interval) String() string {
	return fmt.Sprintf("%d%s", i.diatonic+1, i.quality)
}

func (i Interval) Transpose(other Interval) Interval {
	newDiatonic := normalizeDiatonic(i.diatonic + other.diatonic)
	ic := i.Chromatic()
	oc := other.Chromatic()
	newChromatic := normalizeChromatic(ic + oc) //i.Chromatic() + other.Chromatic())
	return New(newDiatonic, qualityFromDiatonicAndChromatic(newDiatonic, newChromatic))
}

var (
	diatonicToChromatic = [7]int{0, 2, 4, 5, 7, 9, 11}
)

func chromaticFromDiatonicAndQuality(diatonic int, q Quality) int {
	chromatic := diatonicToChromatic[diatonic]

	switch q.Kind() {
	case QualityKindPerfect, QualityKindMajor:
		return chromatic
	case QualityKindMinor:
		return chromatic - 1
	case QualityKindDiminished:
		value := chromatic - q.Size()
		if diatonic == 1 || diatonic == 2 || diatonic == 5 || diatonic == 6 {
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
	if diatonic != 1 && diatonic != 2 && diatonic != 5 && diatonic != 6 {
		return Interval{}, fmt.Errorf("only 1, 2, 5, and 6 can be major, but got %d", diatonic)
	}

	return New(diatonic, MajorQuality()), nil
}

func minorErr(diatonic int) (Interval, error) {
	if diatonic != 1 && diatonic != 2 && diatonic != 5 && diatonic != 6 {
		return Interval{}, fmt.Errorf("only 1, 2, 5, and 6 can be minor, but got %d", diatonic)
	}

	return New(diatonic, MinorQuality()), nil
}

func normalizeDiatonic(diatonic int) int {
	return (diatonic + 7) % 7
}

func normalizeChromatic(chromatic int) int {
	return (chromatic + 12) % 12
}

func perfectErr(diatonic int) (Interval, error) {
	if diatonic != 0 && diatonic != 3 && diatonic != 4 {
		return Interval{}, fmt.Errorf("only 0, 3, and 4 can be perfect, but got %d", diatonic)
	}

	return New(diatonic, PerfectQuality()), nil
}

func qualityFromDiatonicAndChromatic(diatonic, chromatic int) Quality {
	diff := chromatic - diatonicToChromatic[diatonic]

	switch diatonic {
	case 0:
		fallthrough
	case 3:
		fallthrough
	case 4:
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
	}

	if diff > 6 {
		diff -= 12
	} else if diff < -6 {
		diff += 12
	}

	if diff > 0 {
		q, _ := AugmentedQuality(diff)
		return q
	}

	q, _ := DiminishedQuality(-diff)
	return q
}
