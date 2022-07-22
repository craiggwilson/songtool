package songio

import (
	"io"

	"github.com/craiggwilson/songtool/pkg/theory"
)

type Reader interface {
	ReadSong(*theory.Config, io.Reader) (Song, error)
}

type Writer interface {
	WriteSong(theory.Config, Song, io.Writer) (int, error)
}
