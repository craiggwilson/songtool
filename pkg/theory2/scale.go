package theory2

import (
	"sort"

	"github.com/craiggwilson/songtool/pkg/theory2/interval"
	"github.com/craiggwilson/songtool/pkg/theory2/scale"
)

type ScaleMeta struct {
	Name      string              `json:"name"`
	Intervals []interval.Interval `json:"intervals"`
}

type ScaleParser interface {
	ParseScale(string) (scale.Scale, error)
}

func ParseScale(name string) (scale.Scale, error) {
	return std.ParseScale(name)
}

type ScaleProvider interface {
	LookupScale(string) (ScaleMeta, bool)
	ListScales() []ScaleMeta
}

func ListScales() []ScaleMeta {
	return std.ListScales()
}

func LookupScale(name string) (ScaleMeta, bool) {
	return std.LookupScale(name)
}

func SortScaleMetas(scaleMetas []ScaleMeta) {
	sort.Slice(scaleMetas, func(i, j int) bool {
		return scaleMetas[i].Name < scaleMetas[j].Name
	})
}
