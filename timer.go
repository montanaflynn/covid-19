package main

import (
	"fmt"
	"time"
)

type timer struct {
	name      string
	startTime time.Time
}

func (t *timer) end() *timer {
	fmt.Println(t.name, "took", time.Now().Sub(t.startTime))
	return t
}

func (t *timer) reset(name string) *timer {
	t.name = name
	t.startTime = time.Now()
	return t
}

func newTimer(name string) *timer {
	return &timer{
		name:      name,
		startTime: time.Now(),
	}
}
