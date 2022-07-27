package theory2_test

import (
	"testing"

	"github.com/craiggwilson/songtool/pkg/theory2"
	"github.com/craiggwilson/songtool/pkg/theory2/interval"
	"github.com/craiggwilson/songtool/pkg/theory2/note"
	"github.com/craiggwilson/songtool/pkg/theory2/scale"
	"github.com/stretchr/testify/require"
)

func TestParseScale(t *testing.T) {
	testCases := []struct {
		name           string
		expected       scale.Scale
		expectedErrMsg string
	}{
		{
			name:     "Cb",
			expected: scale.Generate("Cb Major", note.CFlat, interval.Scales.Ionian...),
		},
		{
			name:     "Cb Major",
			expected: scale.Generate("Cb Major", note.CFlat, interval.Scales.Ionian...),
		},
		{
			name:     "Cb Chromatic",
			expected: scale.Generate("Cb Chromatic", note.CFlat, interval.Scales.Chromatic...),
		},
		{
			name:     "C",
			expected: scale.Generate("C Major", note.C, interval.Scales.Ionian...),
		},
		{
			name:     "C Major",
			expected: scale.Generate("C Major", note.C, interval.Scales.Ionian...),
		},
		{
			name:     "C Chromatic",
			expected: scale.Generate("C Chromatic", note.C, interval.Scales.Chromatic...),
		},
		{
			name:     "C#",
			expected: scale.Generate("C# Major", note.CSharp, interval.Scales.Ionian...),
		},
		{
			name:     "C# Major",
			expected: scale.Generate("C# Major", note.CSharp, interval.Scales.Ionian...),
		},
		{
			name:     "C# Chromatic",
			expected: scale.Generate("C# Chromatic", note.CSharp, interval.Scales.Chromatic...),
		},
		{
			name:           "H Major",
			expected:       scale.Scale{},
			expectedErrMsg: `expected natural note name at position 0: expected one of ["C" "D" "E" "F" "G" "A" "B"], but got "H"`,
		},
		{
			name:           "C# Unknown",
			expected:       scale.Scale{},
			expectedErrMsg: `unknown scale name "Unknown"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := theory2.ParseScale(tc.name)
			if len(tc.expectedErrMsg) > 0 {
				require.EqualError(t, err, tc.expectedErrMsg)
			} else {
				require.Nil(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}
