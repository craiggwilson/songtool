package theory_test

import (
	"testing"

	"github.com/craiggwilson/songtools/theory"
	"github.com/stretchr/testify/require"
)

func TestParseChord(t *testing.T) {
	testCases := []struct {
		name     string
		expected theory.Chord
	}{
		{
			name: "A",
			expected: theory.Chord{
				Root: theory.Note{
					Name:        "A",
					DegreeClass: 5,
					PitchClass:  9,
					Accidentals: 0,
				},
				Intervals: []int{1, 4, 7},
			},
		},
		{
			name: "Amaj9",
			expected: theory.Chord{
				Root: theory.Note{
					Name:        "A",
					DegreeClass: 5,
					PitchClass:  9,
					Accidentals: 0,
				},
				Intervals: []int{1, 4, 7, 11, 14},
				Suffix:    "maj9",
			},
		},
		{
			name: "Gbbmsus2add6",
			expected: theory.Chord{
				Root: theory.Note{
					Name:        "Gbb",
					DegreeClass: 4,
					PitchClass:  5,
					Accidentals: -2,
				},
				Intervals: []int{1, 2, 7, 9},
				Suffix:    "msus2add6",
			},
		},
		{
			name: "C#m7/G#",
			expected: theory.Chord{
				Root: theory.Note{
					Name:        "C#",
					DegreeClass: 0,
					PitchClass:  1,
					Accidentals: 1,
				},
				Intervals: []int{1, 3, 7, 10},
				Suffix:    "m7",
				Base: theory.Note{
					Name:        "G#",
					DegreeClass: 4,
					PitchClass:  8,
					Accidentals: 1,
				},
			},
		},
	}

	cfg := theory.DefaultConfig()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := theory.ParseChord(&cfg, tc.name)
			require.Nil(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}
