package interval

import (
	"fmt"
	"strings"
)

func AugmentedQuality(size int) (Quality, error) {
	if size < 0 {
		return Quality{}, fmt.Errorf("augmented quality sizes must be > 1, but was %d", size)
	}

	return Quality{QualityKindAugmented, size - 1}, nil
}

func DiminishedQuality(size int) (Quality, error) {
	if size < 1 {
		return Quality{}, fmt.Errorf("dimished quality sizes must be > 1, but was %d", size)
	}

	return Quality{QualityKindDiminished, size - 1}, nil
}

func MajorQuality() Quality {
	return Quality{QualityKindMajor, 0}
}

func MinorQuality() Quality {
	return Quality{QualityKindMinor, 0}
}

func PerfectQuality() Quality {
	return Quality{QualityKindPerfect, 0}
}

// The quality of an interval. The zero value is P1.
type Quality struct {
	kind QualityKind
	size int
}

func (q Quality) Kind() QualityKind {
	return q.kind
}

func (q Quality) Size() int {
	return q.size + 1
}

func (q Quality) String() string {
	return strings.Repeat(q.kind.String(), q.size+1)
}

type QualityKind int

const (
	QualityKindPerfect QualityKind = iota
	QualityKindMajor
	QualityKindMinor
	QualityKindAugmented
	QualityKindDiminished
)

func (k QualityKind) String() string {
	switch k {
	case QualityKindPerfect:
		return "P"
	case QualityKindMajor:
		return "M"
	case QualityKindMinor:
		return "m"
	case QualityKindAugmented:
		return "a"
	case QualityKindDiminished:
		return "d"
	default:
		panic(fmt.Sprintf("unsupported quality %d", k))
	}
}
