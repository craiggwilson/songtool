package theory_test

import (
	"testing"

	"github.com/craiggwilson/songtool/pkg/theory"
	"github.com/stretchr/testify/require"
)

func TestGenerateScale(t *testing.T) {
	testCases := []struct {
		name      string
		root      theory.Note
		intervals []theory.Interval
		expected  theory.Scale
	}{
		{
			name:      "C Major",
			root:      theory.MustNote(theory.ParseNote("C")),
			intervals: theory.IonianIntervals,
			expected: theory.Scale{
				Name: "C Major",
				Notes: []theory.Note{
					{
						Name:        "C",
						DegreeClass: 0,
						PitchClass:  0,
						Accidentals: 0,
					},
					{
						Name:        "D",
						DegreeClass: 1,
						PitchClass:  2,
						Accidentals: 0,
					},
					{
						Name:        "E",
						DegreeClass: 2,
						PitchClass:  4,
						Accidentals: 0,
					},
					{
						Name:        "F",
						DegreeClass: 3,
						PitchClass:  5,
						Accidentals: 0,
					},
					{
						Name:        "G",
						DegreeClass: 4,
						PitchClass:  7,
						Accidentals: 0,
					},
					{
						Name:        "A",
						DegreeClass: 5,
						PitchClass:  9,
						Accidentals: 0,
					},
					{
						Name:        "B",
						DegreeClass: 6,
						PitchClass:  11,
						Accidentals: 0,
					},
				},
			},
		},
		{
			name:      "C# Major",
			root:      theory.MustNote(theory.ParseNote("C#")),
			intervals: theory.IonianIntervals,
			expected: theory.Scale{
				Name: "C# Major",
				Notes: []theory.Note{
					{
						Name:        "C#",
						DegreeClass: 0,
						PitchClass:  1,
						Accidentals: 1,
					},
					{
						Name:        "D#",
						DegreeClass: 1,
						PitchClass:  3,
						Accidentals: 1,
					},
					{
						Name:        "E#",
						DegreeClass: 2,
						PitchClass:  5,
						Accidentals: 1,
					},
					{
						Name:        "F#",
						DegreeClass: 3,
						PitchClass:  6,
						Accidentals: 1,
					},
					{
						Name:        "G#",
						DegreeClass: 4,
						PitchClass:  8,
						Accidentals: 1,
					},
					{
						Name:        "A#",
						DegreeClass: 5,
						PitchClass:  10,
						Accidentals: 1,
					},
					{
						Name:        "B#",
						DegreeClass: 6,
						PitchClass:  0,
						Accidentals: 1,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := theory.GenerateDiatonicScale(tc.name, tc.root, tc.intervals)

			require.Equal(t, tc.expected, actual)
		})
	}
}
