package cflp

import (
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

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
	// SA Algorithm.
	// Simulated Annealing(SA) Searching to solve the problem.
	SA Algorithm = "sa"
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
	log.Println("running in", alg)
	start := float64(time.Now().UnixNano()) / 1e9
	var sol *Solution
	switch alg {
	case Greedy:
		sol = s.solveByGreedy()
	case BruteForce:
		sol = s.solveByBruteForce()
	case SA:
		sol = s.solveBySA()
	default:
		sol = nil
	}
	sol.RunTime = float64(time.Now().UnixNano())/1e9 - start
	return sol
}

func (s *Solver) initSolution() *Solution {
	sol := &Solution{
		Problem: s.p,
		X:       make([]bool, s.p.N),
		Y:       make([][]int, s.p.N),
	}
	for i := 0; i < s.p.N; i++ {
		sol.Y[i] = make([]int, s.p.M)
	}
	return sol
}

func (s *Solver) solveBySA() *Solution {
	// parameters
	bestSol := s.initSolution()
	bestSol.Shuffle()
	for !bestSol.Valid() {
		bestSol.Shuffle()
	}
	bestSol.Assign()
	bestCost := bestSol.Cost()
	nExtIter := 1000
	nInnIter := 1000
	T := 100.0

	// searching
	for ext := 0; ext < nExtIter; ext++ {
		for inn := 0; inn < nInnIter; inn++ {
			sol := CopySolution(bestSol)
			sol.RandomAreaOperate()
			for !sol.Valid() {
				sol.RandomAreaOperate()
			}
			sol.Assign()
			costNext := sol.Cost()
			if math.Min(1, math.Exp(bestCost-costNext)/T) >= rand.Float64() {
				bestSol = sol
				bestCost = costNext
			}
		}
		log.Printf("ExtIter=%d(%d) T=%.3f bestCost=%.3f\n", ext, nExtIter,
			T, bestCost)
		T *= 0.99
	}

	return bestSol
}

func (s *Solver) solveByBruteForce() *Solution {
	// initialize a solution
	var bestSol *Solution
	bestCost := math.MaxFloat64

	// enumerate all situations
	total := 1 << uint(s.p.N)
	for c := 0; c < total; c++ {
		// output intermidiate results
		if c%2000 == 0 {
			log.Printf("evaluating %s(%d/%d) best=%.3f\n",
				strconv.FormatInt(int64(c), 2), c, total, bestCost)
		}
		// initialize a solution
		sol := s.initSolution()
		// open the facilities
		for i := range sol.X {
			if 1<<uint(i)&c > 0 {
				sol.Open(i)
			}
		}
		// if the solution cannot satisfy all customers, skip
		if !sol.Valid() {
			continue
		}
		// assign and compute the cost
		sol.Assign()
		cost := sol.Cost()
		if cost < bestCost {
			bestCost = cost
			bestSol = sol
		}
	}

	return bestSol
}

func (s *Solver) solveByGreedy() *Solution {
	// initialize a solution
	sol := s.initSolution()

	// find the cheapest facility and turn it on
	minFPos := 0
	minF := math.MaxInt32
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
		minF = math.MaxInt32
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
