package key

type Namer interface {
	NameKey(Key) string
}
