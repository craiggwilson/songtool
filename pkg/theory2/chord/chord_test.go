package chord_test

import (
	"testing"

	"github.com/craiggwilson/songtool/pkg/theory2/chord"
	"github.com/craiggwilson/songtool/pkg/theory2/interval"
	"github.com/stretchr/testify/require"
)

func TestChord_Quality(t *testing.T) {
	testCases := []struct {
		name     string
		chord    chord.Chord
		expected chord.Quality
	}{
		{
			name:     "no intervals",
			chord:    chord.Chord{},
			expected: chord.IndeterminateQuality,
		},
		{
			name: "major triad",
			chord: chord.New(
				interval.Perfect(0),
				interval.Major(2),
				interval.Perfect(4),
			),
			expected: chord.MajorQuality,
		},
		{
			name: "augmented triad",
			chord: chord.New(
				interval.Perfect(0),
				interval.Major(2),
				interval.Augmented(4, 1),
			),
			expected: chord.AugmentedQuality,
		},
		{
			name: "minor triad",
			chord: chord.New(
				interval.Perfect(0),
				interval.Minor(2),
				interval.Perfect(4),
			),
			expected: chord.MinorQuality,
		},
		{
			name: "diminished triad",
			chord: chord.New(
				interval.Perfect(0),
				interval.Minor(2),
				interval.Diminished(4, 1),
			),
			expected: chord.DiminishedQuality,
		},
		{
			name: "no 3rd",
			chord: chord.New(
				interval.Perfect(0),
				interval.Perfect(4),
			),
			expected: chord.IndeterminateQuality,
		},
		{
			name: "major 7th",
			chord: chord.New(
				interval.Perfect(0),
				interval.Major(2),
				interval.Perfect(4),
				interval.Major(6),
			),
			expected: chord.MajorQuality,
		},
		{
			name: "minor 7th",
			chord: chord.New(
				interval.Perfect(0),
				interval.Minor(2),
				interval.Perfect(4),
				interval.Minor(6),
			),
			expected: chord.MinorQuality,
		},
		{
			name: "7th",
			chord: chord.New(
				interval.Perfect(0),
				interval.Major(2),
				interval.Perfect(4),
				interval.Minor(6),
			),
			expected: chord.MajorQuality,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.chord.Quality()
			require.Equal(t, tc.expected, actual)
		})
	}

}
