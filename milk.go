package main

import (
	"log"
	"perot.me/meiji-chocolate-puzzle-solver/solver"
)

func main() {
	s, err := solver.New(10, 6, `
		****
		*
		
		***
		 **
		
		****
		 *
		
		 *
		***
		 *
		
		**
		 *
		**
		
		*
		*
		***
		
		  *
		***
		*
		
		**
		 **
		  *
		
		 *
		***
		  *
		
		***
		 *
		 *
		
		  **
		***
		
		*****
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
