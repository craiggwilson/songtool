package interval

import "fmt"

type Quality int

const (
	QualityPerfect Quality = iota
	QualityMajor
	QualityMinor
	QualityAugmented
	QualityDiminished
)

func (q Quality) String() string {
	switch q {
	case QualityPerfect:
		return "P"
	case QualityMajor:
		return "M"
	case QualityMinor:
		return "m"
	case QualityAugmented:
		return "a"
	case QualityDiminished:
		return "d"
	default:
		panic(fmt.Sprintf("unsupported quality %d", q))
	}
}
