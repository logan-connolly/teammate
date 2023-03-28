package event

// Event is a domain event marker.
type Event interface {
	eventName() string
}
