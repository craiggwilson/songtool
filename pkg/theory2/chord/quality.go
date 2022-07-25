package chord

type Quality int

const (
	IndeterminateQuality Quality = iota
	MajorQuality
	MinorQuality
	DiminishedQuality
	AugmentedQuality
)
