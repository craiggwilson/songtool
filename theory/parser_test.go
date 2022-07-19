package theory_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/craiggwilson/songtools/theory"
	"github.com/craiggwilson/songtools/theory/note"
)

func TestParser_ParseNote(t *testing.T) {
	testCases := []struct {
		name        string
		expected    note.Note
		expectedErr error
	}{
		{
			name: "A",
			expected: note.Note{
				Name:        "A",
				DegreeClass: 0,
				PitchClass:  0,
				Accidentals: 0,
			},
		},
		{
			name: "A#",
			expected: note.Note{
				Name:        "A#",
				DegreeClass: 0,
				PitchClass:  0,
				Accidentals: 1,
			},
		},
		{
			name: "Abb",
			expected: note.Note{
				Name:        "Abb",
				DegreeClass: 0,
				PitchClass:  0,
				Accidentals: -2,
			},
		},
		{
			name: "G",
			expected: note.Note{
				Name:        "G",
				DegreeClass: 6,
				PitchClass:  11,
				Accidentals: 0,
			},
		},
		{
			name: "G#",
			expected: note.Note{
				Name:        "G#",
				DegreeClass: 6,
				PitchClass:  11,
				Accidentals: 1,
			},
		},
		{
			name: "Gbb",
			expected: note.Note{
				Name:        "Gbb",
				DegreeClass: 6,
				PitchClass:  11,
				Accidentals: -2,
			},
		},
	}

	p := theory.NewParser(theory.DefaultConfig())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := p.ParseNote(tc.name)
			require.ErrorIs(t, tc.expectedErr, err)

			require.Equal(t, tc.expected, actual)
		})
	}
}
