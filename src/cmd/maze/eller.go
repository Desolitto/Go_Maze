package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Cell struct {
	RightWall  bool
	BottomWall bool
	Set        int // Множество, к которому принадлежит ячейка
}

type Maze struct {
	Rows  int
	Cols  int
	Cells [][]Cell
}

func NewMaze(rows, cols int) *Maze {
	cells := make([][]Cell, rows)
	for i := range cells {
		cells[i] = make([]Cell, cols)
		for j := range cells[i] {
			cells[i][j] = Cell{
				RightWall:  false,
				BottomWall: false,
				Set:        -1, // Уникальное множество, будет назначено позже
			}
		}
	}

	// // Печатаем состояние всех ячеек
	// for i := 0; i < rows; i++ {
	// 	for j := 0; j < cols; j++ {
	// 		cell := cells[i][j]
	// 		fmt.Printf("Cell(%d, %d): RightWall=%t, BottomWall=%t, Set=%d\n", i, j, cell.RightWall, cell.BottomWall, cell.Set)
	// 	}
	// }

	return &Maze{Rows: rows, Cols: cols, Cells: cells}
}

func (m *Maze) Generate(randomNumbers []int) {
	// Инициализация множеств для каждой ячейки
	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			m.Cells[row][col].Set = row*m.Cols + col + 1 // Уникальные значения для каждой ячейки
		}
	}

	// Установка правых стенок
	index := 0
	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			if col < m.Cols-1 {
				fmt.Printf("randomNumbers[index] right = %d\n", randomNumbers[index])

				if randomNumbers[index] == 1 {
					m.Cells[row][col].RightWall = true
				} else {
					// Объединяем множества
					set1 := m.Cells[row][col].Set
					set2 := m.Cells[row][col+1].Set

					if set1 != set2 {
						// Убираем стенку между текущей ячейкой и ячейкой справа
						m.Cells[row][col].RightWall = false

						// Объединяем множества
						for r := 0; r < m.Rows; r++ {
							for c := 0; c < m.Cols; c++ {
								if m.Cells[r][c].Set == set2 {
									m.Cells[r][c].Set = set1 // Обновляем Set для всех ячеек, принадлежащих set2
								}
							}
						}
					}
				}
				index++ // Переход к следующему числу для правых стенок
			}
		}

		// Установка нижних стенок и объединение множеств для последней строки
		if row == m.Rows-1 { // Проверяем, является ли текущая строка последней
			for col := 0; col < m.Cols; col++ {
				m.Cells[row][col].BottomWall = true // Присваиваем нижнюю стенку всем ячейкам в последней строке
			}

			// Убираем стенки между ячейками и объединяем множества
			for col := 0; col < m.Cols-1; col++ {
				set1 := m.Cells[row][col].Set
				set2 := m.Cells[row][col+1].Set

				if set1 != set2 {
					// Убираем стенку между текущей ячейкой и ячейкой справа
					m.Cells[row][col].RightWall = false

					// Объединяем множества
					for r := 0; r < m.Rows; r++ {
						for c := 0; c < m.Cols; c++ {
							if m.Cells[r][c].Set == set2 {
								m.Cells[r][c].Set = set1 // Обновляем Set для всех ячеек, принадлежащих set2
							}
						}
					}
				}
			}
		} else {
			for col := 0; col < m.Cols; col++ {
				if row < m.Rows-1 {
					// Проверяем, нужно ли ставить нижнюю стенку
					set := m.Cells[row][col].Set
					count := 0

					// Считаем количество ячеек в текущем множестве, которые не имеют нижней границы
					for c := 0; c < m.Cols; c++ {
						if m.Cells[row][c].Set == set && !m.Cells[row][c].BottomWall {
							count++
						}
					}
					fmt.Printf("Множества --- %d ячейки ---%d \n", set, count)

					// Если множество содержит более одной ячейки без нижней границы
					if count > 1 {
						// Устанавливаем нижнюю стенку, только если randomNumbers[index] == 1
						fmt.Printf("randomNumbers[index] bottom = %d\n", randomNumbers[index])
						if randomNumbers[index] == 1 {
							m.Cells[row][col].BottomWall = true
						}
					}

					index++ // Переход к следующему числу для нижних стенок
				}
			}
		}
	}

	// Вывод состояния ячеек
	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			fmt.Printf("Cell(%d, %d): RightWall=%v, BottomWall=%v, Set=%d\n",
				row, col, m.Cells[row][col].RightWall, m.Cells[row][col].BottomWall, m.Cells[row][col].Set)
		}
	}
}

func (m *Maze) PrintSets() {
	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			fmt.Printf("Ячейка [%d][%d]: Set = %d\n", row, col, m.Cells[row][col].Set)
		}
	}
}

type Game struct {
	maze     *Maze
	cellSize float64 // Изменено на float64
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	strokeColor := color.RGBA{0, 0, 0, 255}
	fillColor := color.RGBA{255, 255, 255, 255}

	// Рисуем лабиринт
	for y := 0; y < g.maze.Rows; y++ {
		for x := 0; x < g.maze.Cols; x++ {
			// Рисуем ячейку
			ebitenutil.DrawRect(screen, float64(x)*g.cellSize, float64(y)*g.cellSize, g.cellSize, g.cellSize, fillColor)

			// Рисуем правую границу
			if x < g.maze.Cols-1 && g.maze.Cells[y][x].RightWall {
				ebitenutil.DrawLine(screen, float64(x+1)*g.cellSize, float64(y)*g.cellSize, float64(x+1)*g.cellSize, float64(y+1)*g.cellSize, strokeColor)
			}

			// Рисуем нижнюю границу
			if y < g.maze.Rows-1 && g.maze.Cells[y][x].BottomWall {
				ebitenutil.DrawLine(screen, float64(x)*g.cellSize, float64(y+1)*g.cellSize, float64(x+1)*g.cellSize, float64(y+1)*g.cellSize, strokeColor)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {

	randomNumbers := make([]int, 0) // Для 4 строк по 4 столбца
	randomNumbers = append(randomNumbers, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 1, 0)

	fmt.Println(randomNumbers)
	maze := NewMaze(4, 4)
	// r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	// rows, cols := 5, 5
	// numRandomNumbers := rows * cols * 2
	// randomNumbers := make([]int, numRandomNumbers)
	// for i := range randomNumbers {
	// 	randomNumbers[i] = r.Intn(2) // Генерация 0 или 1
	// }
	// fmt.Println(randomNumbers)

	// maze := NewMaze(rows, cols)
	maze.Generate(randomNumbers)
	// Печатаем множества ячеек
	// maze.PrintSets()
	game := &Game{maze: maze, cellSize: 50.0} // Изменено на 40.0
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Maze Generator")
	if err := ebiten.RunGame(game); err != nil {
		fmt.Println(err)
	}
}
