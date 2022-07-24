package interval

type Interval struct {
	Diatonic  int
	Chromatic int
}

func (i Interval) Quality() Quality {
	return -1
}

func (i Interval) String() string {
	return ""
}

func Parse(s string) {

}
