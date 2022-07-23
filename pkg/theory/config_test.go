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
			enharmonic: Sharp,
			expected:   0,
		},
		{
			pitchClass: 1,
			enharmonic: Sharp,
			expected:   0,
		},
		{
			pitchClass: 2,
			enharmonic: Sharp,
			expected:   1,
		},
		{
			pitchClass: 4,
			enharmonic: Sharp,
			expected:   2,
		},
		{
			pitchClass: 5,
			enharmonic: Sharp,
			expected:   3,
		},
		{
			pitchClass: 10,
			enharmonic: Sharp,
			expected:   5,
		},
		{
			pitchClass: 11,
			enharmonic: Sharp,
			expected:   6,
		},
		{
			pitchClass: 0,
			enharmonic: Flat,
			expected:   0,
		},
		{
			pitchClass: 1,
			enharmonic: Flat,
			expected:   1,
		},
		{
			pitchClass: 2,
			enharmonic: Flat,
			expected:   1,
		},
		{
			pitchClass: 4,
			enharmonic: Flat,
			expected:   2,
		},
		{
			pitchClass: 5,
			enharmonic: Flat,
			expected:   3,
		},
		{
			pitchClass: 10,
			enharmonic: Flat,
			expected:   6,
		},
		{
			pitchClass: 11,
			enharmonic: Flat,
			expected:   6,
		},
	}

	cfg := DefaultConfig()

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d %s", tc.pitchClass, tc.enharmonic), func(t *testing.T) {
			actual := cfg.DegreeClassFromPitchClass(tc.pitchClass, tc.enharmonic)
			require.Equal(t, tc.expected, actual)
		})
	}
}
