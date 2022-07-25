package theory2_test

import (
	"testing"

	"github.com/craiggwilson/songtool/pkg/theory2"
	"github.com/craiggwilson/songtool/pkg/theory2/note"
	"github.com/stretchr/testify/require"
)

func TestConfig_NameNote(t *testing.T) {
	cfg := theory2.DefaultConfig()

	testCases := []struct {
		note     note.Note
		expected string
	}{
		{
			note:     note.CFlat,
			expected: "Cb",
		},
		{
			note:     note.C,
			expected: "C",
		},
		{
			note:     note.CSharp,
			expected: "C#",
		},
		{
			note:     note.EFlat,
			expected: "Eb",
		},
		{
			note:     note.E,
			expected: "E",
		},
		{
			note:     note.ESharp,
			expected: "E#",
		},
		{
			note:     note.FFlat,
			expected: "Fb",
		},
		{
			note:     note.F,
			expected: "F",
		},
		{
			note:     note.FSharp,
			expected: "F#",
		},
		{
			note:     note.BFlat,
			expected: "Bb",
		},
		{
			note:     note.B,
			expected: "B",
		},
		{
			note:     note.BSharp,
			expected: "B#",
		},
		{
			note:     note.New(3, 3),
			expected: "Fbb",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			actual := cfg.NameNote(tc.note)
			require.Equal(t, tc.expected, actual)
		})
	}
}
