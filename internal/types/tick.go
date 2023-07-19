package types

// Tick is a timestamp/value pair.
type Tick[T any] struct {
	Time  uint64
	Value T
}
