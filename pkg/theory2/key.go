package theory2

import "github.com/craiggwilson/songtool/pkg/theory2/key"

type KeyNamer interface {
	NameKey(key.Key) string
}

func NameKey(k key.Key) string {
	return std.NameKey(k)
}

type KeyParser interface {
	ParseKey(string) (key.Key, error)
}

func ParseKey(text string) (key.Key, error) {
	return std.ParseKey(text)
}
