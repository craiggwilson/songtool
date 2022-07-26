package theory2

import "github.com/craiggwilson/songtool/pkg/theory2/scale"

type ScaleParser interface {
	ParseScale(name string) (scale.Scale, error)
}

func ParseScale(name string) (scale.Scale, error) {
	return std.ParseScale(name)
}
