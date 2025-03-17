package cave

import "golang.org/x/exp/rand"

func (m *Cave) GenerateCave(initialChance, birthLimit, deathLimit int) {
	m.InitializeGrid(initialChance)
	for i := 0; i < 5; i++ {
		m.ApplyCellularAutomaton(birthLimit, deathLimit)
	}
}

func (m *Cave) InitializeGrid(initialChance int) {
	for y := 1; y < m.Height-1; y++ {
		for x := 1; x < m.Width-1; x++ {
			if rand.Intn(100) < initialChance {
				m.Grid[y][x] = Alive
			} else {
				m.Grid[y][x] = Death
			}
		}
	}
}

func (m *Cave) ApplyCellularAutomaton(birthLimit, deathLimit int) {
	newGrid := make([][]Cell, m.Height)
	for j := range newGrid {
		newGrid[j] = make([]Cell, m.Width)
		copy(newGrid[j], m.Grid[j])
	}

	for y := 1; y < m.Height-1; y++ {
		for x := 1; x < m.Width-1; x++ {
			wallCount := m.CountAliveAround(x, y)
			m.UpdateCellState(x, y, wallCount, birthLimit, deathLimit, newGrid)
		}
	}
	m.Grid = newGrid
}

func (m *Cave) UpdateCellState(x, y, wallCount, birthLimit, deathLimit int, newGrid [][]Cell) {
	if m.Grid[y][x] == Alive {
		if wallCount <= deathLimit {
			newGrid[y][x] = Death
		}
	} else {
		if wallCount >= birthLimit {
			newGrid[y][x] = Alive
		}
	}
}
