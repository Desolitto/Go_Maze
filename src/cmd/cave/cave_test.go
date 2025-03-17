package main

import (
	"go-maze/pkg/cave"
	"testing"
)

// func TestApplyCellularAutomaton(t *testing.T) {
// 	c := &cave.Cave{
// 		Width:  10,
// 		Height: 10,
// 		Grid: [][]cave.Cell{
// 			{cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
// 			{cave.Death, cave.Alive, cave.Alive, cave.Alive, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
// 			{cave.Death, cave.Alive, cave.Alive, cave.Alive, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
// 			{cave.Death, cave.Alive, cave.Alive, cave.Alive, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
// 			{cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
// 			{cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
// 			{cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
// 			{cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
// 			{cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
// 			{cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
// 		},
// 	}

// 	c.ApplyCellularAutomaton(10, 1)

// 	// Проверяем, что после применения клеточного автомата, сетка изменилась
// 	for y := 0; y < c.Height; y++ {
// 		for x := 0; x < c.Width; x++ {
// 			if c.Grid[y][x] == c.Grid[y][x] {
// 				t.Errorf("Ячейка (%d, %d) не изменилась", x, y)
// 			}
// 		}
// 	}
// }

func TestUpdateCellState(t *testing.T) {
	c := &cave.Cave{
		Width:  10,
		Height: 10,
		Grid: [][]cave.Cell{
			{cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
			{cave.Death, cave.Alive, cave.Alive, cave.Alive, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
			{cave.Death, cave.Alive, cave.Alive, cave.Alive, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
			{cave.Death, cave.Alive, cave.Alive, cave.Alive, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
			{cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
			{cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
			{cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
			{cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
			{cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
			{cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death, cave.Death},
		},
	}

	newGrid := make([][]cave.Cell, c.Height)
	for j := range newGrid {
		newGrid[j] = make([]cave.Cell, c.Width)
		copy(newGrid[j], c.Grid[j])
	}

	c.UpdateCellState(1, 1, 2, 4, 3, newGrid)
	c.UpdateCellState(1, 2, 3, 4, 3, newGrid)
	c.UpdateCellState(1, 3, 4, 4, 3, newGrid)

	// Проверяем, что состояние ячеек изменилось в соответствии с правилами
	if newGrid[1][1] != cave.Death {
		t.Errorf("Ячейка (1, 1) должна быть в состоянии Death")
	}
	if newGrid[1][2] != cave.Alive {
		t.Errorf("Ячейка (1, 2) должна быть в состоянии Alive")
	}
	if newGrid[1][3] != cave.Alive {
		t.Errorf("Ячейка (1, 3) должна быть в состоянии Alive")
	}
}
