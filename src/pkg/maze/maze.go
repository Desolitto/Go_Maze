package maze

import "golang.org/x/exp/rand"

type Cell int

const (
	Empty Cell = iota
	Wall
)

type Maze struct {
	Grid          [][]Cell
	Width, Height int
}

// NewMaze создает новый лабиринт с заданными размерами
func NewMaze(width, height int) *Maze {
	maze := &Maze{
		Grid:   make([][]Cell, height),
		Width:  width,
		Height: height,
	}
	for i := range maze.Grid {
		maze.Grid[i] = make([]Cell, width)
	}
	return maze
}

// GenerateCave генерирует пещеру с использованием клеточного автомата
func (m *Maze) GenerateCave(initialChance, birthLimit, deathLimit int) {
	// Инициализация случайных клеток
	for y := 1; y < m.Height-1; y++ {
		for x := 1; x < m.Width-1; x++ {
			if rand.Intn(100) < initialChance {
				m.Grid[y][x] = Wall
			} else {
				m.Grid[y][x] = Empty
			}
		}
	}

	// Пошаговая генерация
	for i := 0; i < 5; i++ { // Количество итераций
		for y := 1; y < m.Height-1; y++ {
			for x := 1; x < m.Width-1; x++ {
				wallCount := m.CountAliveAround(x, y)
				if wallCount < birthLimit {
					m.Grid[y][x] = Empty
				} else if wallCount > deathLimit {
					m.Grid[y][x] = Wall
				}
			}
		}
	}
}

// CountWallsAround считает количество стен вокруг клетки
func (m *Maze) CountAliveAround(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue // Пропускаем саму клетку
			}
			nx, ny := x+dx, y+dy

			// Проверяем, что индексы находятся в пределах границ
			if nx < 0 || nx >= m.Width || ny < 0 || ny >= m.Height {
				count++ // Считаем клетки за границей как стены
			} else if m.Grid[ny][nx] == Wall {
				count++
			}
		}
	}
	return count
}
