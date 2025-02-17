package main

import (
	"fmt"
)

// Function to perform Gauss Elimination
func gaussElimination(matrix [][]float64, n int) []float64 {
	// Forward elimination
	for i := 0; i < n; i++ {
		// Make the diagonal element 1
		if matrix[i][i] == 0 {
			fmt.Println("Mathematical Error: Zero pivot element.")
			return nil
		}

		for j := i + 1; j < n; j++ {
			ratio := matrix[j][i] / matrix[i][i]
			for k := 0; k <= n; k++ {
				matrix[j][k] -= ratio * matrix[i][k]
			}
		}
	}

	// Back substitution
	solution := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		solution[i] = matrix[i][n]
		for j := i + 1; j < n; j++ {
			solution[i] -= matrix[i][j] * solution[j]
		}
		solution[i] /= matrix[i][i]
	}

	return solution
}

func main() {
	matrix := [][]float64{
		{1, 1, -1, -2},
		{2, -1, 1, 5},
		{-1, 2, 2, 1},
	}
	n := len(matrix)
	solution := gaussElimination(matrix, n)
	if solution != nil {
		fmt.Println("Solution:")
		for i, val := range solution {
			fmt.Printf("x%d = %.2f\n", i+1, val)
		}
	}
}
