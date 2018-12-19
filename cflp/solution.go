package cflp

import (
	"fmt"
	"io"
	"math"
	"math/rand"
)

// AreaOperator is an operator to change the solution to a near state.
type AreaOperator int

const (
	// OpFlip AreaOperator.
	// Randomly open/close a facility.
	OpFlip AreaOperator = iota
	// OpRangeFlip AreaOperator.
	// Randomly open/close a sequence of facilities.
	OpRangeFlip
	// OpReverse AreaOperator.
	// Randomly reverse a sequence of facilities.
	OpReverse
)

// Solution represents a solution of a CFLP.
type Solution struct {
	*Problem
	X       []bool
	Y       [][]int
	RunTime float64
}

// CopySolution performs a deep copy.
// The problem will not be deeply copied.
func CopySolution(sol *Solution) *Solution {
	ret := &Solution{
		Problem: sol.Problem,
		X:       make([]bool, sol.N),
		Y:       make([][]int, sol.N),
	}
	copy(ret.X, sol.X)
	for i := range ret.Y {
		ret.Y[i] = make([]int, sol.M)
		copy(ret.Y[i], sol.Y[i])
	}
	return ret
}

// Valid determines if the current total capacity can satisfy all demands.
func (s *Solution) Valid() bool {
	totalCapacity := 0
	for i, open := range s.X {
		if open {
			totalCapacity += s.Capacities[i]
		}
	}
	return s.TotalDemand <= totalCapacity
}

// Open opens a facility.
func (s *Solution) Open(i int) {
	s.X[i] = true
}

// RandomAreaOperate randomly picks an operator and operates the solution.
func (s *Solution) RandomAreaOperate() {
	op := AreaOperator(rand.Intn(3))
	s.AreaOperate(op)
}

// Shuffle randomly generates a state of the facilities.
func (s *Solution) Shuffle() {
	state := rand.Intn(1 << uint(s.N))
	for i := range s.X {
		s.X[i] = 1<<uint(i)&state > 0
	}
}

// AreaOperate operates the solution by the given operator.
func (s *Solution) AreaOperate(op AreaOperator) {
	if op == OpFlip {
		pos := rand.Intn(s.N)
		s.X[pos] = !s.X[pos]
	} else if op == OpRangeFlip {
		posX := rand.Intn(s.N)
		posY := rand.Intn(s.N-posX) + posX
		for i := posX; i <= posY; i++ {
			s.X[i] = !s.X[i]
		}
	} else if op == OpReverse {
		posX := rand.Intn(s.N)
		posY := rand.Intn(s.N-posX) + posX
		for i := posX; i < posY/2; i++ {
			s.X[i] = s.X[posY-i]
		}
	}
}

// Assign assigns demands to facility based on current opened facilities.
func (s *Solution) Assign() {
	totalDemand := s.TotalDemand
	filled := make([]int, s.N)
	demands := make([]int, s.M)
	for j := range s.Demands {
		demands[j] = s.Demands[j]
	}
	for i := range s.Y {
		for j := range s.Y[i] {
			s.Y[i][j] = 0
		}
	}
	for totalDemand > 0 {
		// find the cheapest cost
		minI := 0
		minJ := 0
		minCost := math.MaxInt32
		for i := range s.Costs {
			if !s.X[i] || filled[i] == s.Capacities[i] {
				continue
			}
			for j := range s.Costs[i] {
				if s.Costs[i][j] < minCost && demands[j] > 0 {
					minI = i
					minJ = j
					minCost = s.Costs[i][j]
				}
			}
		}

		// assign
		var fill int
		if s.Capacities[minI]-filled[minI] <= demands[minJ] {
			// the facility will be full
			fill = s.Capacities[minI] - filled[minI]
		} else {
			// the facility has more capacity then the demand
			fill = demands[minJ]
		}
		demands[minJ] -= fill
		totalDemand -= fill
		s.Y[minI][minJ] += fill
		filled[minI] += fill
	}
}

// Cost computes the total cost.
func (s *Solution) Cost() float64 {
	cost := 0.0
	for i, x := range s.X {
		if x {
			cost += float64(s.FixedCosts[i])
		}
	}
	for i := range s.Y {
		for j := range s.Y[i] {
			cost += (float64(s.Y[i][j]) / float64(s.Demands[j])) *
				float64(s.Costs[i][j])
		}
	}
	return cost
}

// Display prints the solution.
func (s *Solution) Display(w io.Writer) {
	fmt.Fprintln(w, s.Cost())
	if s.X[0] {
		fmt.Fprint(w, 1)
	} else {
		fmt.Fprint(w, 0)
	}
	for i := 1; i < s.N; i++ {
		if s.X[i] {
			fmt.Fprintf(w, " 1")
		} else {
			fmt.Fprintf(w, " 0")
		}
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "%d", s.Y[0][0])
	for i := range s.Y {
		for j := range s.Y[i] {
			if i == 0 && j == 0 {
				continue
			}
			fmt.Fprintf(w, " %d", s.Y[i][j])
		}
	}
	fmt.Fprintf(w, "\n")
}
