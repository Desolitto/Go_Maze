package cave

import "fmt"

func (m *Cave) Print() {
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if m.Grid[y][x] == Alive {
				fmt.Print("1 ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}
}
