package cave

func (m *Cave) CountAliveAround(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if (nx < 0 || nx >= m.Width || ny < 0 || ny >= m.Height) || m.Grid[ny][nx] == Alive {
				count++
			}
		}
	}
	return count
}
