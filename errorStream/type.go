package errorStream

type Error struct {
	Error func() string
}
