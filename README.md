# Richter

> a tiny rules engine in Go

Richter is a tiny rules engine written as an experiment to analyze and evaluate arbitrary board states in trading card games.

Its main goal is as a starting point for a full-bodied approach to board state management in a larger simulator or arena.

## Design

The engine is designed around two phases: Analysis and Evaluation.

Analysis is focused on coming up with possible moves, while evaluation is focused around validating and applying a set of actions to the state.

Rules are composed of a condition function that is passed the current state and the card in question. From there, a bool is returned - true if the rule is met, false if the rule is not met and thus is an invalid action. Rules allow arbitrary logic to be written into conditions, enabling it to check the state for any possible value at evaluation time, as well as being able to peer up its own stack and sibling stacks.

Stacks are sets of cards that are being analyzed or evaluated. A set of stacks is a player's boardstate, and a collection of those players is a game.

## Implementation

The `RunProcess` function starts the rules engine and offers read only access to the Analyzer, while safely offering atomic and exclusive write access to the Evaluation phase.
