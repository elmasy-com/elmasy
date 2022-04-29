package portscan

import "time"

type State uint8

const (
	OPEN State = iota
	CLOSED
	FILTERED
)

func (s State) String() string {

	switch s {
	case OPEN:
		return "open"
	case CLOSED:
		return "closed"
	case FILTERED:
		return "filtered"
	default:
		return "unknown"
	}
}

// SuggestTimeout suggests a global timeout based on the number of ports
func SuggestTimeout(n int) time.Duration {
	return time.Duration(500+(n*5)) * time.Millisecond
}
