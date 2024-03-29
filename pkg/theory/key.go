package theory

import (
	"github.com/craiggwilson/songtool/pkg/theory/key"
)

func NameKey(k key.Key) string {
	return std.NameKey(k)
}

func ParseKey(text string) (key.Named, error) {
	return std.ParseKey(text)
}
