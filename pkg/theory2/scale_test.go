package theory2_test

import (
	"testing"

	"github.com/craiggwilson/songtool/pkg/theory2"
	"github.com/craiggwilson/songtool/pkg/theory2/interval"
	"github.com/craiggwilson/songtool/pkg/theory2/note"
	"github.com/stretchr/testify/require"
)

func TestGenerateScale(t *testing.T) {
	testCases := []struct {
		name      string
		root      note.Note
		intervals []interval.Interval
		expected  theory2.Scale
	}{
		{
			name:      "C Major",
			root:      note.C,
			intervals: interval.Scales.Ionian,
			expected: theory2.NewScale("C Major",
				note.C,
				note.D,
				note.E,
				note.F,
				note.G,
				note.A,
				note.B,
			),
		},
		{
			name:      "C# Major",
			root:      note.CSharp,
			intervals: interval.Scales.Ionian,
			expected: theory2.NewScale("C# Major",
				note.CSharp,
				note.DSharp,
				note.ESharp,
				note.FSharp,
				note.GSharp,
				note.ASharp,
				note.BSharp,
			),
		},
		{
			name:      "C# Chromatic",
			root:      note.CSharp,
			intervals: interval.Scales.Chromatic,
			expected: theory2.NewScale("C# Chromatic",
				note.CSharp,
				note.D,
				note.DSharp,
				note.E,
				note.ESharp,
				note.FSharp,
				note.G,
				note.GSharp,
				note.A,
				note.ASharp,
				note.B,
				note.BSharp,
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := theory2.GenerateScale(tc.name, tc.root, tc.intervals...)
			require.Equal(t, tc.expected, actual)
		})
	}
}
