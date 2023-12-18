package main

import (
	"fmt"
	"log"
)

// Richter is a tiny rules engine in Go, meant to be an extensible starting point.
// It is designed with asynchronous handling as a first-priority, so careful attention
// to ownership of memory is paid.
func main() {
	var in = make(chan []Card)
	var out = make(chan []Card)
	var errs = make(chan error)

	var state = State{}

	go Process(state, in, out, errs, Config{})

	go func() {
		in <- []Card{{Name: "foo"}, {Name: "bar"}}
	}()

	// block while we receive errors, like a server would.
	// TODO: handle context cancellations
	for err := range errs {
		log.Printf("[ERROR]: %+v", err)
		break
	}
}

type Card struct {
	Name string
}

type Config struct{}

type State struct{}

func Process(state State, in <-chan []Card, out chan<- []Card, errors chan error, config Config) {
	for stack := range in {
		fmt.Printf("stack: %v\n", stack)
		fmt.Printf("state: %+v\n", state)
		errors <- fmt.Errorf("not impl")
	}
}
