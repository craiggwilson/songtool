package eval

import (
	"github.com/craiggwilson/songtool/pkg/songio"
	"github.com/craiggwilson/songtool/pkg/theory"
)

type Context struct {
	Theory *theory.Theory
	Meta   *songio.Meta
	Lines  []songio.Line
}
