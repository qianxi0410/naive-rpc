package config

import "context"

// watch file changes
type Watcher interface {
	Watch(ctx context.Context, fp string) (<-chan Event, error)
}

type EventType int

const (
	Update EventType = iota
	Delete
)

type Event struct {
	typ  EventType
	meta interface{}
}
