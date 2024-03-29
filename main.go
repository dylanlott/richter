package main

import (
	"fmt"
	"log"
)

// Richter is a tiny rules engine in Go, meant to be an extensible starting point.
// It is designed with asynchronous handling as a first-priority, so careful attention
// to ownership of memory is paid.

func main() {
	RunProcess()
}

func RunProcess() {
	var in = make(chan []Card)
	var out = make(chan []Card)
	var errs = make(chan error)
	var state = State{}

	go Process(state, in, out, errs, Config{})

	// add test data
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
	ID      string
	Owner   string
	Name    string
	Tapped  bool
	Cost    int32
	Attack  int32
	Defense int32
}

type Config struct{}

// Rules are the basic unit that is checking if an action can be applied
// to a given State
type Rule struct {
	Name           string
	Condition      func(state State, card Card) bool
	Transformation func(state State, card Card) State
}

// State
type State struct {
	// Board maps player IDs to a map of cards in zones for each player
	Board map[string]map[string][]Card
}

func RunEvaluation(
	state State,
	in chan []Card,
	out chan State,
	rules []Rule,
	errs chan error,
) {
	for stack := range in {
		// process each stack as it comes in
		for _, card := range stack {
			// check if rule applies to the card
			for _, rule := range rules {
				// check the card's condition in the state
				if rule.Condition(state, card) {
					// apply transformation
					log.Printf("applying state transformation: %+v", state)
					state = rule.Transformation(state, card)
				}
			}
		}
	}
}

// Apply checks if each action in the set obeys the rules of the game.
// If it does, it applies that action to the state.
// If it doesn't, it logs the invalid action and doesn't apply it.
// This function must have exclusive access to State.
func Apply(
	state State,
	actions []Action,
) (State, error) {
	var internalState State = state
	for _, a := range actions {
		valid := a.Rule.Condition(internalState, a.Card)
		if valid {
			updated := a.Rule.Transformation(internalState, a.Card)
			internalState = updated
		} else {
			return state, fmt.Errorf("failed to validate rule: %s", a.Rule.Name)
		}
	}
	return internalState, nil
}

// Actions are composed of Rules, which hold a condition and, possibly, a transformation.
// If they contain a transformation at evaluation time, they will be applied to the state.
// If they don't contain a transformation at evaluation time, they're effectively noops.
type Action struct {
	Rule     Rule
	Card     Card
	TargetID string
	Player   string
	Zone     string
}

type Analysis struct {
	ValidActions []Action
}

func RunAnalysis(
	in <-chan State,
	out chan<- Analysis,
	errors chan error,
	config Config,
) {
	for s := range in {
		for player, board := range s.Board {
			fmt.Printf("analyzing player: %s", player)
			for _, v := range board {
				fmt.Printf("analyzing card: %v\n", v)
			}
		}
	}
}

// Analyze takes a state object and a set of rules and analyzes possible moves for
// each player of the game.
func Analyze(
	state State,
	rules []Rule,
) Analysis {
	var analysis = Analysis{}
	for player, board := range state.Board {
		for zone, cards := range board {
			for _, card := range cards {
				for _, rule := range rules {
					if rule.Condition(state, card) {
						analysis.ValidActions = append(analysis.ValidActions, Action{
							Rule:     rule,
							TargetID: card.ID,
							Player:   player,
							Zone:     zone,
							Card:     card,
						})
					}
				}
			}
		}
	}
	return analysis
}

// Process handles continually analyzing and applying actions to the game.
func Process(
	state State,
	in <-chan []Card,
	out chan<- []Card,
	errors chan error,
	config Config,
) {
	for stack := range in {
		fmt.Printf("stack: %v\n", stack)
		fmt.Printf("state: %+v\n", state)

		rules := []Rule{}
		analysis := Analyze(state, rules)
		fmt.Printf("analysis: %v\n", analysis)

		errors <- fmt.Errorf("not impl")
	}
}
