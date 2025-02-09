package maze

import (
	"fmt"
	"log"
)

type MazeSolving struct {
	mazeInfo      *Maze
	startX        int
	startY        int
	endX          int
	endY          int
	solvingMatrix [][]int
	currentX      int
	currentY      int
	isStartSet    bool
	isEndSet      bool
	path          []Point // Для хранения найденного пути
}

type Point struct {
	X int
	Y int
}

func NewMazeSolving(maze *Maze, startX, startY, endX, endY int) *MazeSolving {
	solvingMatrix := make([][]int, maze.Rows)
	for i := range solvingMatrix {
		solvingMatrix[i] = make([]int, maze.Cols)
	}
	return &MazeSolving{
		mazeInfo:      maze,
		startX:        startX,
		startY:        startY,
		endX:          endX,
		endY:          endY,
		solvingMatrix: solvingMatrix,
		currentX:      startX,
		currentY:      startY,
		isStartSet:    startX >= 0 && startY >= 0,
		isEndSet:      endX >= 0 && endY >= 0,
		path:          []Point{},
	}
}

// Метод Solve запускает решение лабиринта
func (m *MazeSolving) Solve() error {
	if !m.isStartSet || !m.isEndSet {
		return fmt.Errorf("начальная и конечная точки должны быть установлены")
	}
	m.path = []Point{} // Очищаем предыдущий путь
	log.Println("Запуск решения лабиринта...")
	if m.dfs(m.startX, m.startY) {
		log.Println("Путь найден:", m.path)
		return nil
	}
	return fmt.Errorf("не удалось найти путь")
}

func (m *MazeSolving) dfs(x, y int) bool {
	if x < 0 || x >= m.mazeInfo.Cols || y < 0 || y >= m.mazeInfo.Rows {
		log.Printf("Выход за границы: (%d, %d)\n", x, y)
		return false // Выход за границы
	}
	if m.solvingMatrix[y][x] == 1 {
		log.Printf("Уже посещенная клетка: (%d, %d)\n", x, y)
		return false // Уже посещенная клетка
	}
	if x == m.endX && y == m.endY {
		m.path = append(m.path, Point{X: x, Y: y}) // Добавляем конечную точку в путь
		log.Printf("Найдена конечная точка: (%d, %d)\n", x, y)
		return true // Найден путь
	}

	// Помечаем клетку как посещенную
	m.solvingMatrix[y][x] = 1
	m.path = append(m.path, Point{X: x, Y: y}) // Добавляем текущую точку в путь
	log.Printf("Проверяем клетку: (%d, %d)\n", x, y)

	// Пробуем двигаться в четырех направлениях, проверяя наличие стен
	if x < m.mazeInfo.Cols-1 && !m.mazeInfo.Cells[y][x].Right { // вправо
		if m.dfs(x+1, y) {
			return true
		}
	}
	if x > 0 && !m.mazeInfo.Cells[y][x-1].Right { // влево
		if m.dfs(x-1, y) {
			return true
		}
	}
	if y < m.mazeInfo.Rows-1 && !m.mazeInfo.Cells[y][x].Bottom { // вниз
		if m.dfs(x, y+1) {
			return true
		}
	}
	if y > 0 && !m.mazeInfo.Cells[y-1][x].Bottom { // вверх
		if m.dfs(x, y-1) {
			return true
		}
	}

	m.path = m.path[:len(m.path)-1] // Удаляем текущую точку из пути, если путь не найден
	log.Printf("Возвращаемся из клетки: (%d, %d)\n", x, y)
	return false
}

// Метод для получения найденного пути
func (m *MazeSolving) GetPath() []Point {
	return m.path
}
