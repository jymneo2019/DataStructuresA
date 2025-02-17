package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReadMatrix reads an augmented matrix from a CSV file and determines `n` automatically
func ReadMatrix(filename string) ([][]float64, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var matrix [][]float64

	// Read all lines into the matrix
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // Skip empty lines
		}

		fields := strings.Split(line, ",") // Split by comma
		row := make([]float64, len(fields))

		for i, val := range fields {
			row[i], err = strconv.ParseFloat(strings.TrimSpace(val), 64)
			if err != nil {
				return nil, 0, fmt.Errorf("invalid number in matrix: %v", err)
			}
		}
		matrix = append(matrix, row)
	}

	// Determine `n` based on the number of rows
	n := len(matrix)
	if n == 0 || len(matrix[0]) != n+1 {
		return nil, 0, fmt.Errorf("invalid matrix format")
	}

	return matrix, n, nil
}

// GaussianElimination solves a system of linear equations using Gaussian elimination
func GaussianElimination(matrix [][]float64, n int) []float64 {
	// Forward elimination
	for i := 0; i < n; i++ {
		// Partial Pivoting: Find the largest element in column i
		maxRow := i
		for j := i + 1; j < n; j++ {
			if abs(matrix[j][i]) > abs(matrix[maxRow][i]) {
				maxRow = j
			}
		}

		// Swap rows
		matrix[i], matrix[maxRow] = matrix[maxRow], matrix[i]

		// Check for zero pivot
		if abs(matrix[i][i]) < 1e-9 {
			fmt.Println("Matrix is singular or nearly singular.")
			return nil
		}

		// Eliminate column elements below pivot
		for j := i + 1; j < n; j++ {
			ratio := matrix[j][i] / matrix[i][i]
			for k := i; k <= n; k++ {
				matrix[j][k] -= ratio * matrix[i][k]
			}
		}
	}

	// Back substitution
	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		x[i] = matrix[i][n]
		for j := i + 1; j < n; j++ {
			x[i] -= matrix[i][j] * x[j]
		}
		x[i] /= matrix[i][i]
	}

	return x
}

// abs returns the absolute value of a float64 number
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	// Read matrix from file
	matrix, n, err := ReadMatrix("matrix.csv")
	if err != nil {
		fmt.Println("Error reading matrix:", err)
		return
	}

	fmt.Println("Original Augmented Matrix:")
	for _, row := range matrix {
		fmt.Println(row)
	}

	// Solve the system using Gaussian Elimination
	solution := GaussianElimination(matrix, n)

	// Print solution
	if solution != nil {
		fmt.Println("\nSolution:")
		for i, val := range solution {
			fmt.Printf("x%d = %.2f\n", i+1, val)
		}
	} else {
		fmt.Println("No unique solution exists.")
	}
}
