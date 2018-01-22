package main

import (
	"log"
	"perot.me/meiji-chocolate-puzzle-solver/solver"
)

func main() {
	// This creates a 3x2 board with 3 pieces:
	// ┌───┐
	// │A A│
	// └───┘             (A gets rotated)
	//                          ↓
	//   ┌─┐                   ┌─┬─┬─┐
	//   │B│                   │A│C│B│
	// ┌─┘ │    ⇒ Solve() ⇒    │ ├─┘ │
	// │B B│                   │A│B B│
	// └───┘                   └─┴───┘
	//
	// ┌─┐
	// │C│
	// └─┘
	s, err := solver.New(3, 2, `
		**
		
		 *
		**
		
		*
	`)
	if err != nil {
		log.Fatalf("Cannot initialize puzzle: %v", err)
	}
	solution, err := s.Solve()
	if err != nil {
		log.Fatalf("Unsolvable: %v", err)
	}
	solution.Print()
}
