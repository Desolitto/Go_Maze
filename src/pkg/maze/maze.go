package maze

import (
	"fmt"

	"golang.org/x/exp/rand"
)

type Cell int

const (
	Death Cell = iota
	Alive
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

// GenerateCave генерирует пещеру с использованием клеточного автомата// GenerateCave генерирует пещеру с использованием клеточного автомата
func (m *Maze) GenerateCave(initialChance, birthLimit, deathLimit int) {
	// Инициализация случайных клеток
	for y := 1; y < m.Height-1; y++ {
		for x := 1; x < m.Width-1; x++ {
			if rand.Intn(100) < initialChance {
				m.Grid[y][x] = Alive
			} else {
				m.Grid[y][x] = Death
			}
		}
	}

	// Пошаговая генерация
	for i := 0; i < 5; i++ { // Количество итераций
		newGrid := make([][]Cell, m.Height)
		for j := range newGrid {
			newGrid[j] = make([]Cell, m.Width)
			copy(newGrid[j], m.Grid[j]) // Копируем текущее состояние
		}

		for y := 1; y < m.Height-1; y++ {
			for x := 1; x < m.Width-1; x++ {
				wallCount := m.CountAliveAround(x, y)
				if m.Grid[y][x] == Alive { // Если клетка живая
					if wallCount <= deathLimit {
						newGrid[y][x] = Death // Клетка умирает
					}
				} else { // Если клетка мертвая
					if wallCount >= birthLimit {
						newGrid[y][x] = Alive // Клетка становится живой
					}
				}
			}
		}
		m.Grid = newGrid // Обновляем состояние карты
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
			if (nx < 0 || nx >= m.Width || ny < 0 || ny >= m.Height) || m.Grid[ny][nx] == Alive {
				count++
			}
		}
	}
	return count
}

// func (m *Maze) CountAliveAround(x, y int) int {
// 	count := 0
// 	for dy := -1; dy <= 1; dy++ {
// 		for dx := -1; dx <= 1; dx++ {
// 			if dx == 0 && dy == 0 {
// 				continue // Пропускаем саму клетку
// 			}
// 			nx, ny := x+dx, y+dy

// 			// Проверяем, что индексы находятся в пределах границ
// 			if nx >= 0 && nx < m.Width && ny >= 0 && ny < m.Height {
// 				if m.Grid[ny][nx] == Alive {
// 					count++
// 				}
// 			}
// 		}
// 	}
// 	return count
// }

func (m *Maze) Print() {
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if m.Grid[y][x] == Alive {
				fmt.Print("1 ") // Живая клетка
			} else {
				fmt.Print("0 ") // Мертвая клетка
			}
		}
		fmt.Println() // Переход на новую строку после каждой строки матрицы
	}
}
