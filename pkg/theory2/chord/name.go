package chord

type Namer interface {
	NameChord(Chord) string
}
