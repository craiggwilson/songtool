package theory_test

import (
	"fmt"
	"testing"

	"github.com/craiggwilson/songtool/pkg/theory"
	"github.com/stretchr/testify/require"
)

func TestParseNote(t *testing.T) {
	testCases := []struct {
		name           string
		expected       theory.Note
		expectedErrMsg string
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
		{
			name:           "Ab#",
			expected:       theory.Note{},
			expectedErrMsg: `expected EOF at position 2, but had "#"`,
		},
		{
			name:           "A#b",
			expected:       theory.Note{},
			expectedErrMsg: `expected EOF at position 2, but had "b"`,
		},
		{
			name:           "H",
			expected:       theory.Note{},
			expectedErrMsg: `expected natural note name at position 0: expected one of ['C' 'D' 'E' 'F' 'G' 'A' 'B'], but got 'H'`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := theory.ParseNote(nil, tc.name)
			if len(tc.expectedErrMsg) > 0 {
				require.EqualError(t, err, tc.expectedErrMsg)
			}
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestTransposeNote(t *testing.T) {
	testCases := []struct {
		from               theory.Note
		pitchClassInterval int
		enharmonic         theory.Enharmonic
		expected           theory.Note
	}{
		{
			from:               mustParseNote(nil, "C"),
			pitchClassInterval: 2,
			enharmonic:         theory.EnharmonicSharp,
			expected:           mustParseNote(nil, "D"),
		},
		{
			from:               mustParseNote(nil, "C"),
			pitchClassInterval: 1,
			enharmonic:         theory.EnharmonicSharp,
			expected:           mustParseNote(nil, "C#"),
		},
		{
			from:               mustParseNote(nil, "C"),
			pitchClassInterval: 1,
			enharmonic:         theory.EnharmonicFlat,
			expected:           mustParseNote(nil, "Db"),
		},
		{
			from:               mustParseNote(nil, "C#"),
			pitchClassInterval: 0,
			enharmonic:         theory.EnharmonicFlat,
			expected:           mustParseNote(nil, "Db"),
		},
		{
			from:               mustParseNote(nil, "Db"),
			pitchClassInterval: 0,
			enharmonic:         theory.EnharmonicSharp,
			expected:           mustParseNote(nil, "C#"),
		},
		{
			from:               mustParseNote(nil, "C"),
			pitchClassInterval: -1,
			enharmonic:         theory.EnharmonicSharp,
			expected:           mustParseNote(nil, "B"),
		},
		{
			from:               mustParseNote(nil, "C"),
			pitchClassInterval: -1,
			enharmonic:         theory.EnharmonicFlat,
			expected:           mustParseNote(nil, "B"),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s +- (%d, %d)", tc.from.Name, tc.pitchClassInterval, tc.enharmonic), func(t *testing.T) {
			actual := theory.TransposeNote(nil, tc.from, theory.IntervalFromStep(nil, tc.from, tc.pitchClassInterval, tc.enharmonic))
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestTransposeNoteDirect(t *testing.T) {
	testCases := []struct {
		from                theory.Note
		degreeClassInterval int
		pitchClassInterval  int
		expected            theory.Note
	}{
		{
			from:                mustParseNote(nil, "C"),
			degreeClassInterval: 1,
			pitchClassInterval:  2,
			expected:            mustParseNote(nil, "D"),
		},
		{
			from:                mustParseNote(nil, "C"),
			degreeClassInterval: 1,
			pitchClassInterval:  1,
			expected:            mustParseNote(nil, "Db"),
		},
		{
			from:                mustParseNote(nil, "C"),
			degreeClassInterval: 0,
			pitchClassInterval:  1,
			expected:            mustParseNote(nil, "C#"),
		},
		{
			from:                mustParseNote(nil, "C"),
			degreeClassInterval: 0,
			pitchClassInterval:  -1,
			expected:            mustParseNote(nil, "Cb"),
		},
		{
			from:                mustParseNote(nil, "C"),
			degreeClassInterval: -1,
			pitchClassInterval:  -1,
			expected:            mustParseNote(nil, "B"),
		},
		{
			from:                mustParseNote(nil, "C"),
			degreeClassInterval: -1,
			pitchClassInterval:  -2,
			expected:            mustParseNote(nil, "Bb"),
		},
		{
			from:                mustParseNote(nil, "C"),
			degreeClassInterval: -2,
			pitchClassInterval:  -2,
			expected:            mustParseNote(nil, "A#"),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s +- (%d, %d)", tc.from.Name, tc.degreeClassInterval, tc.pitchClassInterval), func(t *testing.T) {
			actual := theory.TransposeNote(nil, tc.from, theory.Interval{tc.degreeClassInterval, tc.pitchClassInterval})
			require.Equal(t, tc.expected, actual)
		})
	}
}

func mustParseNote(cfg *theory.Config, text string) theory.Note {
	n, err := theory.ParseNote(cfg, text)
	if err != nil {
		panic(err)
	}

	return n
}
