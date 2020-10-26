package main

import (
	"fmt"
	"time"
)

type Timer struct {
	name  string
	start time.Time
}

func NewTimer(name string) Timer {
	fmt.Printf("> %v\n", name)
	return Timer{name, time.Now()}
}

func (t Timer) Done() {
	fmt.Printf("< %v %s\n", t.name, time.Since(t.start))
}
