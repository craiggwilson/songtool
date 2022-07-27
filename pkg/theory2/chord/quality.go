package chord

type Quality int

const (
	QualityIndeterminate Quality = iota
	QualityMajor
	QualityMinor
	QualityDiminished
	QualityAugmented
)
