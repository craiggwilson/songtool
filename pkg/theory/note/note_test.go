package note_test

import (
	"fmt"
	"testing"

	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/craiggwilson/songtool/pkg/theory/note"
	"github.com/stretchr/testify/require"
)

func TestNote_Step(t *testing.T) {
	testCases := []struct {
		note     note.Note
		step     int
		expected interval.Interval
	}{
		{
			note:     note.A,
			step:     1,
			expected: interval.NewWithChromatic(0, 1),
		},
		{
			note:     note.ASharp,
			step:     1,
			expected: interval.NewWithChromatic(1, 1),
		},
		{
			note:     note.A,
			step:     2,
			expected: interval.NewWithChromatic(1, 2),
		},
		{
			note:     note.A,
			step:     3,
			expected: interval.NewWithChromatic(2, 3),
		},
		{
			note:     note.A,
			step:     -1,
			expected: interval.NewWithChromatic(0, -1),
		},
		{
			note:     note.ASharp,
			step:     -1,
			expected: interval.NewWithChromatic(0, -1),
		},
		{
			note:     note.A,
			step:     -2,
			expected: interval.NewWithChromatic(-1, -2),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v step %d", tc.note, tc.step), func(t *testing.T) {
			actual := tc.note.Step(tc.step)
			require.Equal(t, tc.expected, actual)
		})
	}
}
