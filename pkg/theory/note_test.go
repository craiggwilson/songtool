package theory_test

import (
	"testing"

	theory2 "github.com/craiggwilson/songtool/pkg/theory"
	"github.com/craiggwilson/songtool/pkg/theory/note"
	"github.com/stretchr/testify/require"
)

func TestNameNote(t *testing.T) {
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
			actual := theory2.NameNote(tc.note)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestParseNote(t *testing.T) {
	testCases := []struct {
		name           string
		expected       note.Note
		expectedErrMsg string
	}{
		{
			name:     "A",
			expected: note.A,
		},
		{
			name:     "A#",
			expected: note.ASharp,
		},
		{
			name:     "Abb",
			expected: note.New(5, 7),
		},
		{
			name:     "C",
			expected: note.C,
		},
		{
			name:     "C#",
			expected: note.CSharp,
		},
		{
			name:     "Cbb",
			expected: note.New(0, 10),
		},
		{
			name:           "Ab#",
			expected:       note.Note{},
			expectedErrMsg: `expected EOF at position 2, but had "#"`,
		},
		{
			name:           "A#b",
			expected:       note.Note{},
			expectedErrMsg: `expected EOF at position 2, but had "b"`,
		},
		{
			name:           "H",
			expected:       note.Note{},
			expectedErrMsg: `expected natural note name at position 0: expected one of ["C" "D" "E" "F" "G" "A" "B"], but got "H"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := theory2.ParseNote(tc.name)
			if len(tc.expectedErrMsg) > 0 {
				require.EqualError(t, err, tc.expectedErrMsg)
			} else {
				require.Nil(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}
