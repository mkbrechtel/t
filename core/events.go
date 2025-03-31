package core

type Event interface {
	apply(*AppState) error
}
