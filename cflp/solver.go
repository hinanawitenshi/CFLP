package cflp

// Algorithm are algorithms available for solving the problem.
type Algorithm string

const (
	// Greedy Algorithm.
	// Greedy Algorithm opens facilities from the cheapest one and until
	// all demands can be satisfied.
	Greedy Algorithm = "greedy"
	// BruteForce Algorithm.
	// BruteForce Algorithm enumerates all possible situations(2^N), which
	// may be slow when the number of facilities is large.
	BruteForce Algorithm = "brute-force"
)

// Solver solves a CFLP.
type Solver struct {
	p *Problem
}

// NewSolver creates a new Solver to solve the problem in pPath.
func NewSolver(pPath string) *Solver {
	return &Solver{
		p: NewProblem(pPath),
	}
}

// Solve solves the problem by the given algorithm.
func (s *Solver) Solve(alg Algorithm) *Solution {
	switch alg {
	case Greedy:
		return s.solveByGreedy()
	default:
		return nil
	}
}

func (s *Solver) solveByGreedy() *Solution {
	// initialize a solution
	sol := &Solution{
		Problem: s.p,
		X:       make([]bool, s.p.N),
		Y:       make([][]int, s.p.N),
	}
	for i := 0; i < s.p.N; i++ {
		sol.Y[i] = make([]int, s.p.M)
	}

	// find the cheapest facility and turn it on
	minFPos := 0
	minF := int(^uint(0) >> 1)
	for i, f := range s.p.FixedCosts {
		if f < minF {
			minFPos = i
			minF = f
		}
	}
	sol.Open(minFPos)

	// open from cheaper until all demands can be satisfied
	for !sol.Valid() {
		minFPos = 0
		minF = int(^uint(0) >> 1)
		for i, f := range s.p.FixedCosts {
			if f < minF && !sol.X[i] {
				minFPos = i
				minF = f
			}
		}
		sol.Open(minFPos)
	}

	// assign the demands
	sol.Assign()

	return sol
}
