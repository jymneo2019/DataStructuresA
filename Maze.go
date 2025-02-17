package main

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Cell struct {
	X, Y int
}

type Maze struct {
	Grid          [][]int
	Width, Height int
}

func NewMaze(width, height int) *Maze {
	grid := make([][]int, height)
	for i := range grid {
		grid[i] = make([]int, width)
	}
	return &Maze{Grid: grid, Width: width, Height: height}
}

func (m *Maze) GenerateDFS(x, y int) {
	m.Grid[y][x] = 1 // Mark the starting cell as open
	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	for _, dir := range directions {
		newX, newY := x+dir[0]*2, y+dir[1]*2
		if newX >= 0 && newX < m.Width && newY >= 0 && newY < m.Height && m.Grid[newY][newX] == 0 {
			m.Grid[y+dir[1]][x+dir[0]] = 1
			m.Grid[newY][newX] = 1
			m.GenerateDFS(newX, newY)
		}
	}
}

func (m *Maze) Dijkstra(start, end Cell) []Cell {
	distances := make([][]float64, m.Height)
	previous := make([][]*Cell, m.Height)
	for i := range distances {
		distances[i] = make([]float64, m.Width)
		previous[i] = make([]*Cell, m.Width)
		for j := range distances[i] {
			distances[i][j] = math.Inf(1)
		}
	}
	distances[start.Y][start.X] = 0

	pq := &PriorityQueue{}
	heap.Push(pq, &Item{Value: start, Priority: 0})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*Item).Value

		if current.X == end.X && current.Y == end.Y {
			var path []Cell
			for c := &current; c != nil; c = previous[c.Y][c.X] {
				path = append(path, *c)
			}
			// Reverse path to correct order
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}
			return path
		}

		for _, dir := range [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
			newX, newY := current.X+dir[0], current.Y+dir[1]
			if newX >= 0 && newX < m.Width && newY >= 0 && newY < m.Height && m.Grid[newY][newX] == 1 {
				newDist := distances[current.Y][current.X] + 1
				if newDist < distances[newY][newX] {
					distances[newY][newX] = newDist
					previous[newY][newX] = &Cell{current.X, current.Y} // Store a copy
					heap.Push(pq, &Item{Value: Cell{newX, newY}, Priority: newDist})
				}
			}
		}
	}
	return nil
}

type Item struct {
	Value    Cell
	Priority float64
	Index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Priority < pq[j].Priority }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i]; pq[i].Index, pq[j].Index = i, j }

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.Index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[0]
	*pq = old[1:n]
	return item
}

func main() {
	maze := NewMaze(21, 21)
	maze.GenerateDFS(1, 1)
	start := Cell{1, 1}
	end := Cell{19, 19}
	path := maze.Dijkstra(start, end)

	if path != nil {
		fmt.Println("Path from start to end:")
		for _, cell := range path {
			fmt.Printf("(%d, %d) ", cell.X, cell.Y)
		}
		fmt.Println()
	} else {
		fmt.Println("No path found.")
	}
}
