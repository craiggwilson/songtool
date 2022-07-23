package theory_test

import (
	"testing"

	"github.com/craiggwilson/songtool/pkg/theory"
	"github.com/stretchr/testify/require"
)

func TestGenerateKeys(t *testing.T) {
	testCases := []struct {
		kind     theory.KeyKind
		expected []theory.Key
	}{
		{
			kind: theory.KeyMajor,
			expected: []theory.Key{
				{
					Note: theory.Note{
						Name:        "Cb",
						DegreeClass: 0,
						PitchClass:  11,
						Accidentals: -1,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "C",
						DegreeClass: 0,
						PitchClass:  0,
						Accidentals: 0,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "C#",
						DegreeClass: 0,
						PitchClass:  1,
						Accidentals: 1,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "Db",
						DegreeClass: 1,
						PitchClass:  1,
						Accidentals: -1,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "D",
						DegreeClass: 1,
						PitchClass:  2,
						Accidentals: 0,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "D#",
						DegreeClass: 1,
						PitchClass:  3,
						Accidentals: 1,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "Eb",
						DegreeClass: 2,
						PitchClass:  3,
						Accidentals: -1,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "E",
						DegreeClass: 2,
						PitchClass:  4,
						Accidentals: 0,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "E#",
						DegreeClass: 2,
						PitchClass:  5,
						Accidentals: 1,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "Fb",
						DegreeClass: 3,
						PitchClass:  4,
						Accidentals: -1,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "F",
						DegreeClass: 3,
						PitchClass:  5,
						Accidentals: 0,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "F#",
						DegreeClass: 3,
						PitchClass:  6,
						Accidentals: 1,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "Gb",
						DegreeClass: 4,
						PitchClass:  6,
						Accidentals: -1,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "G",
						DegreeClass: 4,
						PitchClass:  7,
						Accidentals: 0,
					},
					Kind: theory.KeyMajor,
				}, {
					Note: theory.Note{
						Name:        "G#",
						DegreeClass: 4,
						PitchClass:  8,
						Accidentals: 1,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "Ab",
						DegreeClass: 5,
						PitchClass:  8,
						Accidentals: -1,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "A",
						DegreeClass: 5,
						PitchClass:  9,
						Accidentals: 0,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "A#",
						DegreeClass: 5,
						PitchClass:  10,
						Accidentals: 1,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "Bb",
						DegreeClass: 6,
						PitchClass:  10,
						Accidentals: -1,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "B",
						DegreeClass: 6,
						PitchClass:  11,
						Accidentals: 0,
					},
					Kind: theory.KeyMajor,
				},
				{
					Note: theory.Note{
						Name:        "B#",
						DegreeClass: 6,
						PitchClass:  0,
						Accidentals: 1,
					},
					Kind: theory.KeyMajor,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.kind), func(t *testing.T) {
			actual := theory.GenerateKeys(tc.kind)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestParseKey(t *testing.T) {
	testCases := []struct {
		name     string
		expected theory.Key
	}{
		{
			name: "A",
			expected: theory.Key{
				Note: theory.Note{
					Name:        "A",
					DegreeClass: 5,
					PitchClass:  9,
					Accidentals: 0,
				},
				Kind: theory.KeyMajor,
			},
		},
		{
			name: "A#",
			expected: theory.Key{
				Note: theory.Note{
					Name:        "A#",
					DegreeClass: 5,
					PitchClass:  10,
					Accidentals: 1,
				},
				Kind: theory.KeyMajor,
			},
		},
		{
			name: "Abb",
			expected: theory.Key{
				Note: theory.Note{
					Name:        "Abb",
					DegreeClass: 5,
					PitchClass:  7,
					Accidentals: -2,
				},
				Kind: theory.KeyMajor,
			},
		},
		{
			name: "Am",
			expected: theory.Key{
				Note: theory.Note{
					Name:        "A",
					DegreeClass: 5,
					PitchClass:  9,
					Accidentals: 0,
				},
				Suffix: "m",
				Kind:   theory.KeyMinor,
			},
		},
		{
			name: "A#m",
			expected: theory.Key{
				Note: theory.Note{
					Name:        "A#",
					DegreeClass: 5,
					PitchClass:  10,
					Accidentals: 1,
				},
				Suffix: "m",
				Kind:   theory.KeyMinor,
			},
		},
		{
			name: "Abbm",
			expected: theory.Key{
				Note: theory.Note{
					Name:        "Abb",
					DegreeClass: 5,
					PitchClass:  7,
					Accidentals: -2,
				},
				Suffix: "m",
				Kind:   theory.KeyMinor,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := theory.ParseKey(tc.name)
			require.Nil(t, err)

			require.Equal(t, tc.expected, actual)
		})
	}
}
