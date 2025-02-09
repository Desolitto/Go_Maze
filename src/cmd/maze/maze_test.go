package main

import (
	"go-maze/pkg/maze"
	"testing"
)

func TestInitializeSets(t *testing.T) {
	// Create a new Maze instance
	mazeInstance := maze.NewMaze(3, 3)

	// Initialize sets
	mazeInstance.InitializeSets()

	// Validate the sets
	for row := 0; row < mazeInstance.Rows; row++ {
		for col := 0; col < mazeInstance.Cols; col++ {
			expectedSet := row*mazeInstance.Cols + col + 1
			if mazeInstance.Cells[row][col].Set != expectedSet {
				t.Errorf("Expected Set[%d][%d] = %d, got %d", row, col, expectedSet, mazeInstance.Cells[row][col].Set)
			}
		}
	}
}

// func TestCopyPreviousRow(t *testing.T) {
// 	maze := &Maze{Rows: 2, Cols: 2}
// 	maze.Cells = make([][]Cell, maze.Rows)
// 	for i := range maze.Cells {
// 		maze.Cells[i] = make([]Cell, maze.Cols)
// 	}

// 	maze.Cells[0][0].Right = true
// 	maze.Cells[0][0].Bottom = false
// 	maze.Cells[0][0].Set = 1

// 	currentSetCount := 2
// 	maze.copyPreviousRow(1, &currentSetCount)

// 	if maze.Cells[1][0].Right != true {
// 		t.Errorf("Expected Right = true, got %v", maze.Cells[1][0].Right)
// 	}
// 	if maze.Cells[1][0].Bottom != false {
// 		t.Errorf("Expected Bottom = false, got %v", maze.Cells[1][0].Bottom)
// 	}
// 	if maze.Cells[1][0].Set != 1 {
// 		t.Errorf("Expected Set = 1, got %d", maze.Cells[1][0].Set)
// 	}
// }

// func TestSetFirstRowSets(t *testing.T) {
// 	maze := &Maze{Rows: 1, Cols: 3}
// 	maze.Cells = make([][]Cell, maze.Rows)
// 	for i := range maze.Cells {
// 		maze.Cells[i] = make([]Cell, maze.Cols)
// 	}

// 	currentSetCount := 0
// 	maze.setFirstRowSets(&currentSetCount)

// 	for col := 0; col < maze.Cols; col++ {
// 		expectedSet := col + 1
// 		if maze.Cells[0][col].Set != expectedSet {
// 			t.Errorf("Expected Set[0][%d] = %d, got %d", col, expectedSet, maze.Cells[0][col].Set)
// 		}
// 	}
// }

// func TestProcessRightWalls(t *testing.T) {
// 	maze := &Maze{Rows: 1, Cols: 3}
// 	maze.Cells = make([][]Cell, maze.Rows)
// 	for i := range maze.Cells {
// 		maze.Cells[i] = make([]Cell, maze.Cols)
// 	}

// 	randomNumbers := []int{1, 0, 1}
// 	index := 0
// 	maze.processRightWalls(0, randomNumbers, &index)

// 	if maze.Cells[0][0].Right != true {
// 		t.Errorf("Expected Right[0][0] = true, got %v", maze.Cells[0][0].Right)
// 	}
// 	if maze.Cells[0][1].Right != false {
// 		t.Errorf("Expected Right[0][1] = false, got %v", maze.Cells[0][1].Right)
// 	}
// 	if maze.Cells[0][2].Right != true {
// 		t.Errorf("Expected Right[0][2] = true, got %v", maze.Cells[0][2].Right)
// 	}
// }

// func TestMergeSets(t *testing.T) {
// 	maze := &Maze{Rows: 2, Cols: 2}
// 	maze.Cells = make([][]Cell, maze.Rows)
// 	for i := range maze.Cells {
// 		maze.Cells[i] = make([]Cell, maze.Cols)
// 	}

// 	maze.Cells[0][0].Set = 1
// 	maze.Cells[0][1].Set = 2

// 	maze.mergeSets(1, 2)

// 	for row := 0; row < maze.Rows; row++ {
// 		for col := 0; col < maze.Cols; col++ {
// 			if maze.Cells[row][col].Set != 1 {
// 				t.Errorf("Expected Set[%d][%d] = 1, got %d", row, col, maze.Cells[row][col].Set)
// 			}
// 		}
// 	}
// }
