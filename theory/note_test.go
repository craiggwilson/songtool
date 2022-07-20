package theory_test

import (
	"testing"

	"github.com/craiggwilson/songtools/theory"
	"github.com/stretchr/testify/require"
)

func TestParseNote(t *testing.T) {
	testCases := []struct {
		name        string
		expected    theory.Note
		expectedErr error
	}{
		{
			name: "A",
			expected: theory.Note{
				Name:        "A",
				DegreeClass: 5,
				PitchClass:  9,
				Accidentals: 0,
			},
		},
		{
			name: "A#",
			expected: theory.Note{
				Name:        "A#",
				DegreeClass: 5,
				PitchClass:  10,
				Accidentals: 1,
			},
		},
		{
			name: "Abb",
			expected: theory.Note{
				Name:        "Abb",
				DegreeClass: 5,
				PitchClass:  7,
				Accidentals: -2,
			},
		},
		{
			name: "C",
			expected: theory.Note{
				Name:        "C",
				DegreeClass: 0,
				PitchClass:  0,
				Accidentals: 0,
			},
		},
		{
			name: "C#",
			expected: theory.Note{
				Name:        "C#",
				DegreeClass: 0,
				PitchClass:  1,
				Accidentals: 1,
			},
		},
		{
			name: "Cbb",
			expected: theory.Note{
				Name:        "Cbb",
				DegreeClass: 0,
				PitchClass:  10,
				Accidentals: -2,
			},
		},
	}

	cfg := theory.DefaultConfig()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := theory.ParseNote(&cfg, tc.name)
			require.ErrorIs(t, err, tc.expectedErr)

			require.Equal(t, tc.expected, actual)
		})
	}
}
