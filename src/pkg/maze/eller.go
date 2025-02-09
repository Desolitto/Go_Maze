package maze

func (m *Maze) copyPreviousRow(row int, currentSetCount *int) {
	for col := 0; col < m.Cols; col++ {
		m.Cells[row][col].Right = m.Cells[row-1][col].Right
		m.Cells[row][col].Bottom = m.Cells[row-1][col].Bottom
		m.Cells[row][col].Set = m.Cells[row-1][col].Set
	}
	for col := 0; col < m.Cols; col++ {
		m.Cells[row][col].Right = false
		if m.Cells[row-1][col].Bottom {
			m.Cells[row][col].Set = 0
			m.Cells[row][col].Bottom = false
		}
	}
	for col := 0; col < m.Cols; col++ {
		if m.Cells[row][col].Set == 0 {
			m.Cells[row][col].Set = (*currentSetCount)
			(*currentSetCount)++
		}
	}

}

func (m *Maze) InitializeSets() {
	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			m.Cells[row][col].Set = row*m.Cols + col + 1
		}
	}
}

func (m *Maze) setFirstRowSets(currentSetCount *int) {
	for col := 0; col < m.Cols; col++ {
		m.Cells[0][col].Set = col + 1
		(*currentSetCount)++
	}
}

func (m *Maze) processRightWalls(row int, randomNumbers []int, index *int) {
	for col := 0; col < m.Cols-1; col++ {
		if randomNumbers[*index] == 1 {
			m.Cells[row][col].Right = true
		} else {
			set1 := m.Cells[row][col].Set
			set2 := m.Cells[row][col+1].Set
			if set1 != set2 {
				m.mergeSets(set1, set2)
			} else {
				m.Cells[row][col].Right = true
			}
		}
		(*index)++
	}
}

func (m *Maze) processBottomWalls(row int, randomNumbers []int, index *int) {
	for col := 0; col < m.Cols; col++ {
		set := m.Cells[row][col].Set
		count := 0
		for c := 0; c < m.Cols; c++ {
			if m.Cells[row][c].Set == set && !m.Cells[row][c].Bottom {
				count++
			}
		}

		if count > 1 {
			if randomNumbers[*index] == 1 {
				m.Cells[row][col].Bottom = true
			}
		}
		*index++
	}
}
func (m *Maze) addBottomWalls(row int) {
	for col := 0; col < m.Cols; col++ {
		m.Cells[row][col].Bottom = true
	}
}

func (m *Maze) GenerateEller(randomNumbers []int) {
	m.InitializeSets()

	currentSetCount := 1
	m.setFirstRowSets(&currentSetCount)
	index := 0
	for row := 0; row < m.Rows; row++ {

		if row > 0 {
			m.copyPreviousRow(row, &currentSetCount)
		}

		m.processRightWalls(row, randomNumbers, &index)
		m.processBottomWalls(row, randomNumbers, &index)

		if row == m.Rows-1 {
			m.addBottomWalls(row)
			m.mergeLastRowSets(row)
		}
	}

}

func (m *Maze) mergeLastRowSets(row int) {
	for col := 0; col < m.Cols-1; col++ {
		set1 := m.Cells[row][col].Set
		set2 := m.Cells[row][col+1].Set
		if set1 != set2 {
			m.Cells[row][col].Right = false
			m.mergeSets(set1, set2)
		}
	}
}

func (m *Maze) mergeSets(set1, set2 int) {
	for r := 0; r < m.Rows; r++ {
		for c := 0; c < m.Cols; c++ {
			if m.Cells[r][c].Set == set2 {
				m.Cells[r][c].Set = set1
			}
		}
	}
}

/* ================== old code ========================== */

// func (m *Maze) Initialize(rows, cols int) {
// 	m.Rows = rows
// 	m.Cols = cols
// 	m.Cells = make([][]Cell, rows)

// 	for y := 0; y < rows; y++ {
// 		m.Cells[y] = make([]Cell, cols)
// 		for x := 0; x < cols; x++ {
// 			// Устанавливаем все стенки по умолчанию
// 			m.Cells[y][x].Right = true
// 			m.Cells[y][x].Bottom = true
// 		}
// 	}
// }

// func (m *Maze) Generate(x, y int) {
// 	visited := make([][]bool, m.Rows)
// 	for i := range visited {
// 		visited[i] = make([]bool, m.Cols)
// 	}

// 	stack := []struct{ x, y int }{{x, y}}
// 	visited[y][x] = true

// 	directions := []struct {
// 		dx, dy int
// 	}{
// 		{1, 0},  // вправо
// 		{0, 1},  // вниз
// 		{-1, 0}, // влево
// 		{0, -1}, // вверх
// 	}

// 	for len(stack) > 0 {
// 		curr := stack[len(stack)-1]
// 		stack = stack[:len(stack)-1]

// 		rand.Shuffle(len(directions), func(i, j int) {
// 			directions[i], directions[j] = directions[j], directions[i]
// 		})

// 		for _, dir := range directions {
// 			newX, newY := curr.x+dir.dx, curr.y+dir.dy
// 			if newX >= 0 && newX < m.Cols && newY >= 0 && newY < m.Rows && !visited[newY][newX] {
// 				if dir.dx == 1 { // вправо
// 					m.Cells[curr.y][curr.x].Right = false
// 				} else if dir.dy == 1 { // вниз
// 					m.Cells[curr.y][curr.x].Bottom = false
// 				} else if dir.dx == -1 { // влево
// 					m.Cells[newY][newX].Right = false
// 				} else if dir.dy == -1 { // вверх
// 					m.Cells[newY][newX].Bottom = false
// 				}

// 				visited[newY][newX] = true
// 				stack = append(stack, struct{ x, y int }{newX, newY})
// 			}
// 		}
// 	}
// }
