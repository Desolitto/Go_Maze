package cave

type Cell int

const (
	Death Cell = iota
	Alive
)

type Cave struct {
	Grid          [][]Cell
	Width, Height int
}

func NewCave(width, height int) *Cave {
	cave := &Cave{
		Grid:   make([][]Cell, height),
		Width:  width,
		Height: height,
	}
	for i := range cave.Grid {
		cave.Grid[i] = make([]Cell, width)
	}
	return cave
}
