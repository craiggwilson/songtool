package interval_test

import (
	"testing"

	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/stretchr/testify/require"
)

func TestNote_Chromatic(t *testing.T) {
	testCases := []struct {
		input    interval.Interval
		expected int
	}{
		{
			input:    interval.Perfect(0),
			expected: 0,
		},
		{
			input:    interval.NewWithChromatic(1, 0),
			expected: 0,
		},
		{
			input:    interval.NewWithChromatic(-1, 0),
			expected: 0,
		},
		{
			input:    interval.NewWithChromatic(0, 2),
			expected: 2,
		},
		{
			input:    interval.NewWithChromatic(3, 5),
			expected: 5,
		},
		{
			input:    interval.NewWithChromatic(4, 7),
			expected: 7,
		},
		{
			input:    interval.NewWithChromatic(4, 8),
			expected: 8,
		},
		{
			input:    interval.NewWithChromatic(5, 0),
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input.String(), func(t *testing.T) {
			actual := tc.input.Chromatic()
			require.Equal(t, tc.expected, actual)
		})
	}
}

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
			by:       interval.Perfect(10),
			expected: interval.Perfect(3),
		},
		{
			input:    interval.Perfect(0),
			by:       interval.NewWithChromatic(-1, -1),
			expected: interval.Major(6),
		},
		{
			input:    interval.NewWithChromatic(3, 5),
			by:       interval.Diminished(4, 1),
			expected: interval.Diminished(0, 1),
		},
		{
			input:    interval.NewWithChromatic(4, 7),
			by:       interval.NewWithChromatic(1, 0),
			expected: interval.NewWithChromatic(5, 7),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input.String()+" by "+tc.by.String(), func(t *testing.T) {
			actual := tc.input.Transpose(tc.by)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func Test_RoundTrip(t *testing.T) {
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
			text:     "4d",
			expected: interval.NewWithChromatic(3, 4),
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
			text:     "6d",
			expected: interval.NewWithChromatic(5, 7),
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
			text:     "6a",
			expected: interval.NewWithChromatic(5, 10),
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
			expected: interval.NewWithChromatic(5, 5),
		},
		{
			text:     "6aaa",
			expected: interval.NewWithChromatic(5, 0),
		},
		{
			text:     "9m",
			expected: interval.NewWithChromatic(8, 13),
		},
		{
			text:     "9M",
			expected: interval.NewWithChromatic(8, 14),
		},
		{
			text:     "9a",
			expected: interval.NewWithChromatic(8, 15),
		},
		{
			text:     "11P",
			expected: interval.NewWithChromatic(10, 17),
		},
		{
			text:     "11a",
			expected: interval.NewWithChromatic(10, 18),
		},
		{
			text:     "13m",
			expected: interval.NewWithChromatic(12, 20),
		},
		{
			text:     "13M",
			expected: interval.NewWithChromatic(12, 21),
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
			text:           "3P",
			expectedErrStr: "only 1, 4, 5, and 11 can be perfect, but got 3",
		},
		{
			text:           "5M",
			expectedErrStr: "only 2, 3, 6, 7, 9, 10, and 13 can be major, but got 5",
		},
		{
			text:           "5m",
			expectedErrStr: "only 2, 3, 6, 7, 9, 10, and 13 can be minor",
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
