package interval_test

import (
	"fmt"
	"testing"

	"github.com/craiggwilson/songtool/pkg/theory2/interval"
	"github.com/stretchr/testify/require"
)

func TestInterval_Transpose(t *testing.T) {
	testCases := []struct {
		input    interval.Interval
		by       interval.Interval
		expected interval.Interval
	}{
		{
			input:    interval.Perfect(0),
			by:       interval.Perfect(0),
			expected: interval.Perfect(0),
		},
		{
			input:    interval.Perfect(0),
			by:       interval.Perfect(4),
			expected: interval.Perfect(4),
		},
		{
			input:    interval.Perfect(0),
			by:       interval.NewWithChromatic(-1, -1),
			expected: interval.Major(6),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input.String()+" by "+tc.by.String(), func(t *testing.T) {
			actual := tc.input.Transpose(tc.by)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestFromStep(t *testing.T) {
	testCases := []struct {
		step     int
		expected interval.Interval
	}{
		{
			step:     -7,
			expected: interval.NewWithChromatic(3, 5),
		},
		{
			step:     -1,
			expected: interval.NewWithChromatic(6, 11),
		},
		{
			step:     0,
			expected: interval.NewWithChromatic(0, 0),
		},
		{
			step:     1,
			expected: interval.NewWithChromatic(1, 1),
		},
		{
			step:     2,
			expected: interval.NewWithChromatic(1, 2),
		},
		{
			step:     3,
			expected: interval.NewWithChromatic(2, 3),
		},
		{
			step:     4,
			expected: interval.NewWithChromatic(2, 4),
		},
		{
			step:     5,
			expected: interval.NewWithChromatic(3, 5),
		},
		{
			step:     6,
			expected: interval.NewWithChromatic(4, 6),
		},
		{
			step:     7,
			expected: interval.NewWithChromatic(4, 7),
		},
		{
			step:     8,
			expected: interval.NewWithChromatic(5, 8),
		},
		{
			step:     9,
			expected: interval.NewWithChromatic(5, 9),
		},
		{
			step:     10,
			expected: interval.NewWithChromatic(6, 10),
		},
		{
			step:     11,
			expected: interval.NewWithChromatic(6, 11),
		},
		{
			step:     12,
			expected: interval.NewWithChromatic(0, 0),
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("step %d", tc.step), func(t *testing.T) {
			actual := interval.FromStep(tc.step)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestParseRoundTrip(t *testing.T) {
	testCases := []struct {
		text           string
		expected       interval.Interval
		expectedErrStr string
	}{
		{
			text:     "1P",
			expected: interval.NewWithChromatic(0, 0),
		},
		{
			text:     "2m",
			expected: interval.NewWithChromatic(1, 1),
		},
		{
			text:     "2M",
			expected: interval.NewWithChromatic(1, 2),
		},
		{
			text:     "3m",
			expected: interval.NewWithChromatic(2, 3),
		},
		{
			text:     "3M",
			expected: interval.NewWithChromatic(2, 4),
		},
		{
			text:     "4P",
			expected: interval.NewWithChromatic(3, 5),
		},
		{
			text:     "4a",
			expected: interval.NewWithChromatic(3, 6),
		},
		{
			text:     "5d",
			expected: interval.NewWithChromatic(4, 6),
		},
		{
			text:     "5P",
			expected: interval.NewWithChromatic(4, 7),
		},
		{
			text:     "6m",
			expected: interval.NewWithChromatic(5, 8),
		},
		{
			text:     "6M",
			expected: interval.NewWithChromatic(5, 9),
		},
		{
			text:     "7m",
			expected: interval.NewWithChromatic(6, 10),
		},
		{
			text:     "7M",
			expected: interval.NewWithChromatic(6, 11),
		},
		{
			text:     "6ddd",
			expected: interval.NewWithChromatic(5, 6),
		},
		{
			text:     "6aaa",
			expected: interval.NewWithChromatic(5, 0),
		},
		{
			text:           "1",
			expectedErrStr: "intervals must contain at least 2 characters, but had 1",
		},
		{
			text:           "P",
			expectedErrStr: "intervals must contain at least 2 characters, but had 1",
		},
		{
			text:           "9P",
			expectedErrStr: "expected a number between 1 and 7, but got 9",
		},
		{
			text:           "3P",
			expectedErrStr: "only 1, 4, and 5 can be perfect, but got 3",
		},
		{
			text:           "5M",
			expectedErrStr: "only 2, 3, 6, and 7 can be major, but got 5",
		},
		{
			text:           "5m",
			expectedErrStr: "only 2, 3, 6, and 7 can be minor, but got 5",
		},
		{
			text:           "5PP",
			expectedErrStr: "perfect interval quality has no size",
		},
		{
			text:           "2MM",
			expectedErrStr: "major interval quality has no size",
		},
		{
			text:           "2mm",
			expectedErrStr: "minor interval quality has no size",
		},
		{
			text:           "2aad",
			expectedErrStr: "cannot mix interval qualities; expected 'a', but got 'd' at pos 3",
		},
		{
			text:           "2dda",
			expectedErrStr: "cannot mix interval qualities; expected 'd', but got 'a' at pos 3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.text, func(t *testing.T) {
			actual, err := interval.Parse(tc.text)
			if len(tc.expectedErrStr) > 0 {
				require.ErrorContains(t, err, tc.expectedErrStr)
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.expected, actual)
				require.Equal(t, tc.text, actual.String())
			}
		})
	}
}
