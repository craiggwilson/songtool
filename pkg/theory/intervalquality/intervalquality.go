package intervalquality

import (
	"fmt"
	"strings"
)

func Perfect() IntervalQuality {
	return IntervalQuality{PerfectKind, 0}
}

func Major() IntervalQuality {
	return IntervalQuality{MajorKind, 0}
}

func Minor() IntervalQuality {
	return IntervalQuality{MinorKind, 0}
}

func Augmented(size int) IntervalQuality {
	return IntervalQuality{AugmentedKind, size - 1}
}

func Diminished(size int) IntervalQuality {
	return IntervalQuality{DiminishedKind, size - 1}
}

// The quality of an interval. The zero value is P1.
type IntervalQuality struct {
	Kind Kind
	Size int
}

func (q IntervalQuality) String() string {
	return strings.Repeat(q.Kind.String(), q.Size+1)
}

type Kind int

const (
	PerfectKind Kind = iota
	MajorKind
	MinorKind
	AugmentedKind
	DiminishedKind
)

func (k Kind) String() string {
	switch k {
	case PerfectKind:
		return "P"
	case MajorKind:
		return "M"
	case MinorKind:
		return "m"
	case AugmentedKind:
		return "a"
	case DiminishedKind:
		return "d"
	default:
		panic(fmt.Sprintf("unsupported quality %d", k))
	}
}
