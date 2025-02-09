package maze

import (
	"bufio"
	"fmt"
	"go-maze/config"
	"log"
	"os"

	"github.com/sqweek/dialog"
)

func LoadMaze(filename string, windowWidth, windowHeight int) (*Maze, float32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при открытии файла: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Читаем размеры лабиринта
	if !scanner.Scan() {
		return nil, 0, fmt.Errorf("ошибка при чтении размеров лабиринта: %v", scanner.Err())
	}
	var rows, cols int
	_, err = fmt.Sscanf(scanner.Text(), "%d %d", &rows, &cols)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при парсинге размеров лабиринта: %v", err)
	}
	fmt.Printf("Размеры лабиринта: %d строк, %d столбцов\n", rows, cols)

	maze := &Maze{
		Rows:  rows,
		Cols:  cols,
		Cells: make([][]Cell, rows),
	}

	for i := range maze.Cells {
		maze.Cells[i] = make([]Cell, cols)
	}
	// Чтение матриц стенок
	if err := ReadWalls(scanner, maze, rows, cols, "правых"); err != nil {
		return nil, 0, fmt.Errorf("ошибка при чтении правых стен: %v", err)
	}
	// Пропускаем пустую строку между матрицами
	if !scanner.Scan() {
		return nil, 0, fmt.Errorf("ошибка при чтении пустой строки между матрицами: %v", scanner.Err())
	}

	if err := ReadWalls(scanner, maze, rows, cols, "нижних"); err != nil {
		return nil, 0, fmt.Errorf("ошибка при чтении нижних стен: %v", err)
	}

	// Вычисляем размер ячейки
	cellSize := float32(windowWidth) / float32(cols)

	fmt.Println("Загрузка лабиринта завершена успешно.")
	return maze, cellSize, nil
}

// SaveMaze сохраняет лабиринт в файл в указанном формате
func (m *Maze) SaveMaze(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = fmt.Fprintf(file, "%d %d\n", m.Rows, m.Cols); err != nil {
		return err
	}

	if err := m.WriteWalls(file, true); err != nil {
		return err
	}

	if _, err = fmt.Fprintln(file); err != nil {
		return err
	}

	return m.WriteWalls(file, false)
}

func (g *Game) isInsideButton(x, y float32, buttonY float32, buttonHeight float32) bool {
	buttonX := float32(0)
	buttonWidth := float32(config.SceneHeight + config.BorderThickness*2)
	return x >= buttonX && x <= buttonX+buttonWidth && y >= buttonY && y <= buttonY+buttonHeight
}

func (g *Game) ShowFileSelector() {

	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Ошибка при получении текущей директории:", err)
		return
	}

	filename, err := dialog.File().
		Filter("Text files", "txt").
		SetStartDir(currentDir).
		Load()

	if err != nil {
		log.Println("Ошибка при выборе файла:", err)
		return
	}

	// Загружаем лабиринт из выбранного файла
	mazeNew, cellSize, err := LoadMaze(filename, config.SceneWidth, config.SceneHeight)
	if err != nil {
		log.Println("Ошибка при загрузке лабиринта:", err)
		return
	}

	// Обновляем состояние игры с новым лабиринтом
	g.maze = mazeNew
	g.cellSize = cellSize
	g.ResetGame()

}
