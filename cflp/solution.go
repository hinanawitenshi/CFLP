package cflp

import "fmt"

// Solution represents a solution of a CFLP.
type Solution struct {
	*Problem
	X []bool
	Y [][]int
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

// Assign assigns demands to facility based on current opened facilities.
func (s *Solution) Assign() {
	totalDemand := s.TotalDemand
	filled := make([]int, s.N)
	demands := make([]int, s.M)
	for j := range s.Demands {
		demands[j] = s.Demands[j]
	}
	for totalDemand > 0 {
		// find the cheapest cost
		minI := 0
		minJ := 0
		minCost := int(^uint(0) >> 1)
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
func (s *Solution) Display() {
	fmt.Println(s.Cost())
	if s.X[0] {
		fmt.Print(1)
	} else {
		fmt.Print(0)
	}
	for i := 1; i < s.N; i++ {
		if s.X[i] {
			fmt.Printf(" 1")
		} else {
			fmt.Printf(" 0")
		}
	}
	fmt.Printf("\n")
	for i := range s.Y {
		fmt.Printf("%d", s.Y[i][0])
		for j := range s.Y[i] {
			if j == 0 {
				continue
			}
			fmt.Printf("\t%d", s.Y[i][j])
		}
		fmt.Printf("\n")
	}
}
