package theory_test

import (
	"fmt"
	"testing"

	"github.com/craiggwilson/songtool/pkg/theory"
	theory2 "github.com/craiggwilson/songtool/pkg/theory"
	"github.com/craiggwilson/songtool/pkg/theory/chord"
	"github.com/craiggwilson/songtool/pkg/theory/interval"
	"github.com/craiggwilson/songtool/pkg/theory/note"
	"github.com/stretchr/testify/require"
)

func TestNameChord(t *testing.T) {
	testCases := []struct {
		chord    chord.Chord
		expected string
	}{
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Major(2),
				interval.Perfect(4),
			),
			expected: "C",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Minor(2),
				interval.Perfect(4),
			),
			expected: "Cm",
		},
		// {
		// 	chord: chord.New(note.C, nil,
		// 		interval.Perfect(0),
		// 		interval.Perfect(4),
		// 	),
		// 	expected: "C5",
		// },
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Major(2),
				interval.Augmented(4, 1),
			),
			expected: "Caug",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Minor(2),
				interval.Diminished(4, 1),
			),
			expected: "Cdim",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Minor(2),
				interval.Diminished(4, 1),
				interval.Minor(6),
			),
			expected: "Cm7b5",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Minor(2),
				interval.Diminished(4, 1),
				interval.Diminished(6, 1),
			),
			expected: "Cdim7",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Perfect(3),
				interval.Perfect(4),
			),
			expected: "Csus",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Major(1),
				interval.Perfect(4),
			),
			expected: "C2",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Major(2),
				interval.Perfect(4),
				interval.Major(5),
			),
			expected: "C6",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Minor(2),
				interval.Perfect(4),
				interval.Major(5),
			),
			expected: "Cm6",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Major(2),
				interval.Perfect(4),
				interval.Major(5),
				interval.Major(8),
			),
			expected: "C6add9",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Minor(2),
				interval.Perfect(4),
				interval.Major(5),
				interval.Major(8),
			),
			expected: "Cm6add9",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Major(1),
				interval.Major(2),
				interval.Perfect(4),
			),
			expected: "Cadd2",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Major(1),
				interval.Minor(2),
				interval.Perfect(4),
			),
			expected: "Cmadd2",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Major(2),
				interval.Perfect(3),
				interval.Perfect(4),
			),
			expected: "Cadd4",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Minor(2),
				interval.Perfect(3),
				interval.Perfect(4),
			),
			expected: "Cmadd4",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Major(2),
				interval.Perfect(4),
				interval.Major(8),
			),
			expected: "Cadd9",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Minor(2),
				interval.Perfect(4),
				interval.Major(8),
			),
			expected: "Cmadd9",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Major(2),
				interval.Perfect(4),
				interval.Perfect(10),
			),
			expected: "Cadd11",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Major(2),
				interval.Perfect(4),
				interval.Major(12),
			),
			expected: "Cadd13",
		},
		{
			chord: chord.New(note.C, nil,
				interval.Perfect(0),
				interval.Perfect(3),
				interval.Perfect(4),
				interval.Minor(6),
			),
			expected: "C7sus",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.chord.Intervals()), func(t *testing.T) {
			actual := theory.NameChord(tc.chord)
			require.Equal(t, tc.expected, actual)
		})
	}
}

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
