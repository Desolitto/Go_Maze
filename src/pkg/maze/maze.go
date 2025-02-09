package maze

type Cell struct {
	Right  bool
	Bottom bool
	Set    int
}

type Maze struct {
	Rows, Cols int
	Cells      [][]Cell
}

func NewMaze(rows, cols int) *Maze {
	// ebiten.SetWindowSize(config.SceneWidth+int(config.BorderThickness*2), config.SceneHeight+config.ButtonHeight+int(config.BorderThickness))
	cells := make([][]Cell, rows)
	for i := range cells {
		cells[i] = make([]Cell, cols)
		for j := range cells[i] {
			cells[i][j] = Cell{Right: false, Bottom: false, Set: -1}
		}
	}
	return &Maze{Rows: rows, Cols: cols, Cells: cells}
}
