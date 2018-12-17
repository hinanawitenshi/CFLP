// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/hinanawitenshi/CFLP/cflp"
	"github.com/spf13/cobra"
)

var solveInput string
var solveAlg string

// solveCmd represents the solve command
var solveCmd = &cobra.Command{
	Use:   "solve",
	Short: "Solve a capacitated facility location problem.",
	Long: `Solve a capacitated facility location problem with an input problem
and an algorithm specified.`,
	Run: func(cmd *cobra.Command, args []string) {
		solver := cflp.NewSolver(solveInput)
		sol := solver.Solve(cflp.Algorithm(solveAlg))
		sol.Display()
	},
}

func init() {
	rootCmd.AddCommand(solveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// solveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// solveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	solveCmd.Flags().StringVarP(&solveInput, "input", "i", "", `the path of a 
problem, like "res/instances/p1"`)
	solveCmd.Flags().StringVarP(&solveAlg, "alg", "a", "", `the algorithm, from
{greedy, brute-force}`)
}
