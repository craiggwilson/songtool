package theory

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_DegreeClassFromPitchClass(t *testing.T) {
	testCases := []struct {
		pitchClass PitchClass
		enharmonic Enharmonic
		expected   DegreeClass
	}{
		{
			pitchClass: 0,
			enharmonic: EnharmonicSharp,
			expected:   0,
		},
		{
			pitchClass: 1,
			enharmonic: EnharmonicSharp,
			expected:   0,
		},
		{
			pitchClass: 2,
			enharmonic: EnharmonicSharp,
			expected:   1,
		},
		{
			pitchClass: 4,
			enharmonic: EnharmonicSharp,
			expected:   2,
		},
		{
			pitchClass: 5,
			enharmonic: EnharmonicSharp,
			expected:   3,
		},
		{
			pitchClass: 10,
			enharmonic: EnharmonicSharp,
			expected:   5,
		},
		{
			pitchClass: 11,
			enharmonic: EnharmonicSharp,
			expected:   6,
		},
		{
			pitchClass: 0,
			enharmonic: EnharmonicFlat,
			expected:   0,
		},
		{
			pitchClass: 1,
			enharmonic: EnharmonicFlat,
			expected:   1,
		},
		{
			pitchClass: 2,
			enharmonic: EnharmonicFlat,
			expected:   1,
		},
		{
			pitchClass: 4,
			enharmonic: EnharmonicFlat,
			expected:   2,
		},
		{
			pitchClass: 5,
			enharmonic: EnharmonicFlat,
			expected:   3,
		},
		{
			pitchClass: 10,
			enharmonic: EnharmonicFlat,
			expected:   6,
		},
		{
			pitchClass: 11,
			enharmonic: EnharmonicFlat,
			expected:   6,
		},
	}

	cfg := DefaultConfig()

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d %d", tc.pitchClass, tc.enharmonic), func(t *testing.T) {
			actual := cfg.DegreeClassFromPitchClass(tc.pitchClass, tc.enharmonic)
			require.Equal(t, tc.expected, actual)
		})
	}
}
