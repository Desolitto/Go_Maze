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

func (m *Maze) initializeSets() {
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
	m.initializeSets()

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
