package note_test

import (
	"fmt"
	"testing"

	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/craiggwilson/songtool/pkg/theory/note"
	"github.com/stretchr/testify/require"
)

func TestNote_Enharmonic(t *testing.T) {
	testCases := []struct {
		note     note.Note
		expected interval.Interval
	}{
		{
			note:     note.A,
			expected: interval.NewWithChromatic(0, 0),
		},
		{
			note:     note.ASharp,
			expected: interval.NewWithChromatic(1, 0),
		},
		{
			note:     note.ESharp,
			expected: interval.NewWithChromatic(1, 0),
		},
		{
			note:     note.FFlat,
			expected: interval.NewWithChromatic(-1, 0),
		},
		{
			note:     note.New(note.F.DegreeClass(), note.F.PitchClass()+2),
			expected: interval.NewWithChromatic(1, 0),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.note), func(t *testing.T) {
			actual := tc.note.Enharmonic()
			require.Equal(t, tc.expected, actual, "expected %s, but got %s", tc.expected, actual)
		})
	}
}

func TestNote_Step(t *testing.T) {
	testCases := []struct {
		note     note.Note
		step     int
		expected interval.Interval
	}{
		{
			note:     note.C,
			step:     1,
			expected: interval.NewWithChromatic(0, 1), // To C#
		},
		{
			note:     note.CSharp,
			step:     1,
			expected: interval.NewWithChromatic(1, 1), // To D
		},
		{
			note:     note.DFlat,
			step:     1,
			expected: interval.NewWithChromatic(0, 1), // To D
		},
		{
			note:     note.B,
			step:     1,
			expected: interval.NewWithChromatic(1, 1), // To C
		},
		{
			note:     note.BFlat,
			step:     1,
			expected: interval.NewWithChromatic(0, 1), // To B
		},
		{
			note:     note.C,
			step:     2,
			expected: interval.NewWithChromatic(1, 2), // To D
		},
		{
			note:     note.C,
			step:     3,
			expected: interval.NewWithChromatic(1, 3), // To D#
		},
		{
			note:     note.C,
			step:     4,
			expected: interval.NewWithChromatic(2, 4), // To E
		},
		{
			note:     note.C,
			step:     -1,
			expected: interval.NewWithChromatic(-1, -1), // To B
		},
		{
			note:     note.CSharp,
			step:     -1,
			expected: interval.NewWithChromatic(0, -1), // To C
		},
		{
			note:     note.DFlat,
			step:     -1,
			expected: interval.NewWithChromatic(-1, -1), // To C
		},
		{
			note:     note.C,
			step:     -2,
			expected: interval.NewWithChromatic(-1, -2), // To Bb
		},
		{
			note:     note.C,
			step:     -3,
			expected: interval.NewWithChromatic(-2, -3), // To A
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v step %d", tc.note, tc.step), func(t *testing.T) {
			actual := tc.note.Step(tc.step)
			require.Equal(t, tc.expected, actual, "expected %s, but got %s", tc.expected, actual)
		})
	}
}
