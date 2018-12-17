package cflp

import (
	"fmt"
	"log"
	"os"
)

// Problem contains data of a CFLP.
type Problem struct {
	N, M        int
	Capacities  []int
	FixedCosts  []int
	Demands     []int
	TotalDemand int
	Costs       [][]int
}

// NewProblem loads a problem from fpath.
func NewProblem(fpath string) *Problem {
	log.Printf("loading %s...\n", fpath)

	// open the file
	f, err := os.Open(fpath)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	// read the first line: the number of facilities n and customers m
	var n, m int
	fmt.Fscanf(f, "%d%d", &n, &m)

	// create the Problem instance
	ret := &Problem{
		N:           n,
		M:           m,
		Capacities:  make([]int, n),
		FixedCosts:  make([]int, n),
		Demands:     make([]int, m),
		TotalDemand: 0,
		Costs:       make([][]int, n),
	}
	for i := 0; i < n; i++ {
		ret.Costs[i] = make([]int, m)
	}

	// read the capacities and fixed costs
	for i := 0; i < n; i++ {
		_, err = fmt.Fscanf(f, "%d%d", &ret.Capacities[i], &ret.FixedCosts[i])
		if err != nil {
			log.Panic(err)
		}
	}
	log.Println("Capacities..OK, FixedCosts..OK")

	// read the demands
	for i := 0; i < m; i++ {
		var val float64
		_, err = fmt.Fscanf(f, "%f", &val)
		if err != nil {
			log.Panic(err)
		}
		ret.Demands[i] = int(val)
		ret.TotalDemand += int(val)
	}
	log.Println("Demands..OK")

	// read the costs
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			var val float64
			_, err = fmt.Fscanf(f, "%f", &val)
			if err != nil {
				log.Panic(err)
			}
			ret.Costs[i][j] = int(val)
		}
	}
	log.Println("Costs..OK")

	fmt.Fprintf(os.Stderr, "loading complete(n=%d,m=%d)\n", n, m)
	fmt.Fprintf(os.Stderr, "capacities:")
	for _, s := range ret.Capacities {
		fmt.Fprintf(os.Stderr, " %d", s)
	}
	fmt.Fprintf(os.Stderr, "\nfixed costs:")
	for _, f := range ret.FixedCosts {
		fmt.Fprintf(os.Stderr, " %d", f)
	}
	fmt.Fprintf(os.Stderr, "\ndemands:")
	for _, d := range ret.Demands {
		fmt.Fprintf(os.Stderr, " %d", d)
	}
	fmt.Fprintf(os.Stderr, "\ntotalDemands: %d", ret.TotalDemand)
	fmt.Fprintf(os.Stderr, "\nc_{ij}(i=0,...,n)(j=0,...,m):\n")
	for _, line := range ret.Costs {
		fmt.Fprintf(os.Stderr, "%d", line[0])
		for c := 1; c < m; c++ {
			fmt.Fprintf(os.Stderr, " %d", line[c])
		}
		fmt.Fprintf(os.Stderr, "\n")
	}

	return ret
}
