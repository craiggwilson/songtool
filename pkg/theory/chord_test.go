package theory_test

import (
	"testing"

	theory2 "github.com/craiggwilson/songtool/pkg/theory"
	"github.com/craiggwilson/songtool/pkg/theory/chord"
	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/craiggwilson/songtool/pkg/theory/note"
	"github.com/stretchr/testify/require"
)

func TestParseChord(t *testing.T) {
	testCases := []struct {
		name           string
		expected       chord.Parsed
		expectedErrMsg string
	}{
		{
			name: "A",
			expected: chord.Parsed{
				Chord: chord.New(note.A, nil, interval.Perfect(0), interval.Major(2), interval.Perfect(4)),
			},
		},
		{
			name: "Am",
			expected: chord.Parsed{
				Chord:  chord.New(note.A, nil, interval.Perfect(0), interval.Minor(2), interval.Perfect(4)),
				Suffix: "m",
			},
		},
		{
			name: "Aaug",
			expected: chord.Parsed{
				Chord:  chord.New(note.A, nil, interval.Perfect(0), interval.Major(2), interval.Augmented(4, 1)),
				Suffix: "aug",
			},
		},
		{
			name: "Adim",
			expected: chord.Parsed{
				Chord:  chord.New(note.A, nil, interval.Perfect(0), interval.Minor(2), interval.Diminished(4, 1)),
				Suffix: "dim",
			},
		},
		{
			name: "Asus",
			expected: chord.Parsed{
				Chord:  chord.New(note.A, nil, interval.Perfect(0), interval.Perfect(3), interval.Perfect(4)),
				Suffix: "sus",
			},
		},
		{
			name: "Asus2",
			expected: chord.Parsed{
				Chord:  chord.New(note.A, nil, interval.Perfect(0), interval.Major(1), interval.Perfect(4)),
				Suffix: "sus2",
			},
		},
		{
			name: "Asus4",
			expected: chord.Parsed{
				Chord:  chord.New(note.A, nil, interval.Perfect(0), interval.Perfect(3), interval.Perfect(4)),
				Suffix: "sus4",
			},
		},
		{
			name: "A/B",
			expected: chord.Parsed{
				Chord:             chord.New(note.A, &note.B, interval.Perfect(0), interval.Major(2), interval.Perfect(4)),
				BaseNoteDelimiter: "/",
			},
		},
		{
			name: "Am7/B",
			expected: chord.Parsed{
				Chord:             chord.New(note.A, &note.B, interval.Perfect(0), interval.Minor(2), interval.Perfect(4), interval.Minor(6)),
				Suffix:            "m7",
				BaseNoteDelimiter: "/",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := theory2.ParseChord(tc.name)
			if len(tc.expectedErrMsg) > 0 {
				require.EqualError(t, err, tc.expectedErrMsg)
			} else {
				require.Nil(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}
