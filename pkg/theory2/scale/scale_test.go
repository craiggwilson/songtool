package scale_test

import (
	"testing"

	"github.com/craiggwilson/songtool/pkg/theory2/interval"
	"github.com/craiggwilson/songtool/pkg/theory2/note"
	"github.com/craiggwilson/songtool/pkg/theory2/scale"
	"github.com/stretchr/testify/require"
)

func TestGenerateScale(t *testing.T) {
	testCases := []struct {
		name      string
		root      note.Note
		intervals []interval.Interval
		expected  scale.Scale
	}{
		{
			name:      "C Major",
			root:      note.C,
			intervals: interval.Scales.Ionian,
			expected: scale.New("C Major",
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
			expected: scale.New("C# Major",
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
			name:      "C Chromatic",
			root:      note.C,
			intervals: interval.Scales.Chromatic,
			expected: scale.New("C Chromatic",
				note.C,
				note.DFlat, //note.CSharp,
				note.D,
				note.EFlat, //note.DSharp,
				note.E,
				note.F,
				note.GFlat, //note.FSharp,
				note.G,
				note.AFlat, // note.GSharp,
				note.A,
				note.BFlat, //note.ASharp,
				note.B,
			),
		},
		{
			name:      "C# Chromatic",
			root:      note.CSharp,
			intervals: interval.Scales.Chromatic,
			expected: scale.New("C# Chromatic",
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
		{
			name:      "F Chromatic",
			root:      note.F,
			intervals: interval.Scales.Chromatic,
			expected: scale.New("F Chromatic",
				note.F,
				note.GFlat, //note.FSharp,
				note.G,
				note.AFlat, //note.GSharp,
				note.A,
				note.BFlat, //note.ASharp,
				note.CFlat,
				note.C,
				note.DFlat, // note.DFlat,
				note.D,
				note.EFlat, //note.EFlat,
				note.E,
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := scale.Generate(tc.name, tc.root, tc.intervals...)
			require.Equal(t, tc.expected, actual)
		})
	}
}
