package maze

type Cell int

const (
	Empty Cell = iota
	Wall
)

type Maze struct {
	Grid [][]Cell
}

func NewMaze(width, height int) *Maze {
	maze := &Maze{
		Grid: make([][]Cell, height),
	}
	for i := range maze.Grid {
		maze.Grid[i] = make([]Cell, width)
		for j := range maze.Grid[i] {
			if i%2 == 0 && j%2 == 0 {
				maze.Grid[i][j] = Empty
			} else {
				maze.Grid[i][j] = Wall
			}
		}
	}
	return maze
}
