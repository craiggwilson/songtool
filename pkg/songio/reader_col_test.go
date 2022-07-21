package songio_test

import (
	"strings"
	"testing"

	"github.com/craiggwilson/songtool/pkg/songio"
	"github.com/stretchr/testify/require"
)

func TestReadChordsOverLyrics(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{
			input: "#title=Hello World\n#key=G\n\n[Verse 1]\nG  C  D\nRunning to say Hello\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			lines, err := songio.ReadAllLines(songio.ReadChordsOverLyrics(nil, strings.NewReader(tc.input)))
			require.Nil(t, err)

			require.Equal(t, 6, len(lines))
		})
	}

}
