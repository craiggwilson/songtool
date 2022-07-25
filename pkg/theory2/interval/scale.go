package interval

var Scales = struct {
	Chromatic []Interval
	Ionian    []Interval
}{
	Chromatic: []Interval{
		Perfect(0),
		Minor(1),
		Major(1),
		Minor(2),
		Major(2),
		Perfect(3),
		Diminished(4, 1),
		Perfect(4),
		Minor(5),
		Major(5),
		Minor(6),
		Major(6),
	},
	Ionian: []Interval{
		Perfect(0),
		Major(1),
		Major(2),
		Perfect(3),
		Perfect(4),
		Major(5),
		Major(6),
	},
}
