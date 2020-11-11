package main

import (
	"fmt"
	"time"
)

// Timer stores initial timer state and can print time since.
type Timer struct {
	name  string
	start time.Time
}

// NewTimer creates a new timer and prints that it started.
func NewTimer(name string) Timer {
	fmt.Printf("> %v\n", name)
	return Timer{name, time.Now()}
}

// Done prints the time since the timer value was created.
func (t Timer) Done() {
	fmt.Printf("< %v %s\n", t.name, time.Since(t.start))
}
