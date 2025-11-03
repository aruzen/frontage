package pkg

type Size struct {
	Width  int
	Height int
}

type Point struct {
	X int
	Y int
}

func Make2D[T interface{}](size Size, init T) [][]T {
	data := make([]T, size.Height*size.Width)
	grid := make([][]T, size.Width)
	for i := range grid {
		grid[i] = data[i*size.Width : (i+1)*size.Width]
		for j, _ := range grid[i] {
			grid[i][j] = init
		}
	}
	return grid
}

func Copy2D[T interface{}](size Size, o [][]T) [][]T {
	data := make([]T, size.Height*size.Width)
	grid := make([][]T, size.Width)
	for i := range grid {
		grid[i] = data[i*size.Width : (i+1)*size.Width]
		for j := range grid[i] {
			grid[i][j] = o[i][j]
		}
	}
	return grid
}
