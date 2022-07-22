package songio

import (
	"io"

	"github.com/craiggwilson/songtool/pkg/theory"
)

type ChordsOverLyricsFormat struct{}

func (ChordsOverLyricsFormat) Read(cfg *theory.Config, src io.Reader) (Song, error) {
	return ReadChordsOverLyrics(cfg, src), nil
}

func (ChordsOverLyricsFormat) Write(_ *theory.Config, src Song, dst io.Writer) (int, error) {
	return WriteChordsOverLyrics(src, dst)
}
