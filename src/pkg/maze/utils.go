package maze

import (
	"bufio"
	"fmt"
	"os"
)

func ReadWalls(scanner *bufio.Scanner, maze *Maze, rows, cols int, wallType string) error {
	for y := 0; y < rows; y++ {
		if !scanner.Scan() {
			return fmt.Errorf("ошибка при чтении стенок %s в строке %d: %v", wallType, y, scanner.Err())
		}
		for x := 0; x < cols; x++ {
			var wall int
			_, err := fmt.Sscanf(scanner.Text()[x*2:x*2+1], "%d", &wall) // Предполагаем, что данные разделены пробелами
			if err != nil {
				return fmt.Errorf("ошибка при парсинге стенки %s в строке %d, столбце %d: %v", wallType, y, x, err)
			}
			if wall == 1 {
				if wallType == "правых" {
					maze.Cells[y][x].Right = true
				} else {
					maze.Cells[y][x].Bottom = true
				}
			}
		}
	}
	return nil
}

func (m *Maze) WriteWalls(file *os.File, isRight bool) error {
	for y := 0; y < m.Rows; y++ {
		for x := 0; x < m.Cols; x++ {
			var wall bool
			if isRight {
				wall = m.Cells[y][x].Right
			} else {
				wall = m.Cells[y][x].Bottom
			}
			if _, err := fmt.Fprintf(file, "%d ", boolToInt(wall)); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(file); err != nil {
			return err
		}
	}
	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
