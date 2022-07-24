package interval_test

import (
	"testing"

	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/stretchr/testify/require"
)

func TestParseRoundTrip(t *testing.T) {
	testCases := []struct {
		text           string
		expected       interval.Interval
		expectedErrStr string
	}{
		{
			text:     "1P",
			expected: interval.Interval{0, 0},
		},
		{
			text:     "2m",
			expected: interval.Interval{1, 1},
		},
		{
			text:     "2M",
			expected: interval.Interval{1, 2},
		},
		{
			text:     "3m",
			expected: interval.Interval{2, 3},
		},
		{
			text:     "3M",
			expected: interval.Interval{2, 4},
		},
		{
			text:     "4P",
			expected: interval.Interval{3, 5},
		},
		{
			text:     "4a",
			expected: interval.Interval{3, 6},
		},
		{
			text:     "5d",
			expected: interval.Interval{4, 6},
		},
		{
			text:     "5P",
			expected: interval.Interval{4, 7},
		},
		{
			text:     "6m",
			expected: interval.Interval{5, 8},
		},
		{
			text:     "6M",
			expected: interval.Interval{5, 9},
		},
		{
			text:     "7m",
			expected: interval.Interval{6, 10},
		},
		{
			text:     "7M",
			expected: interval.Interval{6, 11},
		},
		{
			text:     "6ddd",
			expected: interval.Interval{5, 6},
		},
		{
			text:     "6aaa",
			expected: interval.Interval{5, 0},
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
